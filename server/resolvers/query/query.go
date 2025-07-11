package query

import (
	"context"

	"github.com/codecrafter404/bubble/graph"
)

// GetPermission is the resolver for the getPermission field.
func GetPermission(ctx context.Context) (graph.User, error) {
	return graph.UserUser, nil
}

// GetOrder is the resolver for the getOrder field.
func GetOrder(ctx context.Context, id int) (*graph.Order, error) {
	panic("not implemented")
}

// GetItems is the resolver for the getItems field.
func GetItems(ctx context.Context) ([]*graph.Item, error) {
	panic("not implemented")
}

// GetCustomItems is the resolver for the getCustomItems field.
func GetCustomItems(ctx context.Context) ([]*graph.CustomItem, error) {
	panic("not implemented")
}
