package utils

import (
	"database/sql"
	"fmt"

	"github.com/codecrafter404/bubble/graph/model"
)

func MigrateDb(connection *sql.DB) error {
	migration := `
CREATE TABLE IF NOT EXISTS item(
	id INTEGER NOT NULL PRIMARY KEY,
	name TEXT NOT NULL,
	price REAL NOT NULL,
	image TEXT NOT NULL,
	available INTEGER NOT NULL,
	identifier TEXT NOT NULL,
	oneoff INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS custom_item(
	id INTEGER NOT NULL PRIMARY KEY,
	name TEXT NOT NULL,
	depends_on INTEGER,
	exclusive INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS custom_item_item_link(
	custom_item_id INTEGER NOT NULL,
	item_id INTEGER NOT NULL,
	PRIMARY KEY(custom_item_id, item_id)
);

CREATE TABLE IF NOT EXISTS orders(
	id INTEGER NOT NULL PRIMARY KEY,
	timestamp TEXT NOT NULL,
	identifier TEXT NOT NULL,
	state TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders_items_link(
	order_id INTEGER NOT NULL,
	quantity INTEGER NOT NULL,
	item_id INTEGER NOT NULL,
	PRIMARY KEY(order_id, quantity, item_id)
);
CREATE TABLE IF NOT EXISTS orders_custom_items_link(
	order_id INTEGER NOT NULL,
	quantity INTEGER NOT NULL,
	custom_item_id INTEGER NOT NULL,
	item_id INTEGER NOT NULL
);
	`
	_, err := connection.Exec(migration)
	if err != nil {
		return fmt.Errorf("Failed to run migration: %w", err)
	}
	return nil
}
func QueryCustomItems(db *sql.DB) ([]model.CustomItem, error) {
	rows, err := db.Query(`SELECT custom_item.id, custom_item.name, custom_item.depends_on, custom_item.exclusive, item.id, item.name, item.price, item.image, item.available, item.identifier, item.oneoff
		FROM custom_item INNER JOIN custom_item_item_link ON custom_item.id=custom_item_item_link.custom_item_id
		INNER JOIN item ON custom_item_item_link.item_id=item.id`)
	if err != nil {
		return []model.CustomItem{}, fmt.Errorf("Failed to query custom_items: %w", err)
	}
	res, err := parseCustomItemRows(rows)
	if err != nil {
		return []model.CustomItem{}, fmt.Errorf("Failed to parse custom item: %w", err)
	}
	return res, nil
}
func QueryCustomItem(db *sql.DB, id int) (*model.CustomItem, error) {
	rows, err := db.Query(`SELECT custom_item.id, custom_item.name, custom_item.depends_on, custom_item.exclusive, item.id, item.name, item.price, item.image, item.available, item.identifier, item.oneoff
		FROM custom_item INNER JOIN custom_item_item_link ON custom_item.id=custom_item_item_link.custom_item_id
		INNER JOIN item ON custom_item_item_link.item_id=item.id
		WHERE custom_item.id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("Failed to query custom_items: %w", err)
	}
	items, err := parseCustomItemRows(rows)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse custom item rows: %w", err)
	}
	if len(items) != 1 {
		return nil, nil
	}
	return &items[0], nil
}
func parseCustomItemRows(rows *sql.Rows) ([]model.CustomItem, error) {
	customItemMap := make(map[int]model.CustomItem)
	itemMap := make(map[int][]model.Item)

	for rows.Next() {
		var customItem model.CustomItem
		var item model.Item

		err := rows.Scan(&customItem.ID, &customItem.Name, &customItem.DependsOn, &customItem.Exclusive, &item.ID, &item.Name, &item.Price, &item.Image, &item.Available, &item.Identifier, &item.IsOneOff)

		if err != nil {
			return []model.CustomItem{}, fmt.Errorf("Failed to scan item %w", err)
		}

		_, ok := customItemMap[customItem.ID]

		if !ok {
			customItemMap[customItem.ID] = customItem
		}

		a, ok := itemMap[customItem.ID]
		if ok {
			itemMap[customItem.ID] = append(a, item)
		} else {
			itemMap[customItem.ID] = []model.Item{item}
		}
	}
	res := []model.CustomItem{}
	for _, v := range customItemMap {
		x, s := itemMap[v.ID]

		if !s {
			res = append(res, v)
			continue
		}
		for _, item := range x {
			v.Variants = append(v.Variants, &item)
		}

		res = append(res, v)
	}
	return res, nil
}

func QueryOrders(db *sql.DB) ([]model.Order, error) {
	rows, err := db.Query(`SELECT orders.id, orders.identifier, orders.timestamp, orders.state, orders_items_link.quantity, item.id, item.name, item.price, item.image, item.available, item.identifier, item.oneoff FROM orders
		INNER JOIN orders_items_link ON orders.id=orders_items_link.order_id
		INNER JOIN item ON orders_items_link.item_id=item.id`)
	if err != nil {
		return []model.Order{}, fmt.Errorf("Failed to query db: %w", err)
	}

	orderMap := make(map[int]model.Order)
	orderItemMap := make(map[int][]model.OrderItem)

	for rows.Next() {
		var order model.Order
		var item model.OrderItem

		err := rows.Scan(&order.ID, &order.Identifier, &order.Timestamp, &order.State, &item.Quantity, &item.Item.ID, &item.Item.Name, &item.Item.Price, &item.Item.Image, &item.Item.Available, &item.Item.Identifier, &item.Item.IsOneOff)
		if err != nil {
			return []model.Order{}, fmt.Errorf("Failed to scan row: %w", err)
		}

		_, exists := orderMap[order.ID]
		if !exists {
			orderMap[order.ID] = order
		}

		l, exists := orderItemMap[order.ID]
		if exists {
			orderItemMap[order.ID] = append(l, item)
		} else {
			orderItemMap[order.ID] = []model.OrderItem{item}
		}
	}

	//TODO: SELECT orders.id, orders_custom_items_link.quantity, custom_item.id, custom_item.name, custom_item.exclusive, item.id, item.name, item.price, item.image, item.available, item.identifier, item.oneoff FROM orders INNER JOIN orders_custom_items_link ON orders_custom_items_link.order_id=orders.id INNER JOIN custom_item ON orders_custom_items_link.custom_item_id=custom_item.id INNER JOIN item ON orders_custom_items_link.item_id=item.id
	panic("")
}
