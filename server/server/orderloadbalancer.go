package server

import (
	"context"
	"slices"

	"github.com/codecrafter404/bubble/config"
	"github.com/codecrafter404/bubble/ent"
	"github.com/codecrafter404/bubble/ent/order"
	"github.com/codecrafter404/bubble/gql"
	"github.com/codecrafter404/bubble/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func OrderBalancer(ctx context.Context, client *ent.Client, notifier <-chan gql.NotificationType, subscribers *[]*utils.Subscriber, disconnect <-chan uuid.UUID, globalConfig *config.Config) {

	if client == nil || subscribers == nil || globalConfig == nil {
		log.Fatal().Msg("Client, Subscribers and Config should not be nil")
	}

	// query db for unprocessed orders
	lastWorker := uuid.New()

	// implement load balancing strategy -> implement dangling order
	// check if subscriber is free -> after order has been updated

	for {
		select {
		case notification := <-notifier:
			switch notification {
			case gql.NotificationTypeNewSubscriber:
			case gql.NotificationTypeNewOrder:
			case gql.NotificationTypeOrderUpdated:
				// check if order has been compleated
				// if true remove order id from subscriber -> worker is now free

				for _, sub := range *subscribers {
					if sub.Processing != nil {
						res, err := client.Order.Query().Where(order.And(order.IDEQ(*sub.Processing), order.StateIn(order.StateCancelled, order.StateCompleated))).Exist(ctx)
						if err != nil {
							log.Err(err).Msg("Failed to query updated orders")
							continue
						}
						if res {
							sub.Processing = nil
							sub.Reciever <- nil
						}
					}
				}

			default:
				continue
			}
		case uuid := <-disconnect:
			// reset order to created -> add to orders -> sort by submitted -> remove from subscribers
			subscriberIndex := slices.IndexFunc(*subscribers, func(sub *utils.Subscriber) bool {
				return sub.Id == uuid
			})
			if subscriberIndex == -1 {
				log.Error().Str("uuid", uuid.String()).Msg("Failed to find subscriber")

				continue
			}
			orderId := (*subscribers)[subscriberIndex].Processing
			if orderId != nil {
				err := client.Order.UpdateOneID(*orderId).SetState(order.StateCreated).Exec(ctx)
				if err != nil {
					log.Error().Err(err).Int("order_id", *orderId).Msg("Failed to reset order state")
					continue
				}
			}

			*subscribers = slices.Delete((*subscribers), subscriberIndex, subscriberIndex+1)
		}
		// load balance order

		switch globalConfig.OrderConfig.LoadBalancingStrategy {
		case config.LoadBalancingRoundRobin:
			loadBalanceRoundRobin(subscribers, &lastWorker, client, ctx)
		case config.LoadBalancingRandom:
			loadBalanceRandom(subscribers, &lastWorker, *globalConfig)
		}

	}

}

func loadBalanceRandom(subscribers *[]*utils.Subscriber, last *uuid.UUID, globalConfig config.Config) {
	log.Fatal().Msg("Random Load balancing isn't currently implemented")
	// usable := 0
	// for _, sub := range *subscribers {
	// 	if sub.Processing == nil {
	// 		usable += 1
	// 	}
	// }
	//
	// if usable == 0 {
	// 	// none found
	// 	log.Trace().Msg("When a new notification is received the order will be fulfilled")
	// 	return
	// }
	// selected := rand.Intn(usable)
	//
	// count := 0
	//
	// for _, sub := range *subscribers {
	// 	if sub.Reciever == nil {
	// 		if count == selected {
	// 			// check if met requirements
	// 		}else {
	// 			count += 1
	// 		}
	// 	}
	// }

}

func loadBalanceRoundRobin(subscribers *[]*utils.Subscriber, last *uuid.UUID, client *ent.Client, ctx context.Context) {

	lastIdx := slices.IndexFunc(*subscribers, func(sub *utils.Subscriber) bool { return sub.Id == *last })

	if lastIdx == -1 {
		lastIdx = 0
	} else {
		lastIdx++ // in order to select the next
	}

	noFreeOneFound := 0
	for i := lastIdx; ; i++ {
		if len(*subscribers) == 0 {
			log.Trace().Msg("Currently ther're no subscribers available")
			break
		}
		idx := i % len(*subscribers)

		processing, err := client.Order.Query().Where(order.StateEQ(order.StateCreated)).Order(order.BySubmitted(), ent.Asc()).First(ctx)
		if err != nil {
			if _, ok := err.(*ent.NotFoundError); ok {
				break // no orders to process
			} else {
				log.Error().Err(err).Msg("Failed to query order to process")
				continue
			}
		}

		if (*subscribers)[idx].Processing != nil {
			if noFreeOneFound > len(*subscribers) {
				log.Trace().Msg("Currently orders cant be fulfilled")
				break
			}
			noFreeOneFound++
			continue
		} else {
			noFreeOneFound = 0
		}

		processing, err = processing.Update().SetState(order.StatePending).Save(ctx)

		if err != nil {
			log.Error().Err(err).Msg("Failed to update order state")
			continue
		}

		(*subscribers)[idx].Processing = &processing.ID
		(*subscribers)[idx].Reciever <- processing

		*last = (*subscribers)[idx].Id
	}

}
