package mutation

import (
	"context"

	"github.com/codecrafter404/bubble/graph"
)

func CreateOrder(ctx context.Context, order graph.NewOrder) (*graph.Order, error) {
	panic("not implemented")
}

// CreateItems is the resolver for the createItems field.
func CreateItems(ctx context.Context, items []*graph.NewItem) ([]int, error) {
	panic("not implemented")
}

// CreateCustomItems is the resolver for the createCustomItems field.
func CreateCustomItems(ctx context.Context, items []*graph.NewCustomItem) ([]int, error) {
	panic("not implemented")
}
