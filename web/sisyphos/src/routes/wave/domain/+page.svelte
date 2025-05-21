<script lang="ts">
  import { createClient } from "@connectrpc/connect";
  import { createConnectTransport } from "@connectrpc/connect-web";
  import {DomainService} from "$lib/types/wave/v1/domain/service_pb"
  import { create } from '@bufbuild/protobuf';
  import { type Domain, ListRequestSchema } from "$lib/types/wave/v1/domain/message_pb";
    import { SetException } from "$lib/exception/exception.svelte";

  const transport = createConnectTransport({
    baseUrl: "http://127.0.0.1:1871",
  })

  const client = createClient(DomainService, transport)

  let domains: {[key: string]: Domain} = $state({})

  async function list() {
    try {
      const request = create(ListRequestSchema, {});

      const response = await client.list(request)
      domains = response.domains
    } catch (err: any) {
      SetException({title: "Domain List", desc: err.message})
    }
  } 
</script>

<h1>Domain</h1>

<button onclick="{() => list()}" class="bg-orange-700 p-2 rounded-2xl">
  Click me
</button>

{#each Object.entries(domains) as [id, domain]}
  <div>Domain {id} with name {domain.config?.name}</div>
{/each}
