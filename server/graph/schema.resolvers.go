package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/codecrafter404/bubble/graph/model"
	"github.com/codecrafter404/bubble/utils"
	"github.com/jmoiron/sqlx"
)

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context, order model.NewOrder) (*model.Order, error) {
	itemIds := []int{}
	customItemIds := []int{}

	for _, i := range order.Items {
		found := false
		for _, id := range itemIds {
			if id == i.Item {
				found = true
				break
			}
		}
		if !found {
			itemIds = append(itemIds, i.Item)
		} else {
			return nil, fmt.Errorf("Item %d is a duplicate; please use the quantity field to specify different amounts", i.Item)
		}
	}
	for _, i := range order.CustomItems {
		found := false
		for _, id := range customItemIds {
			if id == i.ID {
				found = true
				break
			}
		}
		if !found {
			customItemIds = append(customItemIds, i.ID)
		}
		for _, v := range i.Variants {
			found := false
			for _, id := range itemIds {
				if id == v {
					found = true
					break
				}
			}
			if !found {
				itemIds = append(itemIds, v)
			}
		}
	}
	items := []model.Item{}
	if len(itemIds) != 0 {

		itemQuery, itemArgs, err := sqlx.In("SELECT id, name, price, image, available, identifier, oneoff FROM item WHERE id IN (?);", itemIds)
		if err != nil {
			return nil, fmt.Errorf("Failed to build query for items: %w", err)
		}

		itemRows, err := r.Db.Query(itemQuery, itemArgs...)
		defer itemRows.Close()
		if err != nil {
			return nil, fmt.Errorf("Failed to query items: %w", err)
		}

		for itemRows.Next() {
			var item model.Item
			err := itemRows.Scan(&item.ID, &item.Name, &item.Price, &item.Image, &item.Available, &item.Identifier, &item.IsOneOff)
			if err != nil {
				return nil, fmt.Errorf("Failed to scan rows: %w", err)
			}
			items = append(items, item)
		}

		// validate items
		for _, o := range order.Items {
			found := false
			for _, i := range items {
				if i.ID == o.Item {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("Couldn't find item with id %d", o.Item)
			}
			if o.Quantity <= 0 {
				return nil, fmt.Errorf("Expected to find quantity >= 0 for %d", o.Item)
			}
		}
	}

	customItems, err := utils.QueryCustomItems(r.Db)
	if err != nil {
		return nil, fmt.Errorf("Failed to query custom items: %w", err)
	}
	graphNodes := []utils.GraphNode{}
	for _, i := range customItems {
		graphNodes = append(graphNodes, utils.GraphNode{Id: i.ID, DependsOn: i.DependsOn})
	}
	// validate customitem
	for _, o := range order.CustomItems {
		if o.Quantity <= 0 {
			return nil, fmt.Errorf("Expected to have quantity >= 0 for customitem %d", o.ID)
		}
		found := false
		var customItem model.CustomItem
		for _, i := range customItems {
			if i.ID == o.ID {
				found = true
				customItem = i
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("Couldn't find customitem with id %d", o.ID)
		}

		if customItem.DependsOn != nil {
			return nil, fmt.Errorf("Expected CustomItem %d to be a toplevel, but go part of a tree (dependsOn = %d)", o.ID, *customItem.DependsOn)
		}
		var graphNode utils.GraphNode
		for _, n := range graphNodes {
			if o.ID == n.Id {
				graphNode = n
				break
			}
		}

		deps, successful := graphNode.ResolveDependency(graphNodes, []utils.GraphNode{})
		if !successful {
			return nil, fmt.Errorf("Invalid db state: CustomItems should have wellformed dependencies")
		}

		validationStages := []model.CustomItem{}
		validationStages = append(validationStages, customItem)
		for _, i := range deps {
			var customItem model.CustomItem
			for _, x := range customItems {
				if i.Id == x.ID {
					customItem = x
					break
				}
			}
			validationStages = append(validationStages, customItem)
		}

		for _, stage := range validationStages {
			found := false
			for _, v := range stage.Variants {
				for _, i := range o.Variants {
					if i == v.ID {
						if stage.Exclusive && found {
							return nil, fmt.Errorf("CustomItem %d (dependency for %d) has non-exclusive variant %d", stage.ID, customItem.ID, i)
						}
						found = true
					}
				}
			}
			if !found {
				return nil, fmt.Errorf("CustomItem %d didn't contain at least one variant for depndend CustomItem %d", customItem.ID, stage.ID)
			}
		}

	}

	identifier, err := utils.GetIdentifier(r.Db, 100)
	if err != nil {
		return nil, fmt.Errorf("Failed to find identifier: %w", err)
	}

	id := rand.Intn(999999999)

	tx, err := r.Db.Begin()
	if err != nil {
		return nil, fmt.Errorf("Failed to begin transaction: %w", err)
	}

	cTime := time.Now().UnixMicro()
	_, exErr := tx.Exec("INSERT INTO orders (id, timestamp, identifier, state, total) VALUES (?, ?, ?, ?, ?)", id, cTime, identifier, model.OrderStateCreated, order.Total)
	if exErr != nil {
		return nil, fmt.Errorf("Failed to insert order: %w", exErr)
	}

	for _, i := range order.Items {
		_, err := tx.Exec("INSERT INTO orders_items_link (order_id, quantity, item_id) VALUES (?, ?, ?)", id, i.Quantity, i.Item)
		if err != nil {
			return nil, fmt.Errorf("Failed to link item %d to order %d", i.Item, id)
		}
	}

	for _, c := range order.CustomItems {
		for _, v := range c.Variants {
			_, err := tx.Exec("INSERT INTO orders_custom_items_link (order_id, custom_item_id, item_id, quantity) VALUES (?, ?, ?, ?)", id, c.ID, v, c.Quantity)
			if err != nil {
				return nil, fmt.Errorf("Failed to link custom_order %d to order %d", c.ID, id)
			}
		}
	}
	resCustomItems := []*model.OrderCustomItem{}
	for _, c := range order.CustomItems {
		var cItem model.CustomItem
		for _, x := range customItems {
			if x.ID == c.ID {
				cItem = x
				break
			}
		}
		variants := []*model.Item{}
		for _, v := range c.Variants {
			for _, i := range items {
				if i.ID == v {
					variants = append(variants, &i)
				}
			}
		}
		cItem.Variants = variants
		resCustomItems = append(resCustomItems, &model.OrderCustomItem{Quantity: c.Quantity, CustomItem: &cItem})

	}
	resItems := []*model.OrderItem{}
	for _, x := range order.Items {
		var item model.Item
		for _, z := range items {
			if z.ID == x.Item {
				item = z
				break
			}
		}
		resItems = append(resItems, &model.OrderItem{Quantity: x.Quantity, Item: &item})
	}
	txErr := tx.Commit()
	if txErr != nil {
		return nil, fmt.Errorf("Failed to commit transaction: %w", txErr)
	}
	resOrder := model.Order{
		ID:          id,
		Timestamp:   cTime,
		Identifier:  strconv.Itoa(identifier),
		State:       model.OrderStateCreated,
		Total:       order.Total,
		CustomItems: resCustomItems,
		Items:       resItems,
	}

	r.OrderChannelMux.RLock()

	for _, c := range r.OrderChannel {
		select {
		case c <- id:
			break
		default:
			log.Println("Failed to order event to channel (create)")
			break
		}
	}

	r.OrderChannelMux.RUnlock()

	return &resOrder, nil
}

// UpdateOrder is the resolver for the updateOrder field.
func (r *mutationResolver) UpdateOrder(ctx context.Context, order int, state model.OrderState) (*model.Order, error) {
	_, err := r.Db.Exec("UPDATE orders SET state = ? WHERE id = ?", state, order)
	if err != nil {
		return nil, fmt.Errorf("Failed to update order: %w", err)
	}

	options := utils.OrderQueryOptions{
		FilterIds: &[]int{order},
	}
	res, err := utils.QueryOrders(r.Db, &options)
	if err != nil {
		return nil, fmt.Errorf("Failed to query order: %w", err)
	}

	if len(res) == 0 {
		return nil, nil
	}

	r.OrderChannelMux.RLock()

	for _, c := range r.OrderChannel {
		select {
		case c <- res[0].ID:
			break
		default:
			log.Println("Failed to order event to channel (update)")
			break
		}
	}

	r.OrderChannelMux.RUnlock()
	return &res[0], nil
}

// DeleteOrder is the resolver for the deleteOrder field.
func (r *mutationResolver) DeleteOrder(ctx context.Context, order int) (int, error) {
	res, err := r.Db.Exec("DELETE FROM orders WHERE id = ?; DELETE FROM orders_custom_items_link WHERE order_id = ?; DELETE FROM orders_items_link WHERE order_id = ?", order, order, order)
	if err != nil {
		return 0, fmt.Errorf("Failed to delete order: %w", err)
	}

	if res, err := res.RowsAffected(); err == nil && res >= 0 {
		r.OrderChannelMux.RLock()

		for _, c := range r.OrderChannel {
			select {
			case c <- order:
				break
			default:
				log.Println("Failed to order event to channel (delete)")
				break
			}
		}

		r.OrderChannelMux.RUnlock()
	}

	return order, nil
}

// UpdateItem is the resolver for the updateItem field.
func (r *mutationResolver) UpdateItem(ctx context.Context, id int, item model.UpdateItem) (*model.Item, error) {
	tx, txErr := r.Db.Begin()
	if txErr != nil {
		return nil, fmt.Errorf("Failed to begin transaction: %w", txErr)
	}

	if item.Available != nil {
		_, updateErr := tx.Exec("UPDATE item SET available = ? WHERE id = ?", item.Available, id)
		if updateErr != nil {
			return nil, fmt.Errorf("Failed to update available: %w", updateErr)
		}
	}

	if item.Identifier != nil {
		_, updateErr := tx.Exec("UPDATE item SET identifier = ? WHERE id = ?", item.Identifier, id)
		if updateErr != nil {
			return nil, fmt.Errorf("Failed to update identifier: %w", updateErr)
		}
	}

	if item.Image != nil {
		_, updateErr := tx.Exec("UPDATE item SET image = ? WHERE id = ?", item.Image, id)
		if updateErr != nil {
			return nil, fmt.Errorf("Failed to update image: %w", updateErr)
		}
	}

	if item.IsOneOff != nil {
		_, updateErr := tx.Exec("UPDATE item SET oneoff = ? WHERE id = ?", item.IsOneOff, id)
		if updateErr != nil {
			return nil, fmt.Errorf("Failed to update IsOneOff: %w", updateErr)
		}
	}

	if item.Name != nil {
		_, updateErr := tx.Exec("UPDATE item SET name = ? WHERE id = ?", item.Name, id)
		if updateErr != nil {
			return nil, fmt.Errorf("Failed to update name: %w", updateErr)
		}
	}

	if item.Price != nil {
		_, updateErr := tx.Exec("UPDATE item SET price = ? WHERE id = ?", item.Price, id)
		if updateErr != nil {
			return nil, fmt.Errorf("Failed to update price: %w", updateErr)
		}
	}
	err := tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("Failed to commit transaction: %w", err)
	}

	rows := r.Db.QueryRow("SELECT id, name, price, image, available, identifier, oneoff FROM item WHERE id = ?", id)
	var dest model.Item
	rows.Scan(&dest.ID, &dest.Name, &dest.Price, &dest.Image, &dest.Available, &dest.Identifier, &dest.IsOneOff)

	r.EventChannelMux.RLock()

	msg := model.UpdateEventUpdateItem
	for _, c := range r.EventChannel {
		select {
		case c <- &msg:
			break
		default:
			log.Println("Failed to UpdateEventUpdateItem event to channel")
			break
		}
	}

	r.EventChannelMux.RUnlock()
	return &dest, nil
}

// UpdateCustomItem is the resolver for the updateCustomItem field.
func (r *mutationResolver) UpdateCustomItem(ctx context.Context, id int, item model.UpdateCustomItem) (*model.CustomItem, error) {
	tx, txErr := r.Db.Begin()
	if txErr != nil {
		return nil, fmt.Errorf("Failed to begin transaction: %w", txErr)
	}

	if item.Exclusive != nil {
		_, updateErr := tx.Exec("UPDATE custom_item SET exclusive = ? WHERE id = ?", item.Exclusive, id)
		if updateErr != nil {
			return nil, fmt.Errorf("Failed to update exclusive: %w", updateErr)
		}
	}

	if item.Name != nil {
		_, updateErr := tx.Exec("UPDATE custom_item SET name = ? WHERE id = ?", item.Name, id)
		if updateErr != nil {
			return nil, fmt.Errorf("Failed to update name: %w", updateErr)
		}
	}
	err := tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("Failed to commit transaction: %w", err)
	}
	res, err := utils.QueryCustomItem(r.Db, id)
	if err != nil {
		return nil, fmt.Errorf("Failed to query orders: %w", err)
	}

	r.EventChannelMux.RLock()

	msg := model.UpdateEventUpdateCustomitem
	for _, c := range r.EventChannel {
		select {
		case c <- &msg:
			break
		default:
			log.Println("Failed to UpdateEventUpdateCustomitem event to channel")
			break
		}
	}

	r.EventChannelMux.RUnlock()

	return res, nil
}

