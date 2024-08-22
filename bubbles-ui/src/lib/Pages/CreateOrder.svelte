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
</script>

<div class="flex flex-col min-h-[100vh] 2xl:max-h-[100vh]">
	<NavBar />
	<div class="flex-grow flex flex-col 2xl:flex-row">
		<div class="flex-grow bg-blue-100">
			{#if $res.fetching}
				<p>Loading</p>
			{:else if $res.error}
				<p>Failed to load: {$res.error.message}</p>
			{:else if $res.data}
				<div
					class="grid overflow-x-scroll auto-rows-max grid-cols-1 gap-5"
				>
					{#each $res.data.getItems as item}
						{#if item.isOneOff && item.available}
							<CreateSelectionItem
								{item}
								customItem={null}
							/>
						{/if}
					{/each}
				</div>
			{/if}
		</div>
		<div class="2xl:w-[30vw] bg-orange-200">
			OrderList (Shopping Cart)
		</div>
	</div>
</div>

<style></style>
