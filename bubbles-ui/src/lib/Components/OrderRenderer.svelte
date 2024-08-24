<script lang="ts">
	import type { Order } from "../../generated/graphql";
	export let order: Order;
	export let total: number;
	export let expanded: boolean = false;

	if (expanded) {
		// sort
	}
</script>

<div class="h-full p-3">
	<div
		class="flex flex-col h-full ring-natural-200 ring-4 rounded-md text-lg"
	>
		<!-- Make overflow scrollable -->
		<div class="flex-grow h-full overflow-y-scroll">
			<ul class="w-full p-4">
				{#each order.customItems as item}
					<li class="li-item px-2">
						<span
							class="flex justify-between"
							><span>
								{item.customItem
									.name}</span
							><span>
								{item.quantity}
								x
								{#if !expanded}
									{(() => {
										let res = 0;
										item.customItem.variants.forEach(
											(
												x,
											) =>
												(res +=
													x.price),
										);
										return res.toFixed(
											2,
										);
									})()}€
								{/if}
							</span></span
						>
						<ul
							class="pl-6 list-disc list-inside"
						>
							{#each item.customItem.variants as v}
								<!-- TODO: update me -->
								<li>
									{v.name}
									- x{item.quantity}
								</li>
							{/each}
						</ul>
					</li>
				{/each}
				{#each order.items as item}
					<li
						class="li-item px-2 flex justify-between"
					>
						<span>
							{#if expanded}
								<span
									>({item
										.item
										.identifier})</span
								>
							{/if}
							{item.item.name}</span
						><span
							class={expanded
								? "font-bold"
								: ""}
							>{item.quantity} x

							{#if !expanded}
								{item.item.price.toFixed(
									2,
								)}€
							{/if}
						</span>
					</li>
				{/each}
			</ul>
		</div>
		<h1 class="text-center text-2xl p-4 bg-primary-100">
			{#if expanded}
				<span class="text-neutral-700">Order </span>
				<span class="font-bold text-accent-500"
					>#{order.identifier}</span
				>
			{:else}
				<span class="text-neutral-700">Total </span>
				<span>{total.toFixed(2)}€</span>
			{/if}
		</h1>
	</div>
</div>

<style>
	.li-item:nth-of-type(2n - 1) {
		background-color: #bfc5c5;
		border-radius: 0.375rem /* 6px */;
	}
</style>
