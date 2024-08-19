package utils

import (
	"database/sql"
	"fmt"
	"slices"
	"strconv"
	"time"

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
	total REAL NOT NULL,
	timestamp INTEGER NOT NULL,
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
	defer rows.Close()
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
	defer rows.Close()
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
func QueryOrdersLimited(db *sql.DB, state *model.OrderState, id *int, limit *int, skip *int, sortAsc *bool) ([]*model.Order, error) {
	opts := OrderQueryOptions{
		SortAsc: sortAsc,
		State:   state,
	}
	if id != nil {
		opts.FilterIds = &[]int{*id}
	}
	orders, err := QueryOrders(db, &opts)
	if err != nil {
		return []*model.Order{}, fmt.Errorf("Failed to query db: %w", err)
	}

	skipped := 0
	collected := 0

	res := []*model.Order{}
	for _, o := range orders {
		if skip != nil && skipped < *skip {
			skipped++
			continue
		}

		if limit != nil && collected >= *limit {
			break
		}

		res = append(res, &o)
		collected++
	}
	return res, nil
}

type OrderQueryOptions struct {
	// sort by time asc
	SortAsc   *bool
	FilterIds *[]int
	State     *model.OrderState
}

func queryOrdersTransaction(tx *sql.Tx, options *OrderQueryOptions) ([]model.Order, error) {
	orderQuery := `SELECT orders.id, orders.total, orders.identifier, orders.timestamp, orders.state, orders_items_link.quantity, item.id, item.name, item.price, item.image, item.available, item.identifier, item.oneoff FROM orders
		LEFT JOIN orders_items_link ON orders.id=orders_items_link.order_id
		LEFT JOIN item ON orders_items_link.item_id=item.id`

	orderQueryArgs := []any{}

	if options != nil {
		orderQuery, orderQueryArgs = PrepareQueryOrdersWithOptions(orderQuery, orderQueryArgs, *options)
	}

	rows, err := tx.Query(orderQuery, orderQueryArgs...)
	defer rows.Close()

	if err != nil {
		return []model.Order{}, fmt.Errorf("Failed to query db: %w", err)
	}

	orderMap := make(map[int]model.Order)
	orderItemMap := make(map[int][]model.OrderItem)

	for rows.Next() {
		var order model.Order
		var item model.OrderItem
		scanItem := struct {
			iQuantity   *int
			iId         *int
			iName       *string
			iPrice      *float64
			iImage      *string
			iAvailable  *bool
			iIdentifier *string
			iOneOff     *bool
		}{}
		err := rows.Scan(&order.ID, &order.Total, &order.Identifier, &order.Timestamp, &order.State, &scanItem.iQuantity, &scanItem.iId, &scanItem.iName, &scanItem.iPrice, &scanItem.iImage, &scanItem.iAvailable, &scanItem.iIdentifier, &scanItem.iOneOff)

		if err != nil {
			return []model.Order{}, fmt.Errorf("Failed to scan row1: %w", err)
		}

		_, exists := orderMap[order.ID]
		if !exists {
			orderMap[order.ID] = order
		}

		if scanItem.iQuantity == nil {
			continue
		}

		item = model.OrderItem{
			Quantity: *scanItem.iQuantity,
			Item: &model.Item{
				ID:         *scanItem.iId,
				Name:       *scanItem.iName,
				Price:      *scanItem.iPrice,
				Image:      *scanItem.iImage,
				Available:  *scanItem.iAvailable,
				Identifier: *scanItem.iIdentifier,
				IsOneOff:   *scanItem.iOneOff,
			},
		}

		l, exists := orderItemMap[order.ID]
		if exists {
			orderItemMap[order.ID] = append(l, item)
		} else {
			orderItemMap[order.ID] = []model.OrderItem{item}
		}
	}

	orderCustomItemQuery := `SELECT orders.id, orders_custom_items_link.quantity, custom_item.id, custom_item.name, custom_item.exclusive, item.id, item.name, item.price, item.image, item.available, item.identifier, item.oneoff
		FROM orders
		INNER JOIN orders_custom_items_link ON orders_custom_items_link.order_id=orders.id
		LEFT JOIN custom_item ON orders_custom_items_link.custom_item_id=custom_item.id
		LEFT JOIN item ON orders_custom_items_link.item_id=item.id`
	orderCustomItemQueryArgs := []any{}

	if options != nil {
		orderCustomItemQuery, orderCustomItemQueryArgs = PrepareQueryOrdersWithOptions(orderCustomItemQuery, orderCustomItemQueryArgs, *options)
	}

	orderCustomItem, err := tx.Query(orderCustomItemQuery, orderCustomItemQueryArgs...)
	defer orderCustomItem.Close()

	if err != nil {
		return []model.Order{}, fmt.Errorf("Failed to query orderItems: %w", err)
	}

	customItemMap := make(map[int][]model.OrderCustomItem)
	for orderCustomItem.Next() {
		var order int
		var customItem model.OrderCustomItem
		customItem.CustomItem = &model.CustomItem{}
		var item model.Item

		err := orderCustomItem.Scan(&order, &customItem.Quantity, &customItem.CustomItem.ID, &customItem.CustomItem.Name, &customItem.CustomItem.Exclusive, &item.ID, &item.Name, &item.Price, &item.Image, &item.Available, &item.Identifier, &item.IsOneOff)
		if err != nil {
			return []model.Order{}, fmt.Errorf("Failed to scan row: %w", err)
		}

		x, exists := customItemMap[order]
		if !exists {
			customItem.CustomItem.Variants = append(customItem.CustomItem.Variants, &item)
			customItemMap[order] = []model.OrderCustomItem{customItem}
			continue
		}

		found := false

		res := []model.OrderCustomItem{}

		for _, i := range x {
			if i.CustomItem.ID == customItem.CustomItem.ID && i.Quantity == customItem.Quantity {
				found = true
				i.CustomItem.Variants = append(i.CustomItem.Variants, &item)
			}
			res = append(res, i)
		}

		if !found {
			customItem.CustomItem.Variants = append(customItem.CustomItem.Variants, &item)
			res = append(res, customItem)
		}

		customItemMap[order] = res
	}
	res := []model.Order{}
	for orderId, order := range orderMap {
		if items, exists := orderItemMap[orderId]; exists {
			resItems := []*model.OrderItem{}
			for _, x := range items {
				resItems = append(resItems, &x)
			}
			order.Items = resItems
		}
		if cItems, exists := customItemMap[orderId]; exists {
			resItems := []*model.OrderCustomItem{}
			for _, x := range cItems {
				resItems = append(resItems, &x)
			}
			order.CustomItems = resItems
		}
		res = append(res, order)
	}
	if options != nil && options.SortAsc != nil {
		slices.SortFunc(res, func(a, b model.Order) int {

			aTime := time.UnixMicro(a.Timestamp)
			bTime := time.UnixMicro(b.Timestamp)

			res := aTime.Compare(bTime)
			if !*options.SortAsc {
				res *= -1
			}

			return res
		})
	}
	return res, nil
}
func QueryOrders(db *sql.DB, options *OrderQueryOptions) ([]model.Order, error) {
	tx, err := db.Begin()
	defer tx.Rollback()

	if err != nil {
		return []model.Order{}, fmt.Errorf("Failed to begin transaction: %w", err)
	}
	res, err := queryOrdersTransaction(tx, options)
	if err != nil {
		return []model.Order{}, fmt.Errorf("Failed to query orders: %w", err)
	}
	tErr := tx.Commit()
	if tErr != nil {
		return []model.Order{}, fmt.Errorf("Failed to commit transaction: %w", tErr)
	}
	return res, nil
}
func PrepareQueryOrdersWithOptions(query string, args []any, options OrderQueryOptions) (string, []any) {
	if options.FilterIds != nil || options.State != nil {
		query += " WHERE"
	}
	first := true
	if options.FilterIds != nil {
		for _, x := range *options.FilterIds {
			query += " "
			if !first {
				query += "&& "
			} else {
				first = false
			}
			query += "orders.id = ?"
			args = append(args, x)
		}
	}
	if options.State != nil {
		query += " "
		if !first {
			query += "&& "
		} else {
			first = false
		}

		query += "orders.state = ?"
		args = append(args, options.State)
	}
	return query, args
}
func GetIdentifier(db *sql.DB, max int) (int, error) {
	res, err := db.Query("SELECT identifier FROM orders ORDER BY timestamp DESC LIMIT 1")
	defer res.Close()
	if err != nil {
		return 0, fmt.Errorf("Failed to getIdentifier: %w", err)
	}
	has := res.Next()
	defer res.Close()
	if !has {
		return 1, nil
	}
	var snum string
	scanErr := res.Scan(&snum)
	if scanErr != nil {
		return 0, fmt.Errorf("Failed to scan identifier: %w", err)
	}
	num, err := strconv.Atoi(snum)
	if err != nil {
		return 0, fmt.Errorf("Invalid db state, expected number as identifier, got %s: %w", snum, err)
	}

	num++
	if num > max {
		num = 1
	}

	return num, nil
}
func GetStats(db *sql.DB) (model.Statistics, error) {
	stateRow, err := db.Query("SELECT SUM(total), COUNT(*) FROM orders WHERE state = ?", model.OrderStateCompleated)
	defer stateRow.Close()
	if err != nil {
		return model.Statistics{}, fmt.Errorf("Failed to query stateRow: %w", err)
	}
	var stats model.Statistics
	var totalEarned *float64
	if !stateRow.Next() {
		return model.Statistics{TotalEarned: 0, TotalOrders: 0, TotalOrdersCompleated: 0}, nil
	}
	eErr := stateRow.Scan(&totalEarned, &stats.TotalOrdersCompleated)
	if eErr != nil {
		return model.Statistics{}, fmt.Errorf("Failed to scan stateRows: %w", eErr)
	}
	if totalEarned == nil {
		stats.TotalEarned = 0
	} else {
		stats.TotalEarned = *totalEarned
	}

	totalRows, err := db.Query("SELECT COUNT(*) FROM orders")
	defer totalRows.Close()

	if !totalRows.Next() {
		return model.Statistics{TotalEarned: 0, TotalOrders: 0, TotalOrdersCompleated: 0}, nil
	}

	if err != nil {
		return model.Statistics{}, fmt.Errorf("Failed to query totalRows: %w", err)
	}

	tErr := totalRows.Scan(&stats.TotalOrders)
	if tErr != nil {
		return model.Statistics{}, fmt.Errorf("Failed to scan totalRows: %w", tErr)
	}
	return stats, nil
}

func GetNextOrder(db *sql.DB) (*model.Order, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("Failed to lock db: %w", err)
	}
	defer tx.Rollback()
	res, err := tx.Query("SELECT id FROM orders WHERE state = ? ORDER BY timestamp ASC LIMIT 1", model.OrderStateCreated)
	if err != nil {
		return nil, fmt.Errorf("Failed to query orders")
	}
	defer res.Close()

	var id *int

	if res.Next() {
		err := res.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan id: %w", err)
		}
	}
	if id == nil {
		tErr := tx.Commit()
		if tErr != nil {
			return nil, fmt.Errorf("Failed to commit transaction: %w", tErr)
		}
		return nil, nil // no order available
	}

	_, eErr := tx.Exec("UPDATE orders SET state = ? WHERE id = ?", model.OrderStatePending, id)
	if eErr != nil {
		return nil, fmt.Errorf("Failed to update order state: %w", eErr)
	}

	opts := OrderQueryOptions{
		FilterIds: &[]int{*id},
	}

	order, err := queryOrdersTransaction(tx, &opts)
	if err != nil {
		return nil, fmt.Errorf("Failed to query order: %w", err)
	}

	if len(order) != 1 {
		return nil, fmt.Errorf("Failed to query order: order has been deleted?")
	}

	tErr := tx.Commit()
	if tErr != nil {
		return nil, fmt.Errorf("Failed to commit transaction: %w", tErr)
	}

	return &order[0], nil
}
