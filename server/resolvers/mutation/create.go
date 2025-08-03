package mutation

import (
	"fmt"
	"strconv"
	"time"

	"slices"

	"github.com/codecrafter404/bubble/config"
	"github.com/codecrafter404/bubble/graph"
	"github.com/codecrafter404/bubble/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateOrder(db *gorm.DB, orderChannel chan graph.Order, config config.Config, order graph.NewOrder) (*graph.Order, error) {
	// checking duplicate Items
	var itemIds []int
	for _, i := range order.Items {
		if slices.Contains(itemIds, i.ItemID) {
			return nil, fmt.Errorf("Item %d is duplicate", i.ItemID)
		}
		itemIds = append(itemIds, i.ItemID)
	}

	// checking duplicate custom Items
	var customItemIds []int
	for _, i := range order.CustomItems {
		if slices.Contains(customItemIds, i.CustomItemID) {
			return nil, fmt.Errorf("CustomItem %d is duplicate", i.CustomItemID)
		}
		customItemIds = append(customItemIds, i.CustomItemID)
	}

	// querying all items
	var items []graph.Item
	if err := db.Where("id IN ?", itemIds).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("Failed to query items: %s", err)
	}

	var customItems []graph.CustomItem
	if err := db.Preload(clause.Associations).Where("id IN ?", customItemIds).Find(&customItems).Error; err != nil {
		return nil, fmt.Errorf("Failed to query custom items: %s", err)
	}

	var finalItems []*graph.OrderItem
	// checking for thier existance
	for _, i := range order.Items {
		var found *graph.Item
		for _, j := range items {
			if i.ItemID == j.ID {
				found = &j
				break
			}
		}

		if found == nil {
			return nil, fmt.Errorf("Item %d doesn't exist", i.ItemID)
		}
		if i.Quantity <= 0 {
			return nil, fmt.Errorf("Item %d must at least have one item", i.ItemID)
		}

		finalItems = append(finalItems, &graph.OrderItem{
			Item:     found,
			Quantity: i.Quantity,
		})
	}

	var finalCustomItems []*graph.OrderCustomItem
	var linkedItems []utils.LinkedNode
	for _, i := range customItemIds {
		var found *graph.CustomItem
		for _, j := range customItems {
			if i == j.ID {
				found = &j
				break
			}
		}

		if found == nil {
			return nil, fmt.Errorf("CustomItem %d doesn't exist", i)
		}

		linkedItems = append(linkedItems, utils.LinkedNode{
			Id:        found.ID,
			DependsOn: found.DependsOn,
		})
		var input graph.NewOrderCustomItem
		for _, j := range order.CustomItems {
			if j.CustomItemID == i {
				input = *j
			}
		}

		var finalItems []*graph.Item

		for _, vId := range input.Variants {
			var item *graph.Item
			for _, x := range found.Variants {
				if vId == x.ID {
					item = x
					break
				}
			}
			if item == nil {
				return nil, fmt.Errorf("Item %d doesn't exist or isn't a variant of custom item %d", vId, i)
			}

			finalItems = append(finalItems, item)
		}

		if found.Exclusive && len(finalItems) > 1 {
			return nil, fmt.Errorf("CustomItem %d is not exclusive", i)
		}

		if len(finalItems) <= 0 {
			return nil, fmt.Errorf("CustomItem %d has zero variants", i)
		}

		if input.Quantity <= 0 {
			return nil, fmt.Errorf("Quantity of CustomItem %d is <= 0", i)
		}

		finalCustomItems = append(finalCustomItems, &graph.OrderCustomItem{
			CustomItem: &graph.CustomItem{
				ID:        found.ID,
				DependsOn: found.DependsOn,
				Exclusive: found.Exclusive,
				Name:      found.Name,
				Variants:  finalItems,
			},
			Quantity: input.Quantity,
		})
	}

	// NOTE: check if all dependencies of the customitem exists (check if there is a chain)
	var deps []utils.LinkedNode
	for _, v := range linkedItems {
		// top level node
		if v.DependsOn == nil {
			localDeps, success := v.ResolveDependency(linkedItems, []utils.LinkedNode{}) //NOTE: nu uhh, doesn't work
			if !success {
				return nil, fmt.Errorf("Top level CustomItem %d is cyclic there isn't a compleate dependency path", v.Id)
			}

			// append only existent
			for _, l := range localDeps {
				found := false
				for _, d := range deps {
					if l.Id == d.Id {
						found = true
						break
					}
				}
				if !found {
					deps = append(deps, l)
				}
			}
		}
	}

	for _, i := range finalCustomItems {
		found := false
		for _, j := range deps {
			if i.ID == j.Id {
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("Orphran CustomItem with id %d", i.CustomItem.ID)
		}
	}

	var finalOrder graph.Order
	state := graph.OrderStateCreated
	if order.State != nil {
		state = *order.State
	}

	// identifier
	err := db.Transaction(func(tx *gorm.DB) error {
		identifier := "0"
		var last []graph.Order
		if err := tx.Order("timestamp DESC").Limit(1).Find(&last).Error; err != nil {
			return fmt.Errorf("Failed to query last order: %s", err)
		}

		if len(last) >= 1 {
			x, err := strconv.Atoi(last[0].Identifier)
			if err != nil {
				return fmt.Errorf("Failed to convert identifier")
			}
			identifier = fmt.Sprintf("%d", ((x + 1) % config.OrderConfig.MaximalIdentifiers))
		}

		finalOrder = graph.Order{
			Timestamp:   time.Now().UnixMicro(),
			Identifier:  identifier,
			State:       state,
			Total:       graph.CalculateTotal(finalCustomItems, finalItems),
			Items:       finalItems,
			CustomItems: finalCustomItems,
		}
		if err := tx.Create(&finalOrder).Error; err != nil {
			return fmt.Errorf("Failed to create order (2): %s", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to create order (1): %s", err)
	}

	if state == graph.OrderStateCreated || state == graph.OrderStatePending {
		orderChannel <- finalOrder // NOTE: blocks, when the channel buffer is full
	}

	return &finalOrder, nil

}

// CreateItems is the resolver for the createItems field.
func CreateItems(db *gorm.DB, eventChannel []chan *graph.UpdateEvent, items []*graph.NewItem) ([]int, error) {
	var uniqueItems []*graph.Item
	var ids []int

	for _, x := range items {
		for _, y := range uniqueItems {
			if x.ID == y.ID {
				return nil, fmt.Errorf("ID of item %d is not unique", x.ID)
			}
		}

		uniqueItems = append(uniqueItems, &graph.Item{
			ID:         x.ID,
			Available:  x.Available,
			Identifier: x.Identifier,
			Image:      x.Image,
			IsVariant:  x.IsVariant,
			Name:       x.Name,
			Price:      x.Price,
		})
		ids = append(ids, x.ID)
	}

	if len(uniqueItems) <= 0 {
		return nil, fmt.Errorf("Provide at least one valid Item")
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select(clause.Associations).Delete(&graph.Item{}, ids).Error; err != nil {
			return fmt.Errorf("Failed to clean up old records: %s", err)
		}
		if err := tx.Create(&uniqueItems).Error; err != nil {
			return fmt.Errorf("Failed to create items: %s", err)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to create items: %s", err)
	}

	msg := graph.UpdateEventUpdateItem
	for _, c := range eventChannel {
		c <- &msg // NOTE: blocks, when the channel buffer is full
	}

	return ids, nil
}

// CreateCustomItems is the resolver for the createCustomItems field.
func CreateCustomItems(db *gorm.DB, eventChannel []chan *graph.UpdateEvent, items []*graph.NewCustomItem) ([]int, error) {
	// check if variants exists

	var itemIds []int64
	var customItemDependencies []int
	var nodes []utils.LinkedNode

	for _, x := range items {
		// check variants
		for _, id := range x.Variants {
			for _, ex := range itemIds {
				if ex == int64(id) {
					continue
				}
				itemIds = append(itemIds, int64(id))
			}
		}

		//
		for _, n := range nodes {
			if n.Id == x.ID {
				return nil, fmt.Errorf("Duplicate CustomItem %d", x.ID)
			}

		}

		nodes = append(nodes, utils.LinkedNode{Id: x.ID, DependsOn: x.DependsOn})
		if x.DependsOn != nil {
			customItemDependencies = append(customItemDependencies, *x.DependsOn)
		}
	}

	var existingItemIds []int64

	if err := db.Model(&graph.Item{}).Where("id IN ?", itemIds).Select("id").Scan(&existingItemIds).Error; err != nil {
		return nil, fmt.Errorf("Failed to check if variants exist: %s", err)
	}

	for _, i := range itemIds {
		if !slices.Contains(existingItemIds, i) {
			return nil, fmt.Errorf("Variant %d doesn't exist", i)
		}
	}

	// check if customItems dependencies are present
	var notMet []int

	for _, x := range customItemDependencies {
		met := slices.ContainsFunc(nodes, func(n utils.LinkedNode) bool {
			return n.Id == x
		})

		if !met {
			notMet = append(notMet, x)
		}
	}

	var met []int

	if err := db.Model(&graph.CustomItem{}).Where("id = ?", notMet).Select("id").Scan(&met).Error; err != nil {
		return nil, fmt.Errorf("Failed to query customItems: %s", err)
	}

	for _, nm := range notMet {
		if !slices.Contains(met, nm) {
			return nil, fmt.Errorf("Dependency (%d) of CustomItem can not be satisfied: the CustomItem does'nt exist", nm)
		}
	}

	// check dependency loop
	if !utils.CheckDependencyLoop(nodes) {
		return nil, fmt.Errorf("The supplied CustomItems have circular/unfeasable dependencies")
	}

	// Delete Existing

	var ids []int

	for _, x := range items {
		ids = append(ids, x.ID)
	}

	err := db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Select(clause.Associations).Delete(&graph.CustomItem{}, ids).Error; err != nil {
			return fmt.Errorf("Failed to clean up old customItems: %s", err)
		}

		var res []graph.CustomItem

		for _, x := range items {
			var variants []*graph.Item
			if err := tx.Where("id IN ?", x.Variants).Find(&variants).Error; err != nil {
				return fmt.Errorf("Failed to fetch item (%v): %s", x.Variants, err)
			}

			res = append(res, graph.CustomItem{
				ID:        x.ID,
				DependsOn: x.DependsOn,
				Exclusive: x.Exclusive,
				Name:      x.Name,
				Variants:  variants,
			})
		}

		tx.Create(&res)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to create custom items: %s", err)
	}

	msg := graph.UpdateEventUpdateCustomitem
	for _, c := range eventChannel {
		c <- &msg // NOTE: blocks, when the channel buffer is full
	}

	return ids, nil
}
