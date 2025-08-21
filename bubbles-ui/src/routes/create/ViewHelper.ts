import type { CreateOrderInput, CustomItem, Item } from "../../gql/graphql";

export interface SelectionItemProps {
	id: number;
	name: string;
	notes: string;
	price: number;
	image: string;
	count: number;
	isCustomItem: boolean;
	selected: boolean;
}

export function getNextItems(
	items_input: Array<Item> | undefined,
	customItem_input: Array<any> | undefined,
	currentOrder: CreateOrderInput
): Array<SelectionItemProps> {
	if (items_input == undefined || customItem_input == undefined) {
		return [];
	}

	let items: Array<Item> = items_input.filter((x) => x.inStock);
	let customItems: Array<CustomItem> = customItem_input.map((x) => {
		let variants: Array<Item> = [];

		x.variants.forEach((i: any) => {
			let found = items.find((x) => x.id == i.id);
			if (found != undefined) {
				variants.push(found!);
			}
		});
		return {
			...x,
			variants: variants
		};
	});

	// top level menu, when currently no customitem is incompleate
	let [nextCustomItems, _] = getNextCustomItem(customItems, currentOrder);

	if (nextCustomItems == undefined) {
		let res_customItems: Array<SelectionItemProps> = [];

		customItems.forEach((ci) => {
			if (ci.prev == null) {
				res_customItems.push({
					id: ci.id,
					name: ci.name,
					notes: '',
					price: calculateMinimalPrice(ci, customItems),
					image: '', //TODO: add default image,
					isCustomItem: true,
					count: 0,
					selected: false
				});
			}
		});

		res_customItems = res_customItems.sort((a, b) => a.name.localeCompare(b.name));

		let res_Items: Array<SelectionItemProps> = [];

		items.forEach((item) => {
			let isVariant = customItems.some((x) => x.variants.find((y) => y.id == item.id));
			if (!isVariant && item.inStock) {
				let count = currentOrder!.items.find((x) => x.itemID == item.id)?.quantity ?? 0;
				res_Items.push({
					count: count,
					isCustomItem: false,
					selected: false,
					...item
				});
			}
		});

		res_Items = res_Items.sort((a, b) => a.name.localeCompare(b.name));

		return [...res_customItems, ...res_Items];
	} else {
		return nextCustomItems!.variants.map((x) => {
			let count =
				currentOrder!.customItems
					.find(
						(x) =>
							x.selectedCustomItems.find((y) => y.customItemID == nextCustomItems!.id) !=
							undefined
					)
					?.selectedCustomItems.find(
						(u) => u.selectedVariantIDs.find((y) => y == x.id) != undefined
					)
					?.selectedVariantIDs.filter((z) => z == x.id).length ?? 0;
			return {
				count,
				isCustomItem: false,
				selected: false,
				...x
			};
		});
	}
}

function calculateMinimalPrice(customItem: CustomItem, customItems: Array<CustomItem>): number {
	if (customItem.prev != null) {
		return NaN;
	}

	let total = 0;
	let next: CustomItem | undefined = customItem;

	while (next != undefined) {
		if (next.variants.length >= 1) {
			total += Math.min(...next.variants.map((x) => x.price));
		}
		next = customItems.find((x) => x.prev == next!.id);
	}

	return total;
}

/// return undefined when no customitem is found
export function getNextCustomItem(
	customItems: Array<CustomItem>,
	order: CreateOrderInput
): [CustomItem | undefined, number] {
	let res: CustomItem | undefined;
	let idx = -1;

	order.customItems.forEach((oci, i) => {
		if (res != undefined) {
			return;
		}

		let to_check =
			oci.selectedCustomItems.length == 0
				? customItems.find((x) => x.id == oci.masterCustomItemID)!
				: customItems.find(
					(x) =>
						x.id == oci.selectedCustomItems[oci.selectedCustomItems.length - 1].customItemID
				)!;


		let sel_variants_len = oci.selectedCustomItems.find((x) => x.customItemID == to_check.id)?.selectedVariantIDs.length ?? 0;

		if (to_check.variants.length >= 1 && sel_variants_len == 0) {
			res = to_check;
			idx = i;
			return;
		} // else: if 0 variants are to select

		let next = customItems.find((x) => x.prev != undefined && x.prev! == to_check.id);
		if (next != undefined) {
			res = next;
			idx = i;
			return;
		}
	});

	return [res, idx];
}
