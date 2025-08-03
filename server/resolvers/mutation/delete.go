package mutation

import (
	"fmt"

	"github.com/codecrafter404/bubble/graph"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DeleteOrder is the resolver for the deleteOrder field.
func DeleteOrder(db *gorm.DB, order int) (int, error) {
	//TODO: check admin
	if err := db.Select(clause.Associations).Delete(&graph.Order{}, order).Error; err != nil {
		return 0, fmt.Errorf("Failed to delete order %d: %s", order, err)
	}
	return order, nil
}
