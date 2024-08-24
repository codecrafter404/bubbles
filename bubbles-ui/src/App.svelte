<script lang="ts">
  import Router from "svelte-spa-router";
  import CreateOrder from "./lib/Pages/CreateOrder.svelte";
  import FulfillOrder from "./lib/Pages/FulfillOrder.svelte";
  import Index from "./lib/Pages/Index.svelte";
  import { createClient as createWSClient } from "graphql-ws";
  import {
    cacheExchange,
    Client,
    fetchExchange,
    setContextClient,
    subscriptionExchange,
  } from "@urql/svelte";
  import { load_endpoint, load_ws_endpoint } from "./lib/utils/storage";

  const routes = {
    "/": Index,
    "/create": CreateOrder,
    "/fulfill": FulfillOrder,
  };
  const wsClient = createWSClient({
    url: load_ws_endpoint(),
  });
  const client = new Client({
    url: load_endpoint(),
    exchanges: [
      cacheExchange,
      fetchExchange,
      subscriptionExchange({
        forwardSubscription(request) {
          const input = { ...request, query: request.query || "" };
          return {
            subscribe(sink) {
              const unsubscribe = wsClient.subscribe(input, sink);
              return { unsubscribe };
            },
          };
        },
      }),
    ],
  });
  setContextClient(client);
</script>

<body>
  <Router {routes} />
</body>

<style>
</style>
