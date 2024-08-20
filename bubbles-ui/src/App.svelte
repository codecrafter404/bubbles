<script lang="ts">
  import Router from "svelte-spa-router";
  import CreateOrder from "./lib/Pages/CreateOrder.svelte";
  import FulfillOrder from "./lib/Pages/FulfillOrder.svelte";
  import Index from "./lib/Pages/Index.svelte";

  const routes = {
    "/": Index,
    "/create": CreateOrder,
    "/fulfill/:orderId": FulfillOrder,
  };
</script>

<body>
  <Router {routes} />
</body>

<style>
</style>