// CreateItems is the resolver for the createItems field.
func (r *mutationResolver) CreateItems(ctx context.Context, items []*model.ItemInput) ([]int, error) {
	uniqueIds := []int{}
	for _, i := range items {
		for _, x := range uniqueIds {
			if x == i.ID {
				return []int{}, fmt.Errorf("Id %d isn't unique", x)
			}
		}
		uniqueIds = append(uniqueIds, i.ID)
	}

	tx, err := r.Db.Begin()
	if err != nil {
		return []int{}, fmt.Errorf("Failed to beginn transaction: %w", err)
	}

	itemIDs := []int{}

	_, delErr := tx.Exec("DELETE FROM item;")
	if delErr != nil {
		return []int{}, fmt.Errorf("Failed to clean up old items: %w", delErr)
	}

	for _, item := range items {
		_, err = tx.Exec("INSERT INTO item (id, name, price, image, available, identifier, oneoff) VALUES (?, ?, ?, ?, ?, ?, ?);", item.ID, item.Name, item.Price, item.Image, item.Available, item.Identifier, item.IsOneOff)
		if err != nil {
			return []int{}, fmt.Errorf("Failed to tx.exec insert for %d: %w", item.ID, err)
		}
		itemIDs = append(itemIDs, item.ID)

	}

	err = tx.Commit()
	if err != nil {
		return []int{}, fmt.Errorf("Failed to commit transaction: %w", err)
	}

	r.EventChannelMux.RLock()

	msg := model.UpdateEventUpdateItem
	for _, c := range r.EventChannel {
		select {
		case c <- &msg:
			break
		default:
			log.Println("Failed to UpdateEventUpdateItem event to channel")
			break
		}
	}

	r.EventChannelMux.RUnlock()

	return itemIDs, nil
}

