<script lang="ts">
	import {
		OrderState,
		type CreateOrderInput,
		type CustomItem,
		type Item,
		type Order,
		type OrderCustomItem,
		type OrderItem,
		type SelectedCustomItem
	} from '../../gql/graphql';

	interface Props {
		orderInput: CreateOrderInput | Order;
		customItems: Array<CustomItem>;
		items: Array<Item>;
	}

	let { orderInput, customItems, items }: Props = $props();

	let order: Order = $derived.by(() => {
		if (orderInput.hasOwnProperty('total')) {
			return orderInput as Order;
		}
		let inp = orderInput as CreateOrderInput;
		let res: Order = {
			id: -1,
			identifier: 'NOT SUBMITTED',
			customItems: inp.customItems.map((x) => {
				return {
					...x,
					id: -1,
					masterCustomItem: customItems.find((y) => y.id == x.masterCustomItemID)!,
					selectedCustomItems: x.selectedCustomItems.map((x) => {
						let ret: SelectedCustomItem = {
							id: -1,
							customItem: customItems.find((y) => y.id == x.customItemID)!,
							selectedVariants: x.selectedVariantIDs.map((i) => items.find((y) => y.id == i)!)
						};
						return ret;
					})
				};
			}),
			items: inp.items.map((x) => {
				return {
					...x,
					id: -1,
					item: items.find((y) => y.id == x.itemID)!
				};
			}),
			submitted: '',
			state: OrderState.Created,
			total: 0
		};
		res.total = calculateTotal(res.customItems!, res.items!);
		return res;
	});

	function calculateTotal(customItems: Array<OrderCustomItem>, item: Array<OrderItem>): number {
		let total = 0;
		customItems.forEach((x) => {
			total +=
				x.quantity *
				x
					.selectedCustomItems!.flatMap((x) => x.selectedVariants)
					.map((x) => x.price)
					.reduce((sum, c) => {
						return sum + c;
					}, 0);
		});
		item.forEach((x) => {
			total += x.quantity * x.item.price;
		});
		return total;
	}
</script>

<div>
	<ul>
		{#each order.customItems! as ci}
			<li>
				<p>{ci.quantity}x <span class="font-bold">{ci.masterCustomItem.name}</span></p>
				<ul>
					{#each ci.selectedCustomItems! as sci}
						{#if sci.selectedVariants.length == 1}
							<li>{sci.selectedVariants[0].name}</li>
						{:else if sci.selectedVariants.length > 1}
							<li>
								<ul>
									{#each sci.selectedVariants as sv}
										<li>{sv.name}</li>
									{/each}
								</ul>
							</li>
						{/if}
					{/each}
				</ul>
			</li>
		{/each}
		{#each order.items! as item}
			<li><p>{item.quantity}x <span class="font-bold">{item.item.name}</span></p></li>
		{/each}
	</ul>
</div>
