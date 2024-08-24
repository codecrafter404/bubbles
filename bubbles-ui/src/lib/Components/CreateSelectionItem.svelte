<script lang="ts">
	import { afterUpdate, createEventDispatcher, onMount } from "svelte";
	import type { CustomItem, Item } from "../../generated/graphql";

	export let item: Item | null;

	export let customItem: CustomItem | null;
	export let customItemPriceStart: number | null;
	export let customItemPriceEnd: number | null;

	export let hovered: boolean = true;

	export let quantity: number = 2;

	let currentElem: HTMLDivElement;
	afterUpdate(() => {
		if (hovered) {
			currentElem.scrollIntoView({ behavior: "smooth" });
		}
	});
</script>

<div
	class={"ring-2 ring-opacity-30 aspect-[3/4] flex flex-col overflow-hidden rounded-md z-0 relative shadow-sm" +
		(hovered ? " ring-accent-500 shadow-xl ring-4" : "")}
	bind:this={currentElem}
>
	{#if quantity != 0}
		<div
			class="absolute top-2 right-2 z-10 aspect-square bg-primary-400 h-6 w-6 text-center rounded-md font-bold text-white"
		>
			{quantity}
		</div>
	{/if}
	{#if item != null}
		<img
			src={item.image}
			alt={item.name}
			class="aspect-square w-full object-center object-cover flex-grow"
		/>
		<div class="flex justify-between p-1 bg-primary-100">
			<p class="max-w-[75%]">
				{item ? item.name : customItem?.name}
			</p>
			<p class="font-bold">
				{item.price}€
			</p>
		</div>
	{:else if customItem != null}
		<div
			class="flex-grow flex items-center justify-center bg-natural-500 text-3xl font-bold text-center"
		>
			<p>{customItem.name}</p>
		</div>
		<p class="text-end font-bold bg-primary-100 p-1">
			{customItemPriceStart}€ - {customItemPriceEnd}€
		</p>
	{:else}
		<p>Invalid component Props</p>
	{/if}
</div>