// CreateCustomItems is the resolver for the createCustomItems field.
func (r *mutationResolver) CreateCustomItems(ctx context.Context, items []*model.CustomItemInput) ([]int, error) {
	graphItems := []utils.GraphNode{}
	uniqueIds := []int{}
	for _, i := range items {
		graphItems = append(graphItems, utils.GraphNode{Id: i.ID, DependsOn: i.DependsOn})
		for _, x := range uniqueIds {
			if x == i.ID {
				return []int{}, fmt.Errorf("Id %d isn't unique", x)
			}
		}
		uniqueIds = append(uniqueIds, i.ID)
	}

	if !utils.CheckDependencyLoop(graphItems) {
		return []int{}, fmt.Errorf("The supplied items have circular/unfeasable dependencies")
	}

	itemIds := []int{}

	rows, err := r.Db.Query("SELECT id FROM item WHERE oneoff == 0")
	defer rows.Close()
	if err != nil {
		return []int{}, fmt.Errorf("Failed to query existing item ids for dependency check: %w", err)
	}

	for rows.Next() {
		var res int

		err = rows.Scan(&res)
		if err != nil {
			return []int{}, fmt.Errorf("Failed to scan result fow (dep check): %w", err)
		}
		itemIds = append(itemIds, res)
	}

	for _, i := range items {
		valid := false
		for _, v := range i.Variants {
			found := false
			for _, x := range itemIds {
				if v == x {
					found = true
				}
			}
			if found {
				valid = true
			} else {
				break
			}
		}

		if !valid {
			return []int{}, fmt.Errorf("customitem %d depends on variant which doesn't exist (yet) or is marked as oneoff", i.ID)
		}
	}

	tx, err := r.Db.Begin()
	if err != nil {
		return []int{}, fmt.Errorf("Failed to begin transaction: %w", err)
	}
	insertedIds := []int{}

	_, delErr := tx.Exec("DELETE FROM custom_item_item_link; DELETE FROM custom_item")
	if delErr != nil {
		return []int{}, fmt.Errorf("failed to cleanup item: %w", delErr)
	}

	for _, i := range items {

		_, inErr := tx.Exec("INSERT INTO custom_item(id, name, depends_on, exclusive) VALUES (?, ?, ?, ?)", i.ID, i.Name, i.DependsOn, i.Exclusive)
		if inErr != nil {
			return []int{}, fmt.Errorf("failed to insert customitem %d: %w", i.ID, delErr)
		}

		for _, v := range i.Variants {
			_, vErr := tx.Exec("INSERT INTO custom_item_item_link(custom_item_id, item_id) VALUES (?, ?)", i.ID, v)
			if vErr != nil {
				return []int{}, fmt.Errorf("failed to insert variant_link %d for %d: %w", v, i.ID, delErr)
			}
		}
		insertedIds = append(insertedIds, i.ID)

	}
	comErr := tx.Commit()
	if comErr != nil {
		return []int{}, fmt.Errorf("Failed to commit transaction: %w", comErr)
	}

	r.EventChannelMux.RLock()

	msg := model.UpdateEventUpdateCustomitem
	for _, c := range r.EventChannel {
		select {
		case c <- &msg:
			break
		default:
			log.Println("Failed to UpdateEventUpdateCustomitem event to channel")
			break
		}
	}

	r.EventChannelMux.RUnlock()

	return insertedIds, nil
}

