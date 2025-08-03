package subscription

import (
	"github.com/codecrafter404/bubble/graph"
	"gorm.io/gorm"
)

// Orders is the resolver for the orders field.
func Orders(db *gorm.DB, state *graph.OrderState, id *int, limit *int, skip *int, sortAsc *bool) (<-chan []*graph.Order, error) {
	panic("not implemented")
}

// NextOrder is the resolver for the nextOrder field.
func NextOrder(db *gorm.DB) (<-chan *graph.Order, error) {
	panic("not implemented")
}

// Updates is the resolver for the updates field.
func Updates(db *gorm.DB) (<-chan *graph.UpdateEvent, error) {
	panic("not implemented")
}

// Stats is the resolver for the stats field.
func Stats(db *gorm.DB) (<-chan *graph.Statistics, error) {
	panic("not implemented")
}
