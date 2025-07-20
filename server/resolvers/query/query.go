package query

import (
	"context"
	"fmt"

	"github.com/codecrafter404/bubble/graph"
	"gorm.io/gorm"
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
func GetItems(db *gorm.DB) ([]*graph.Item, error) {
	var items []*graph.Item
	if err := db.Find(&items); err != nil {
		return []*graph.Item{}, fmt.Errorf("Failed to fetch items: %+v", err)
	}
	return items, nil
}

// GetCustomItems is the resolver for the getCustomItems field.
func GetCustomItems(db *gorm.DB) ([]*graph.CustomItem, error) {
	var customitems []*graph.CustomItem

	db.Model(&customitems).Association("Items").Find(&customitems)
}
