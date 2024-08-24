<script lang="ts">
	import { onDestroy, onMount } from "svelte";
	import type { CustomItem, Item, Order } from "../../generated/graphql";
	import CreateSelectionItem from "./CreateSelectionItem.svelte";

	export let columns = 5;

	export let customItems: Array<CustomItem>;
	export let items: Array<Item>;

	export let currentOrder: Order;

	let current = 0;
	let currentIsCustom = false;

	// init current
	select_nth(0);

	function handle_key_down(e: KeyboardEvent) {
		console.log(e);
		let selected = get_current_selected();
		if (e.code == "Numpad8") {
			// up
			let res = selected;
			res -= columns;
			if (res < 0) {
				return;
			}
			select_nth(res);
		}
		if (e.code == "Numpad6") {
			// right
			let res = selected;
			let [cCustomFiltered, cItemFiltered] =
				get_selectable_nodes();
			let len = cCustomFiltered.length + cItemFiltered.length;
			res++;
			if (res >= len) {
				return;
			}
			select_nth(res);
		}
		if (e.code == "Numpad2") {
			// down
			let res = selected;
			let [cCustomFiltered, cItemFiltered] =
				get_selectable_nodes();
			let len = cCustomFiltered.length + cItemFiltered.length;

			res += columns;

			if (res >= len) {
				return;
			}
			select_nth(res);
			// down
		}
		if (e.code == "Numpad4") {
			// left
			let res = selected;
			res--;
			if (res < 0) {
				return;
			}
			select_nth(res);
		}
		if (e.code == "NumpadEnter") {
			// add to order
			add_current_to_order();
		}
		if (e.code == "NumpadSubtract") {
			// remove from order
			remove_current_from_order();
		}
	}

	function get_selectable_nodes(): [Array<CustomItem>, Array<Item>] {
		let cCustomFiltered = customItems.filter(
			(x) => x.dependsOn != null,
		);
		let cItemFiltered = items.filter(
			(x) => x.available && x.isOneOff,
		);
		return [cCustomFiltered, cItemFiltered];
	}
	function get_current_selected(): number {
		let [cCustomFiltered, cItemFiltered] = get_selectable_nodes();
		let idx = -1;
		if (currentIsCustom) {
			idx = cCustomFiltered.findIndex((x) => x.id == current);
		} else {
			idx =
				cItemFiltered.findIndex(
					(x) => x.id == current,
				) +
				customItems.length -
				1;
		}
		return idx;
	}
	function add_current_to_order() {
		let [cCustomFiltered, cItemFiltered] = get_selectable_nodes();
		if (currentIsCustom) {
			let c = customItems.find((x) => x.id == current)!;
			currentOrder.customItems = [
				...currentOrder.customItems,
				{
					customItem: c,
					quantity: 1,
				},
			];
			console.log("Not yet implemented :(");
		} else {
			let c = cItemFiltered.find((x) => x.id == current)!;
			add_more_item_to_order(c);
		}
		console.log(currentOrder);
	}
	function remove_current_from_order() {
		let [cCustomFiltered, cItemFiltered] = get_selectable_nodes();
		if (currentIsCustom) {
			let c = customItems.find((x) => x.id == current);
			console.log("Not yet implemented :(");
		} else {
			let c = cItemFiltered.find((x) => x.id == current)!;
			let newOrderItems = currentOrder.items.map((x) => {
				if (x.item.id == c.id) {
					x.quantity--;
				}
				return x;
			});
			newOrderItems = newOrderItems.filter(
				(x) => x.quantity >= 1,
			);
			currentOrder.items = newOrderItems;
		}
	}
	function add_more_item_to_order(item: Item) {
		let exists = currentOrder.items.find(
			(x) => x.item.id == item.id,
		);
		if (exists) {
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
	function select_nth(n: number) {
		let [cCustomFiltered, cItemFiltered] = get_selectable_nodes();
		let len = cCustomFiltered.length + cItemFiltered.length;

		if (n > len - 1) {
			console.info("Maybee bug?", len, n);
			return;
		}

		if (n < cCustomFiltered.length) {
			// custom item range
			current = cCustomFiltered[n].id;
			currentIsCustom = true;
		} else {
			// custom item range
			current = cItemFiltered[n - cCustomFiltered.length].id;
			currentIsCustom = false;
		}
	}
	function mapCustomItems(
		customItems: Array<CustomItem>,
		item: Array<Item>,
	): Array<[CustomItem, number, number]> {
		let res: Array<[CustomItem, number, number]> = customItems.map(
			(x) => {
				x.variants = x.variants.map(
					(y) => item.find((z) => z.id == y.id)!,
				);
				return [x, 0, 0];
			},
		);
		res = res.map((x) => {
			let deps = getDependencies(
				x[0],
				res.map((x) => x[0]),
			);
			deps.push(x[0]);
			let minPrice = 0;
			let maxPrice = 0;
			deps.forEach((y) => {
				let sorted = y.variants.sort(
					(a, b) => a.price - b.price,
				);
				if (sorted.length == 0) {
					return;
				}
				minPrice += sorted[0].price;
				if (y.exclusive) {
					maxPrice +=
						sorted[sorted.length - 1].price;
				} else {
					sorted.forEach(
						(z) => (maxPrice += z.price),
					);
				}
			});
			return [x[0], minPrice, maxPrice];
		});
		return res;
	}
	// will get stuck if ther's a loop
	function getDependencies(
		root: CustomItem,
		items: Array<CustomItem>,
	): Array<CustomItem> {
		let res: Array<CustomItem> = [];
		let current = root;
		while (current.dependsOn) {
			let x = items.find((x) => x.id == current.dependsOn)!;
			res.push(x);
			current = x;
		}
		return res;
	}

	onMount(() => {
		document.addEventListener("keydown", handle_key_down);
	});
	onDestroy(() => {
		document.removeEventListener("keydown", handle_key_down);
	});
</script>

<div
	class="grid gap-3 items-center"
	style={`grid-template-columns: repeat(${columns}, minmax(0, 1fr))`}
>
	{#each mapCustomItems(customItems, items) as item}
		{#if item[0].dependsOn}
			<CreateSelectionItem
				item={null}
				customItem={item[0]}
				customItemPriceEnd={item[1]}
				customItemPriceStart={item[2]}
				hovered={item[0].id == current &&
					currentIsCustom}
				quantity={(() => {
					let x = currentOrder.customItems.find(
						(x) =>
							x.customItem.id ==
							item[0].id,
					);
					if (x) {
						return x.quantity;
					}
					return 0;
				})()}
			/>
		{/if}
	{/each}
	{#each items as item}
		{#if item.isOneOff && item.available}
			<CreateSelectionItem
				{item}
				customItem={null}
				customItemPriceEnd={null}
				customItemPriceStart={null}
				hovered={item.id == current && !currentIsCustom}
				quantity={(() => {
					let x = currentOrder.items.find(
						(x) => x.item.id == item.id,
					);
					if (x) {
						return x.quantity;
					}
					return 0;
				})()}
			/>
		{/if}
	{/each}
</div>
