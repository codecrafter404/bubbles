<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { Client, cacheExchange, fetchExchange, setContextClient } from '@urql/svelte';

	let { children } = $props();

	const client = new Client({
		url: 'http://localhost:8080/query',
		exchanges: [cacheExchange, fetchExchange]
	});
	setContextClient(client);
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children?.()}
