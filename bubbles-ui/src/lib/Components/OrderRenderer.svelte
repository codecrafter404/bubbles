<script lang="ts">
	import type { Order } from "../../generated/graphql";
	export let order: Order;
</script>

<div class="h-full p-3">
	<div class="flex flex-col h-full ring-natural-200 ring-4 rounded-md">
		<ul class="flex-grow">
			{#each order.items as item}
				<li>{item.item.name} - x{item.quantity}</li>
			{/each}
		</ul>
		<h1 class="text-center text-2xl p-4 bg-primary-100">
			<span class="text-neutral-700">Total </span>
			<span
				>{(() => {
					let sum = 0;
					order.items.forEach(
						(x) =>
							(sum +=
								x.quantity *
								x.item.price),
					);
					order.customItems.forEach((x) => {
						x.customItem.variants.forEach(
							(y) => {
								sum +=
									y.price *
									x.quantity;
							},
						);
					});
					return sum;
				})()}€</span
			>
		</h1>
	</div>
</div>
