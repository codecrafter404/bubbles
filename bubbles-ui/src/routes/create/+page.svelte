<script lang="ts">
	import {
		getContextClient,
		mutationStore,
		queryStore,
		type OperationResultState,
		type OperationResultStore
	} from '@urql/svelte';
	import { graphql } from '../../gql';
	import type { CreateOrderInput, CreateOrderMutation, Item } from '../../gql/graphql';
	import OrderSelectionGrid from './OrderSelectionGrid.svelte';
	import { getNextCustomItem } from './ViewHelper';
	import { KeyboardPressedEvent, ProcessEvent } from '$lib/Keyboard';
	import ItemPricingComponent from '$lib/components/ItemPricingComponent.svelte';

	let graphQLClient = getContextClient();

	let order: CreateOrderInput = $state({
		customItems: [],
		items: []
	});
	const queried = graphql(`
		query QueryItemAndCustomItem {
			items {
				id
				name
				price
				image
				inStock
				notes
			}
			customItems {
				id
				name
				allowOnlyOne
				prev
				variants {
					id
				}
			}
		}
	`);
	const updateOrderGQL = graphql(`
		mutation CreateOrder($order: CreateOrderInput!) {
			createOrder(input: $order) {
				id
				identifier
			}
		}
	`);

	let data = queryStore({
		client: graphQLClient,
		query: queried
	});

	let items = $derived.by(() => {
		if ($data.fetching || $data.error != undefined) {
			return [];
		}
		return $data.data!.items.filter((x) => x.inStock);
	});

	let customItems = $derived.by(() => {
		if ($data.fetching || $data.error != undefined) {
			return [];
		}
		return $data.data!.customItems.map((x) => {
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
	});

	let orderCompleate = $derived.by(() => {
		return getNextCustomItem(customItems, order)[0] == undefined;
	});

	let updatedOrder: OperationResultState<CreateOrderMutation, { order: CreateOrderInput }> | null =
		$state(null);

	function submitOrder() {
		if (order.items.length == 0 && order.customItems.length == 0) return;
		let res = mutationStore({
			client: graphQLClient,
			query: updateOrderGQL,
			variables: {
				order
			}
		});
		res.subscribe((x) => {
			updatedOrder = x;
		});
	}
	let prevOrderId: number | undefined = $state(undefined);

	function handleKeyPressed(event: KeyboardEvent) {
		let res = ProcessEvent(event);
		switch (res) {
			case KeyboardPressedEvent.Confirm:
				// only if error and success
				if (updatedOrder != null && !updatedOrder.fetching) {
					prevOrderId = updatedOrder.data?.createOrder.id ?? prevOrderId;
					order = {
						items: [],
						customItems: []
					};
					updatedOrder = null;
				} else if (orderCompleate) {
					submitOrder();
				}
				break;
			default:
				return;
		}
		event.preventDefault();
	}
</script>

<svelte:body onkeypress={handleKeyPressed} />

<div class="flex flex-row">
	<div
		class={`${updatedOrder != undefined ? 'pointer-events-none select-none' : ''} h-screen w-[70vw]`}
	>
		{#if orderCompleate}
			<p class="font-bold text-green-500">READY</p>
		{:else}
			<p class="font-bold text-red-500">NOT READY</p>
		{/if}

		{#if updatedOrder == null || updatedOrder.fetching}
			{#if $data.fetching}
				<p>Loading...</p>
			{:else if $data.error != undefined}
				<p>Failed to load items & configuration: {JSON.stringify($data.error)}</p>
			{:else}
				<OrderSelectionGrid
					bind:currentOrder={order}
					{items}
					{customItems}
					inputAllowed={updatedOrder == null}
				/>
			{/if}
		{:else if updatedOrder.error}
			<p>Failed to submit order: {JSON.stringify(updatedOrder.error)}</p>
		{:else}
			<p class="text-2xl font-bold">#{updatedOrder.data!.createOrder.identifier}</p>
		{/if}
	</div>
	<div class="grow bg-blue-300">
		<ItemPricingComponent {customItems} {items} orderInput={order} />
	</div>
</div>
