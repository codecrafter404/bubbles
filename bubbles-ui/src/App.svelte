<script lang="ts">
  import Router from "svelte-spa-router";
  import CreateOrder from "./lib/Pages/CreateOrder.svelte";
  import FulfillOrder from "./lib/Pages/FulfillOrder.svelte";
  import Index from "./lib/Pages/Index.svelte";
  import {
    cacheExchange,
    Client,
    fetchExchange,
    setContextClient,
  } from "@urql/svelte";
  import { load_endpoint } from "./lib/utils/storage";

  const routes = {
    "/": Index,
    "/create": CreateOrder,
    "/fulfill/:orderId": FulfillOrder,
  };

  const client = new Client({
    url: load_endpoint(),
    exchanges: [cacheExchange, fetchExchange],
  });
  setContextClient(client);
</script>

<body>
  <Router {routes} />
</body>

<style>
</style>