// GetPermission is the resolver for the getPermission field.
func (r *queryResolver) GetPermission(ctx context.Context) (model.User, error) {
	//TODO: implement authentication
	return model.UserUser, nil
}

// GetOrder is the resolver for the getOrder field.
func (r *queryResolver) GetOrder(ctx context.Context, id int) (*model.Order, error) {
	options := utils.OrderQueryOptions{
		FilterIds: &[]int{id},
	}
	order, err := utils.QueryOrders(r.Db, &options)
	if err != nil {
		return nil, fmt.Errorf("Failed to query order: %w", err)
	}

	if len(order) == 0 {
		return nil, nil
	}

	return &order[0], nil
}

// GetItems is the resolver for the getItems field.
func (r *queryResolver) GetItems(ctx context.Context) ([]*model.Item, error) {
	rows, err := r.Db.Query("SELECT id, name, price, image, available, identifier, oneoff FROM item")
	defer rows.Close()
	if err != nil {
		return []*model.Item{}, fmt.Errorf("Failed to query db: %w", err)
	}
	res := []*model.Item{}
	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Image, &item.Available, &item.Identifier, &item.IsOneOff)
		if err != nil {
			return []*model.Item{}, fmt.Errorf("Failed to scan rows: %w", err)
		}
		res = append(res, &item)
	}
	return res, nil
}

