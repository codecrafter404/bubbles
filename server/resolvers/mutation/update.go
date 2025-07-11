package mutation

import (
	"context"

	"github.com/codecrafter404/bubble/graph"
)

// UpdateOrder is the resolver for the updateOrder field.
func UpdateOrder(ctx context.Context, order int, state graph.OrderState) (*graph.Order, error) {
	panic("not implemented")
}

// UpdateItem is the resolver for the updateItem field.
func UpdateItem(ctx context.Context, id int, item graph.UpdateItem) (*graph.Item, error) {
	panic("not implemented")
}

// UpdateCustomItem is the resolver for the updateCustomItem field.
func UpdateCustomItem(ctx context.Context, id int, item graph.UpdateCustomItem) (*graph.CustomItem, error) {
	panic("not implemented")
}
