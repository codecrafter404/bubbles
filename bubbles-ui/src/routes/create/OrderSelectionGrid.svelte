<script lang="ts">
	import type { CreateOrderInput, CustomItem, Item } from '../../gql/graphql';
	import SelectionItem from './SelectionItem.svelte';
	import { getNextCustomItem, getNextItems } from './ViewHelper';
	import { KeyboardPressedEvent, ProcessEvent } from '$lib/Keyboard';
	import { BREAK } from 'graphql';

	interface Props {
		currentOrder: CreateOrderInput;
		items: Array<Item>;
		customItems: Array<CustomItem>;
		inputAllowed: boolean;
	}

	let { currentOrder = $bindable(), items, customItems, inputAllowed }: Props = $props();

	// should only be run, if data is successfully fetched

	let currently_displaying = $derived.by(() => {
		return getNextItems(items, customItems, currentOrder);
	});

	let currently_selected = $state(0);
	let cols = $state(3);

	function on_key_down(event: KeyboardEvent) {
		if (!inputAllowed) {
			return;
		}
		let res = ProcessEvent(event);

		if (res == undefined) {
			return;
		}

		switch (res!) {
			case KeyboardPressedEvent.MoveUp:
				if (currently_selected - cols >= 0) {
					currently_selected -= cols;
				}
				break;
			case KeyboardPressedEvent.MoveDown:
				if (currently_selected + cols < currently_displaying.length) {
					currently_selected += cols;
				}
				break;
			case KeyboardPressedEvent.MoveLeft:
				if (currently_selected <= 0) {
					currently_selected = currently_displaying.length - 1;
					break;
				}
				currently_selected = (currently_selected - 1) % currently_displaying.length;
				break;
			case KeyboardPressedEvent.MoveRight:
				currently_selected = (currently_selected + 1) % currently_displaying.length;
				break;
			case KeyboardPressedEvent.Add:
				addItem();
				console.log($state.snapshot(currentOrder));
				break;
			case KeyboardPressedEvent.Remove:
				removeItem();
				console.log($state.snapshot(currentOrder));
				break;
			default:
				return;
		}
		event.preventDefault();
	}

	function addItem() {
		let [next, next_idx] = getNextCustomItem(customItems, currentOrder);
		let current = currently_displaying[currently_selected];
		if (next == undefined) {
			if (current.isCustomItem) {
				currentOrder.customItems.push({
					masterCustomItemID: current.id,
					quantity: 1,
					selectedCustomItems: []
				});
			} else {
				let idx = currentOrder.items.findIndex((x) => x.itemID == current.id);
				if (idx == -1) {
					currentOrder.items.push({
						itemID: current.id,
						quantity: 1
					});
				} else {
					currentOrder.items[idx].quantity += 1;
				}
			}
		} else {
			if (current.isCustomItem) {
				console.error('Variants may only be items and not customitems');
			}

			currentOrder.customItems[next_idx].selectedCustomItems.push({
				customItemID: next.id,
				selectedVariantIDs: [current.id]
			});
			let exists_next = customItems.find((x) => x.prev != undefined && x.prev! == next.id);
			if (exists_next == undefined) {
				// dedup
				let mod = currentOrder.customItems[next_idx];
				for (let i = 0; i < currentOrder.customItems.length; i++) {
					if (next_idx == i) continue;

					let current = currentOrder.customItems[i];
					if (
						current.masterCustomItemID == mod.masterCustomItemID &&
						JSON.stringify(current.selectedCustomItems) == JSON.stringify(mod.selectedCustomItems)
					) {
						currentOrder.customItems[i].quantity += 1;
						currentOrder.customItems.splice(next_idx, 1);
						break;
					}
				}
			}
			currently_selected = 0;
		}
	}

	function removeItem() {
		let [next, next_idx] = getNextCustomItem(customItems, currentOrder);
		let current = currently_displaying[currently_selected];

		if (next != undefined) {
			if (currentOrder.customItems[next_idx].selectedCustomItems.length >= 1) {
				currentOrder.customItems[next_idx].selectedCustomItems.pop();
			} else {
				currentOrder.customItems.splice(next_idx, 1);
			}
			return;
		}

		if (current.isCustomItem) {
			let idx = currentOrder.customItems.findLastIndex((x) => x.masterCustomItemID == current.id);
			// no custom items with this id in the current order
			if (idx == -1) {
				return;
			}

			if (currentOrder.customItems[idx].quantity > 1) {
				currentOrder.customItems[idx].quantity -= 1;
			} else {
				currentOrder.customItems.splice(idx, 1);
			}
		} else {
			let idx = currentOrder.items.findIndex((x) => x.itemID == current.id);
			// no items with this id in the current order
			if (idx == -1) {
				return;
			}

			if (currentOrder.items[idx].quantity > 1) {
				currentOrder.items[idx].quantity -= 1;
			} else {
				currentOrder.items.splice(idx, 1);
			}
		}
	}
</script>

<svelte:document on:keydown={on_key_down} />
<div class="grid bg-red-300" style={`grid-template-columns: repeat(${cols}, minmax(0, 1fr))`}>
	{#each currently_displaying as item, i}
		<SelectionItem {...item} selected={i == currently_selected} />
	{/each}
</div>