// GetCustomItems is the resolver for the getCustomItems field.
func (r *queryResolver) GetCustomItems(ctx context.Context) ([]*model.CustomItem, error) {
	items, err := utils.QueryCustomItems(r.Db)
	if err != nil {
		return []*model.CustomItem{}, fmt.Errorf("Failed to query customitem: %w", err)
	}
	res := []*model.CustomItem{}

	for _, i := range items {
		res = append(res, &i)
	}

	return res, nil
}

// Orders is the resolver for the orders field.
func (r *subscriptionResolver) Orders(ctx context.Context, state *model.OrderState, id *int, limit *int, skip *int, sortAsc *bool) (<-chan []*model.Order, error) {
	resChannel := make(chan []*model.Order, 7)
	orderChannel := make(chan int, 7)

	r.OrderChannelMux.RLock()
	r.OrderChannel = append(r.OrderChannel, orderChannel)
	r.OrderChannelMux.RUnlock()

	go func() {

		cleanup := func() {

			r.OrderChannelMux.Lock()
			res := []chan int{}
			for _, c := range r.OrderChannel {
				if c != orderChannel {
					res = append(res, c)
				}
			}
			r.OrderChannel = res
			r.OrderChannelMux.Unlock()
			close(orderChannel)
			close(resChannel)
		}

		fRes, err := utils.QueryOrdersLimited(r.Db, state, id, limit, skip, sortAsc)
		if err != nil {
			log.Println(fmt.Sprintf("ERROR: %s", err))
			cleanup()
			return
		}
		resChannel <- fRes

	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case x := <-orderChannel:
				_ = x

				res, err := utils.QueryOrdersLimited(r.Db, state, id, limit, skip, sortAsc)
				if err != nil {
					log.Println(fmt.Sprintf("ERROR: %s", err))
					break loop
				}
				resChannel <- res
			}
		}
		cleanup()

	}()
	return resChannel, nil
}

