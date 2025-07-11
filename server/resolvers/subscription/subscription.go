package subscription

import (
	"context"

	"github.com/codecrafter404/bubble/graph"
)

// Orders is the resolver for the orders field.
func Orders(ctx context.Context, state *graph.OrderState, id *int, limit *int, skip *int, sortAsc *bool) (<-chan []*graph.Order, error) {
	panic("not implemented")
}

// NextOrder is the resolver for the nextOrder field.
func NextOrder(ctx context.Context) (<-chan *graph.Order, error) {
	panic("not implemented")
}

// Updates is the resolver for the updates field.
func Updates(ctx context.Context) (<-chan *graph.UpdateEvent, error) {
	panic("not implemented")
}

// Stats is the resolver for the stats field.
func Stats(ctx context.Context) (<-chan *graph.Statistics, error) {
	panic("not implemented")
}
