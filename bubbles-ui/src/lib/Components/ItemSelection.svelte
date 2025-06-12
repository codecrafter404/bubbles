<script lang="ts">
	import type { CustomItem, Item, Order } from "../../generated/graphql";
	import ItemSelectionGrid from "./ItemSelectionGrid.svelte";
	import type { SelectableItem } from "../Types/SelectableItem";

	export let columns = 5;

	export let customItems: Array<CustomItem>;
	export let items: Array<Item>;

	export let currentOrder: Order;

	customItems = mapCustomItemVariants(customItems, items);

	// maps the customitem variants ids to thier actual objects
	function mapCustomItemVariants(
		customItems: CustomItem[],
		items: Item[],
	): CustomItem[] {
		return customItems.map((x) => {
			x.variants = x.variants.map(
				(y) => items.find((z) => z.id == y.id)!,
			);
			return x;
		});
	}

	function addCustomItemToOrder(item: CustomItem) {
		let customizing = currentlyCustomizing();
		if (customizing != null)
			throw Error("This should never happen!");
		currentOrder.customItems = [
			...currentOrder.customItems,
			{
				customItem: {
					...item,
					variants: [],
				},
				quantity: 1,
			},
		];
	}

	function removeCustomItemFromOrder(item: CustomItem) {
		let customizing = currentlyCustomizing();
		if (customizing != null)
			throw Error("This should never happen!");
		let first = true;
		currentOrder.customItems = currentOrder.customItems
			.reverse()
			.map((x) => {
				if (x.customItem.id == item.id) {
					if (first) {
						x.quantity -= 1;
					}
				}
				return x;
			})
			.filter((x) => x.quantity > 0);
	}

	function addItemToOrder(item: Item) {
		let exists = currentOrder.items.find(
			(x) => x.item.id == item.id,
		);
		let customizing = currentlyCustomizing();
		if (customizing != null) {
			currentOrder.customItems = currentOrder.customItems.map(
				(x) => {
					let itemChain = getCustomItemChainById(
						x.customItem.id,
					);
					for (let chainItem of itemChain) {
						if (
							chainItem.variants.find(
								(y) =>
									x.customItem.variants.includes(
										y,
									),
							) == null
						) {
							x.customItem.variants =
								[
									...x
										.customItem
										.variants,
									item,
								];
							return x;
						}
					}
					return x;
				},
			);
			return;
		}
		if (exists != null) {
			currentOrder.items = currentOrder.items.map((x) => {
				if (x.item.id == item.id) {
					x.quantity += 1;
				}
				return x;
			});
		} else {
			currentOrder.items = [
				...currentOrder.items,
				{
					item: item,
					quantity: 1,
				},
			];
		}
	}
	function removeItemFromOrder(item: Item) {
		currentOrder.items = currentOrder.items
			.map((x) => {
				if (item.id == x.item.id) x.quantity -= 1;
				console.log(x);
				return x;
			})
			.filter((x) => x.quantity > 0);
	}

	function addToOrder(event: CustomEvent<SelectableItem>) {
		switch (event.detail.data["type"]) {
			case "item":
				let item = event.detail.data["item"] as Item;
				item = items.find((x) => x.id == item.id)!;
				addItemToOrder(item);
				break;
			case "CustomItem":
				let customItem = event.detail.data[
					"item"
				] as CustomItem;
				customItem = customItems.find(
					(x) => x.id == customItem.id,
				)!;
				console.log(
					customItem,
					"this item should not be empty",
					customItems,
				);
				addCustomItemToOrder(customItem);
				break;
			default:
				throw Error("UNIMPLEMENTED case");
		}
	}
	function removeFromOrder(event: CustomEvent<SelectableItem>) {
		switch (event.detail.data["type"]) {
			case "item":
				let item = event.detail.data["item"] as Item;
				removeItemFromOrder(item);
				break;
			case "CustomItem":
				let customItem = event.detail.data[
					"item"
				] as CustomItem;
				removeCustomItemFromOrder(customItem);
				break;
			default:
				throw Error("UNIMPLEMENTED case");
		}
	}

	// will get stuck if ther's a loop
	function getDependencies(
		root: CustomItem,
		items: Array<CustomItem>,
	): Array<CustomItem> {
		console.log("root", root, items);
		let res: Array<CustomItem> = [];
		let current = root;
		while (current != null) {
			let x = items.find((x) => x.dependsOn == current.id);
			if (x != null) res.push(x);
			current = x!;
		}
		return res;
	}

	function getCustomItemChainById(id: number): CustomItem[] {
		let customItem = customItems.find((x) => x.id == id);
		if (customItem == null) return [];
		let itemChain = [
			customItem,
			...getDependencies(customItem, customItems),
		];
		return itemChain;
	}

	/// returns true, if we're currently customizing a CustomItem
	function currentlyCustomizing(): [CustomItem, CustomItem] | null {
		for (let customOrderItem of currentOrder.customItems) {
			let itemChain = getCustomItemChainById(
				customOrderItem.customItem.id,
			);
			console.log("chain", itemChain);
			for (let chainItem of itemChain) {
				if (
					chainItem.variants.find((x) =>
						customOrderItem.customItem.variants.includes(
							x,
						),
					) == null
				) {
					return [
						chainItem,
						customOrderItem.customItem,
					];
				}
			}
		}
		return null;
	}

	function itemToSelectableItem(
		item: Item,
		order: Order,
	): SelectableItem {
		return {
			name: item.name,
			image: item.image,
			price: `${item.price}€`,
			quantity: order.items
				.filter((x) => x.item.id == item.id)
				.reduce((a, b) => a + b.quantity, 0),
			data: { type: "item", item: item },
		};
	}
	function customItemToSelectableItem(
		item: CustomItem,
		order: Order,
	): SelectableItem {
		let chain = getCustomItemChainById(item.id);
		let min = chain
			.map(
				(x) =>
					x.variants.sort(
						(a, b) => a.price - b.price,
					)[0].price,
			)
			.filter((x) => x != null)
			.reduce((s, a) => s + a);
		let max = chain
			.map(
				(x) =>
					x.variants.sort(
						(a, b) => b.price - a.price,
					)[0].price,
			)
			.filter((x) => x != null)
			.reduce((s, a) => s + a);

		return {
			name: item.name,
			image: "",
			quantity: order.customItems
				.filter((x) => x.customItem.id == item.id)
				.map((x) => x.quantity)
				.reduce((a, b) => a + b, 0),
			price: `${min}€ - ${max}€`,
			data: { type: "CustomItem", item: item },
		};
	}

	function getSelectableItems(
		customItems: CustomItem[],
		items: Item[],
		order: Order,
	): SelectableItem[] {
		let customizing = currentlyCustomizing();
		console.log(customizing);
		if (customizing != null)
			return customizing[0].variants.map((x) =>
				itemToSelectableItem(x, order),
			);
		let res = items
			.filter((x) => x.isOneOff && x.available)
			.map((x) => itemToSelectableItem(x, order));
		res = [
			...customItems //TODO: make configurable (if customitems before or after regular ones)
				.filter((x) => x.dependsOn == null)
				.map((x) =>
					customItemToSelectableItem(x, order),
				),
			...res,
		];
		return res;
	}
</script>

<div>
	<!-- TODO: impl add/remove -->
	<ItemSelectionGrid
		items={getSelectableItems(customItems, items, currentOrder)}
		gridCols={columns}
		on:add={addToOrder}
		on:remove={removeFromOrder}
	/>
</div>
