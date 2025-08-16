package gql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/codecrafter404/bubble/config"
	"github.com/codecrafter404/bubble/ent"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	client       *ent.Client
	config       *config.Config
	notification []*chan NotificationType
}

func NewSchema(client *ent.Client, config *config.Config) graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: &Resolver{client, config, make([]*chan NotificationType, 0)},
	})
}