// NextOrder is the resolver for the nextOrder field.
func (r *subscriptionResolver) NextOrder(ctx context.Context) (<-chan *model.Order, error) {
	channel := make(chan int, 7)

	r.OrderChannelMux.Lock()

	r.OrderChannel = append(r.OrderChannel, channel)

	r.OrderChannelMux.Unlock()

	resChannel := make(chan *model.Order, 7)
	go func() {

		notifyUpdate := func(id int) {
			r.OrderChannelMux.RLock()

			for _, c := range r.OrderChannel {
				select {
				case c <- id:
					break
				default:
					log.Println("Failed to order event to channel (subscription)")
					break
				}
			}

			r.OrderChannelMux.RUnlock()

		}
		cleanup := func(order *model.Order) {
			r.OrderChannelMux.Lock()
			res := []chan int{}
			for _, c := range r.OrderChannel {
				if c != channel {
					res = append(res, c)
				}
			}
			r.OrderChannel = res
			r.OrderChannelMux.Unlock()
			close(channel)
			close(resChannel)
			if order != nil {
				_, err := r.Db.Exec("UPDATE orders SET state = IIF(state = ?, ?, state) WHERE id = ?", model.OrderStatePending, model.OrderStateCreated, order.ID)
				if err != nil {
					log.Printf("Failed to update order refrain %d: %v", order.ID, err)
				} else {
					notifyUpdate(order.ID)
				}

			}
			r.OrderNextMux.TryLock()
			r.OrderNextMux.Unlock()
		}

		r.OrderNextMux.Lock()
		order, err := utils.GetNextOrder(r.Db)
		r.OrderNextMux.TryLock() // just to be sure
		r.OrderNextMux.Unlock()

		if err != nil {
			log.Printf("ERROR: Failed to query next order (outer): %v\n", err)
			cleanup(order)
			return
		}

		if order != nil {
			log.Printf("Got order %d\n", order.ID)
			notifyUpdate(order.ID)
		}
		log.Println("Checking for updates")

		resChannel <- order

	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case x := <-channel:
				log.Printf("got event %d\n", x)
				if order != nil && x != order.ID {
					continue
				}

				if order != nil {
					opts := utils.OrderQueryOptions{
						FilterIds: &[]int{order.ID},
					}
					orders, err := utils.QueryOrders(r.Db, &opts)
					if err != nil {
						log.Printf("ERROR: Failed to query order: %v\n", err)
						break loop
					}
					if len(orders) == 1 && (orders[0].State == model.OrderStatePending || orders[0].State == model.OrderStateCreated) {
						order = &orders[0]
						resChannel <- &orders[0]
						continue
					} else if len(orders) == 0 {
						order = nil
					}
				}

				r.OrderNextMux.Lock()
				order, err = utils.GetNextOrder(r.Db)
				r.OrderNextMux.Unlock()

				if err != nil {
					log.Printf("ERROR: Failed to query next order: %v\n", err)
					break loop
				}

				if order != nil {
					notifyUpdate(order.ID)
				}

				resChannel <- order
			}
		}

		cleanup(order)
	}()
	return resChannel, nil
}

