<script lang="ts">
	import { link } from "svelte-spa-router";
	import {
		load_endpoint,
		load_ws_endpoint,
		update_endpoint,
		update_ws_endpoint,
	} from "../utils/storage";
	let graph_url = load_endpoint();
	let graph_ws_url = load_ws_endpoint();
</script>

<div>
	<h1>Useful Links</h1>
	<div class="w-full justify-between inline-grid">
		<a href="/create" use:link>Create</a>
		<a href="/fulfill/" use:link>Fulfill</a>
	</div>
	<br />
	<h1>Setup Endpoint URL</h1>
	<div class="flex flex-col m-2">
		<input
			bind:value={graph_url}
			class="ring-4 ring-offset-inherit ring-primary-500 rounded"
		/>
		<button
			on:click={() => update_endpoint(graph_url)}
			class="bg-primary-500 mt-3 p-3 rounded font-bold"
			>update url</button
		>
	</div>
	<br />
	<h1>Setup Endpoint WS URL</h1>
	<div class="flex flex-col m-2">
		<input
			bind:value={graph_ws_url}
			class="ring-4 ring-offset-inherit ring-primary-500 rounded"
		/>
		<button
			on:click={() => update_ws_endpoint(graph_ws_url)}
			class="bg-primary-500 mt-3 p-3 rounded font-bold"
			>update ws url</button
		>
	</div>
</div>

<style>
	h1 {
		font-size: 26px;
	}
	a {
		text-decoration: underline;
	}
</style>
