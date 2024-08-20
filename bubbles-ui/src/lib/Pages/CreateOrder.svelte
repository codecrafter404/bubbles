<script lang="ts">
	import NavBar from "../Components/NavBar.svelte";
</script>

<div class="flex flex-col min-h-[100vh]">
	<NavBar />
	<div class="flex-grow flex flex-col 2xl:flex-row">
		<div class="flex-grow bg-blue-100">Select Order Items</div>
		<div class="2xl:w-[30vw] bg-orange-200">
			OrderList (Shopping Cart)
		</div>
	</div>
</div>

<style></style>
