<script lang="ts">
	import type { SelectionItemProps } from './ViewHelper';

	let { name, notes, price, image, count, isCustomItem, selected }: SelectionItemProps = $props();

	let price_s = $derived(`${price.toFixed(2)}€`);
</script>

<div>
	<p class={selected ? 'font-bold' : ''}>{name} - {price_s}</p>
</div>
