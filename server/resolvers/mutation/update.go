package mutation

import (
	"github.com/codecrafter404/bubble/graph"
	"gorm.io/gorm"
)

// UpdateOrder is the resolver for the updateOrder field.
func UpdateOrder(db *gorm.DB, order int, state graph.OrderState) (*graph.Order, error) {
	panic("not implemented")
}

// UpdateItem is the resolver for the updateItem field.
func UpdateItem(db *gorm.DB, id int, item graph.UpdateItem) (*graph.Item, error) {
	panic("not implemented")
}

// UpdateCustomItem is the resolver for the updateCustomItem field.
func UpdateCustomItem(db *gorm.DB, id int, item graph.UpdateCustomItem) (*graph.CustomItem, error) {
	panic("not implemented")
}
