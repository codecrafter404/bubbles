<script lang="ts">
	import { getContextClient, queryStore } from "@urql/svelte";
	import * as gql from "../../generated/graphql";
	import {
		OrderState,
		type NewOrder,
		type Order,
	} from "../../generated/graphql";
	import OrderRenderer from "../Components/OrderRenderer.svelte";
	import ItemSelection from "../Components/ItemSelection.svelte";

	let res = queryStore({
		client: getContextClient(),
		query: gql.GetItemsForStoreDocument,
	});

	let currentOrder: Order = {
		items: [],
		customItems: [],
		id: -2,
		identifier: "",
		state: OrderState.Created,
		timestamp: "",
		total: 0.0,
	};
	async function submitOrder(order: Order): Order {
		let newOrder: NewOrder = {
			total: calculateTotal(order),
			customItems: order.customItems.map((x) => {
				return {
					id: x.customItem.id,
					quantity: x.quantity,
					variants: x.customItem.variants.map(
						(z) => z.id,
					),
				};
			}),
			items: order.items.map((x) => {
				return {
					item: x.item.id,
					quantity: x.quantity,
				};
			}),
		};
	}

	function calculateTotal(order: Order): number {
		let sum = 0;
		order.items.forEach((x) => (sum += x.quantity * x.item.price));
		order.customItems.forEach((x) => {
			x.customItem.variants.forEach((y) => {
				sum += y.price * x.quantity;
			});
		});
		return sum;
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
				{#if currentOrder.id == -2}
					<div
						class="overflow-y-scroll max-h-[100vh] h-[100vh] max-w-[70vw] p-3"
					>
						<ItemSelection
							bind:currentOrder
							customItems={$res.data
								.getCustomItems}
							items={$res.data
								.getItems}
						/>
					</div>
				{:else}
					<div></div>
				{/if}
			{/if}
		</div>
		<div class="w-[30vw]">
			<OrderRenderer
				order={currentOrder}
				total={calculateTotal(currentOrder)}
			/>
		</div>
	</div>
</div>

<style></style>
