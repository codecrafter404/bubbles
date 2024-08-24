<script lang="ts">
	import { getContextClient, queryStore } from "@urql/svelte";
	import NavBar from "../Components/NavBar.svelte";
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
