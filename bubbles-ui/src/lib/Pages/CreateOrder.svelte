<script lang="ts">
	import {
		getContextClient,
		mutationStore,
		queryStore,
		type OperationResultStore,
	} from "@urql/svelte";
	import * as gql from "../../generated/graphql";
	import {
		OrderState,
		type NewOrder,
		type Order,
	} from "../../generated/graphql";
	import OrderRenderer from "../Components/OrderRenderer.svelte";
	import ItemSelection from "../Components/ItemSelection.svelte";
	import { onDestroy, onMount } from "svelte";
	import { confirm } from "../utils/NotificationUtils";

	let client = getContextClient();
	let res = queryStore({
		client: client,
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

	let submitedOrder: OperationResultStore<
		gql.SubmitOrderMutation,
		{
			order: gql.NewOrder;
		}
	> | null;
	let lastCreatedId: number | null;

	function submitOrder(order: Order) {
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

		submitedOrder = mutationStore({
			client: client,
			query: gql.SubmitOrderDocument,
			variables: {
				order: newOrder,
			},
		});
		submitedOrder.subscribe((x) => {
			if (x.data != null) {
				lastCreatedId = x.data.createOrder.id;
			}
		});
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

	function handle_key_press(e: KeyboardEvent) {
		if (e.code == "Numpad0") {
			if (submitedOrder == null) {
				if (
					currentOrder.customItems.length >= 1 ||
					currentOrder.items.length >= 1
				) {
					submitOrder(currentOrder);
				}
			} else {
				cleanUp();
			}
		}
		if (e.code == "NumpadDivide") {
			if (lastCreatedId) {
				undoOrder(lastCreatedId);
				lastCreatedId = null;
			}
		}
	}
	function cleanUp() {
		submitedOrder = null;
		currentOrder = {
			items: [],
			customItems: [],
			id: -2,
			identifier: "",
			state: OrderState.Created,
			timestamp: "",
			total: 0.0,
		};
	}
	function undoOrder(id: number) {
		if (
			!confirm(
				"Are you sure, you want to undo the latest order?",
			)
		) {
			return;
		}
		mutationStore({
			client: client,
			query: gql.UndoOrderDocument,
			variables: {
				order: id,
			},
		});
		alert("Order cancelled!");
		cleanUp();
	}

	onMount(() => {
		document.addEventListener("keydown", handle_key_press);
	});
	onDestroy(() => {
		document.removeEventListener("keydown", handle_key_press);
	});
</script>

<div class="flex flex-col min-h-[100vh] max-h-[100vh] h-[100vh]">
	<div class="flex-grow flex flex-col 2xl:flex-row">
		<div class="flex-grow max-h-full">
			{#if $res.fetching}
				<p>Loading</p>
			{:else if $res.error}
				<p>Failed to load: {$res.error.message}</p>
			{:else if $res.data}
				{#if submitedOrder == null}
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
					<div
						class="text-center h-full justify-center flex flex-col"
					>
						{#if $submitedOrder.fetching}
							<p>
								Loading order
								information
							</p>
						{:else if $submitedOrder.error}
							<p>
								Failed to load
								order
								information:
								{$submitedOrder
									.error
									.message}
							</p>
						{:else if $submitedOrder.data}
							<h1 class="text-4xl">
								<span
									class="text-natural-500"
									>Order
								</span><span
									class="text-accent-500 font-bold"
									>#{$submitedOrder
										?.data
										.createOrder
										.identifier}</span
								>
							</h1>
						{/if}
					</div>
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