// Updates is the resolver for the updates field.
func (r *subscriptionResolver) Updates(ctx context.Context) (<-chan *model.UpdateEvent, error) {
	channel := make(chan *model.UpdateEvent, 7)

	r.EventChannelMux.Lock()

	r.EventChannel = append(r.EventChannel, channel)

	r.EventChannelMux.Unlock()

	resChannel := make(chan *model.UpdateEvent, 7)
	go func() {
		for {
			select {
			case <-ctx.Done():
				r.EventChannelMux.Lock()
				res := []chan *model.UpdateEvent{}
				for _, c := range r.EventChannel {
					if c != channel {
						res = append(res, c)
					}
				}
				r.EventChannel = res
				r.EventChannelMux.Unlock()
				close(channel)
				close(resChannel)
				return
			case x := <-channel:
				resChannel <- x
			}
		}
	}()

	return resChannel, nil
}

// Stats is the resolver for the stats field.
func (r *subscriptionResolver) Stats(ctx context.Context) (<-chan *model.Statistics, error) {
	resChannel := make(chan *model.Statistics, 7)
	orderChannel := make(chan int, 7)

	r.OrderChannelMux.RLock()
	r.OrderChannel = append(r.OrderChannel, orderChannel)
	r.OrderChannelMux.RUnlock()

	go func() {

		cleanup := func() {

			r.OrderChannelMux.Lock()
			res := []chan int{}
			for _, c := range r.OrderChannel {
				if c != orderChannel {
					res = append(res, c)
				}
			}
			r.OrderChannel = res
			r.OrderChannelMux.Unlock()
			close(orderChannel)
			close(resChannel)
		}
		stats, err := utils.GetStats(r.Db)
		if err != nil {
			log.Printf("Failed to query stats: %v\n", err)
			cleanup()
			return
		}

		resChannel <- &stats

	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case x := <-orderChannel:
				_ = x

				stats, err := utils.GetStats(r.Db)
				if err != nil {
					log.Printf("Failed to query stats: %v\n", err)
					break loop
				}

				resChannel <- &stats
			}
		}
		cleanup()

	}()
	return resChannel, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
