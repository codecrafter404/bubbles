<script lang="ts">
	import { onDestroy, onMount } from "svelte";
	import {
		getContextClient,
		subscriptionStore,
		mutationStore,
	} from "@urql/svelte";
	import * as gql from "../../generated/graphql";
	import OrderRenderer from "../Components/OrderRenderer.svelte";

	let client = getContextClient();
	let currentOrder = subscriptionStore({
		client: client,
		query: gql.NextOrderDocument,
	});

	let currentOrderId: number | null = null;
	let lastOrderId: number | null = null;

	currentOrder.subscribe((x) => {
		if (x.data != null) {
			if (x.data.nextOrder != null) {
				currentOrderId = x.data.nextOrder.id;
			}
		}
	});

	function handle_keypress(e: KeyboardEvent) {
		if (e.code == "Numpad0") {
			// compleate order
			compleateOrder();
		}
		if (e.code == "NumpadDivide") {
			undoCompleateOrder();
		}
	}

	function compleateOrder() {
		if (currentOrderId == null) {
			return;
		}
		mutationStore({
			client: client,
			query: gql.CompleateOrderDocument,
			variables: {
				order: currentOrderId,
			},
		});
		lastOrderId = currentOrderId;
	}
	function undoCompleateOrder() {
		if (lastOrderId == null) return;

		mutationStore({
			client: client,
			query: gql.UndoCompleateOrderDocument,
			variables: {
				order: lastOrderId,
			},
		});
		lastOrderId = null;
	}
	onMount(() => {
		document.addEventListener("keydown", handle_keypress);
	});
	onDestroy(() => {
		document.removeEventListener("keydown", handle_keypress);
	});
</script>

<div class="items-center flex justify-center">
	<div class="w-[30vw] h-[100vh]">
		{#if $currentOrder.error}
			<p class="text-center">
				Failed to load next order: {$currentOrder.error
					.message}
			</p>
		{:else if !$currentOrder.data}
			<p class="text-center">Loading next order...</p>
		{:else if $currentOrder.data}
			{#if $currentOrder.data.nextOrder}
				<OrderRenderer
					expanded
					total={0.0}
					order={$currentOrder.data.nextOrder}
				/>
			{:else}
				<p>Waiting for new order...</p>
			{/if}
		{/if}
	</div>
</div>

<style></style>
