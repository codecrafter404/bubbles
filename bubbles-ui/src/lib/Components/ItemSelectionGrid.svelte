<script lang="ts">
	import { createEventDispatcher, onDestroy, onMount } from "svelte";
	import type { SelectableItem } from "../Types/SelectableItem";
	import CreateSelectionItem from "./CreateSelectionItem.svelte";

	export let items: SelectableItem[];

	export let gridCols = 5;

	export let currentItem = 0;

	const dispatch = createEventDispatcher<{
		add: SelectableItem;
		remove: SelectableItem;
	}>();

	function handle_key_down(e: KeyboardEvent) {
		if (e.code == "Numpad8") {
			// up
			//TODO: make behavior configurable

			if (currentItem - gridCols >= 0) {
				currentItem -= gridCols;
			}
			// currentItem -= gridCols;
			// if (currentItem < 0) currentItem = 0;
		}
		if (e.code == "Numpad6") {
			// right
			if (currentItem + 1 < items.length) {
				currentItem += 1;
			}
		}
		if (e.code == "Numpad2") {
			// down
			if (currentItem + gridCols < items.length) {
				currentItem += gridCols;
			}
			// } else {
			// 	currentItem = items.length - 1; //TODO: make behavior configurable
			// }
		}
		if (e.code == "Numpad4") {
			// left
			currentItem -= 1;
			if (currentItem < 0) currentItem = 0;
		}
		if (e.code == "NumpadEnter") {
			// add to order
			dispatch("add", items[currentItem]);
		}
		if (e.code == "NumpadSubtract") {
			// remove from order
			dispatch("remove", items[currentItem]);
		}
	}
	onMount(() => {
		document.addEventListener("keydown", handle_key_down);
	});
	onDestroy(() => {
		document.removeEventListener("keydown", handle_key_down);
	});
</script>

<div
	class="grid gap-3 items-center"
	style={`grid-template-columns: repeat(${gridCols}, minmax(0, 1fr))`}
>
	{#each items as item, i}
		<CreateSelectionItem {item} hovered={currentItem == i} />
	{/each}
</div>
