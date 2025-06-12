<script lang="ts">
	import { afterUpdate } from "svelte";
	import type { SelectableItem } from "../Types/SelectableItem";

	export let item: SelectableItem;
	export let hovered: boolean = true;

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
	{#if item.quantity != 0}
		<div
			class="absolute top-2 right-2 z-10 aspect-square bg-primary-400 h-6 w-6 text-center rounded-md font-bold text-white"
		>
			{item.quantity}
		</div>
	{/if}
	<img
		src={item.image == ""
			? "/customitem_placeholder.png"
			: item.image}
		alt={item.name}
		class={`aspect-square w-full object-center ${item.image == "" ? "object-cover" : "object-cover"} flex-grow`}
	/>
	<div class="flex justify-between p-1 bg-primary-100">
		<p class="max-w-[75%]">
			{item.name}
		</p>
		<p class="font-bold">
			{item.price}
		</p>
	</div>
</div>
