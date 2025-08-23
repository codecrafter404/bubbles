package gql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/codecrafter404/bubble/config"
	"github.com/codecrafter404/bubble/ent"
	"github.com/codecrafter404/bubble/utils"
	"github.com/google/uuid"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	client       *ent.Client
	config       *config.Config
	notification []*chan NotificationType
	subscribers  *[]*utils.Subscriber
	disconnect   chan uuid.UUID
}

func NewSchema(client *ent.Client, config *config.Config, subscribers *[]*utils.Subscriber, disconnect chan uuid.UUID, notifier chan NotificationType) graphql.ExecutableSchema {
	notifications := make([]*chan NotificationType, 0)
	notifications = append(notifications, &notifier)
	return NewExecutableSchema(Config{
		Resolvers: &Resolver{client, config, notifications, subscribers, disconnect},
	})
}
