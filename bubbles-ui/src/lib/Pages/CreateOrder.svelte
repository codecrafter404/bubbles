<script lang="ts">
	import { getContextClient, queryStore } from "@urql/svelte";
	import * as gql from "../../generated/graphql";
	import CreateSelectionItem from "../Components/CreateSelectionItem.svelte";
	import type { CustomItem, Item } from "../../generated/graphql";

	let res = queryStore({
		client: getContextClient(),
		query: gql.GetItemsForStoreDocument,
	});
	let current = 1;
	let currentName = "Item 2";
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
		console.log(res);
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
</script>

<div class="flex flex-col min-h-[100vh] max-h-[100vh] h-[100vh]">
	<div class="flex-grow flex flex-col 2xl:flex-row">
		<div class="flex-grow max-h-full">
			{#if $res.fetching}
				<p>Loading</p>
			{:else if $res.error}
				<p>Failed to load: {$res.error.message}</p>
			{:else if $res.data}
				<div
					class="overflow-y-scroll max-h-[100vh] h-[100vh] max-w-[70vw] p-3"
				>
					<div
						class="grid grid-cols-5 gap-3 items-center"
					>
						{#each $res.data.getItems as item}
							{#if item.isOneOff && item.available}
								<CreateSelectionItem
									{item}
									customItem={null}
									customItemPriceEnd={null}
									customItemPriceStart={null}
									hovered={item.id ==
										current &&
										item.name ==
											currentName}
								/>
							{/if}
						{/each}
						{#each mapCustomItems($res.data.getCustomItems, $res.data.getItems) as item}
							{#if item[0].dependsOn}
								<CreateSelectionItem
									item={null}
									customItem={item[0]}
									customItemPriceEnd={item[1]}
									customItemPriceStart={item[2]}
									hovered={item[0]
										.id ==
										current &&
										item[0]
											.name ==
											currentName}
								/>
							{/if}
						{/each}
					</div>
				</div>
			{/if}
		</div>
		<div class="w-[30vw] bg-orange-200">
			OrderList (Shopping Cart)
		</div>
	</div>
</div>

<style></style>
