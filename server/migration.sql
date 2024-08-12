CREATE TABLE IF NOT EXISTS item(
	id INTEGER NOT NULL PRIMARY KEY,
	name TEXT NOT NULL,
	price REAL NOT NULL,
	image TEXT NOT NULL,
	available INTEGER NOT NULL,
	identifier TEXT NOT NULL
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
	item_id INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS orders_custom_items_link(
	order_id INTEGER NOT NULL,
	custom_item_id INTEGER NOT NULL
);
