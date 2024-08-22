<script lang="ts">
	import { createEventDispatcher } from "svelte";
	import type { CustomItem, Item } from "../../generated/graphql";

	export let item: Item | null;
	export let customItem: CustomItem | null;
</script>

<div class="bg-primary-500 w-full h-full row-auto">
	{#if item != null}
		<img src={item.image} alt={item.name} class="object-contain" />
		<p>{item.name}</p>
	{:else if customItem != null}
		<p>{customItem.name}</p>
	{:else}
		<p>Invalid component Props</p>
	{/if}
	<p></p>
</div>
