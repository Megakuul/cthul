<script lang="ts">
  import { ListRequestSchema, type Domain } from "$lib/sdk/types/wave/v1/domain/message_pb";
  import { page } from "$app/state";
  import Spice from "$lib/component/Spice/Spice.svelte";
  import Serial from "$lib/component/Serial/Serial.svelte";
  import Radio from "$lib/component/Radio/Radio.svelte";
  import { DomainPowerState, type DomainStats } from "$lib/sdk/types/wave/v1/domain/stat_pb";
  import type { Node } from "$lib/sdk/types/wave/v1/node/message_pb";
  import { create } from "@bufbuild/protobuf";
  import { DomainClient } from "$lib/client/client.svelte";
  import { SetException } from "$lib/exception/exception.svelte";
  import Link from "$lib/component/Link/Link.svelte";

  import waveDomain from "$lib/assets/wave-domain.svg"
    import { NodeState } from "$lib/sdk/types/wave/v1/node/config_pb";

  let {
    node = $bindable(),
  }: {
    node: Node;
  } = $props();

  let domains: {[key: string]: Domain} | undefined = $state(undefined)

  async function listDomains() {
    try {
      const request = create(ListRequestSchema, {});

      const response = await DomainClient().list(request)
      const localDomains: {[key: string]: Domain} = {};
      for (const [id, domain] of Object.entries(response.domains)) {
        if (domain.node === page.params.id) {
          localDomains[id] = domain
        }
      }
      domains = localDomains
    } catch (err: any) {
      SetException({title: "LIST DOMAINS", desc: err.message})
    }
  }

  $effect.root(() => {
    if (page.params.id !== "new") {
      listDomains()
    }
  })
</script>

<div class="w-full h-[500px] flex flex-row justify-between rounded-xl bg-slate-950/20">
  <div class="relative w-full flex flex-col items-start p-4">
    <div class="absolute top-4 right-4 rounded-lg p-2 bg-slate-950/10">
      {#if node.config?.state === NodeState.HEALTHY}
        <span class="font-bold text-green-700/60">healthy</span>
      {:else if node.config?.state === NodeState.MAINTENANCE}
        <span class="font-bold text-orange-900/60">maintenance</span>
      {:else if node.config?.state === NodeState.DEGRADED}
        <span class="font-bold text-red-800/60">degraded</span>
      {/if}
    </div>
    
    <h1 class="text-2xl">Name: <span class="font-bold">{page.params.id}</span></h1>
    <hr class="w-full my-4">

    <!-- TODO: analyze result data and then visualize with https://github.com/vnau/svelte-gauge -->
    <div class="w-full flex flex-row gap-8 justify-start">
      <h2 class="flex flex-col items-start">
        <span class="text-sm font-bold">CPU</span>
        <span class="text-2xl font-bold opacity-50">
          {node.config?.availableCpu} 
          / {node.config?.allocatedCpu}
        </span>
      </h2>
      <h2 class="flex flex-col items-start">
        <span class="text-sm font-bold">Memory</span>
        <span class="text-2xl font-bold opacity-50">
          {node.config?.availableMemory ?? 0 / (1000 * 1000 * 1000)}
          / {node.config?.allocatedMemory ?? 0 / (1000 * 1000 * 1000)} GB
        </span>
      </h2>
      <h2 class="flex flex-col items-start">
        <span class="text-sm font-bold">Affinity</span>
        {#each node.config?.affinity ?? [] as affinity}
          <span class="text-2xl font-bold opacity-50">
            {affinity}
          </span>
        {/each}
      </h2>
    </div>
  </div>
  <div class="w-full flex flex-col gap-4 items-center justify-center bg-slate-600/10">
    {#if domains}
      {#each Object.entries(domains) as [id, domain]}
        <Link href="/wave/domain/{id}" scale={0.8} class="h-12 flex flex-row justify-between items-center bg-slate-50/20 rounded-lg">
          <h1 class="text-xl font-bold">{domain.config?.name}</h1>
          <img alt="" src={waveDomain}>
        </Link>
      {/each}
    {:else}
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g stroke="currentColor" stroke-width="1"><circle cx="12" cy="12" r="9.5" fill="none" stroke-linecap="round" stroke-width="3"><animate attributeName="stroke-dasharray" calcMode="spline" dur="1.5s" keySplines="0.42,0,0.58,1;0.42,0,0.58,1;0.42,0,0.58,1" keyTimes="0;0.475;0.95;1" repeatCount="indefinite" values="0 150;42 150;42 150;42 150"/><animate attributeName="stroke-dashoffset" calcMode="spline" dur="1.5s" keySplines="0.42,0,0.58,1;0.42,0,0.58,1;0.42,0,0.58,1" keyTimes="0;0.475;0.95;1" repeatCount="indefinite" values="0;-16;-59;-59"/></circle><animateTransform attributeName="transform" dur="2s" repeatCount="indefinite" type="rotate" values="0 12 12;360 12 12"/></g></svg>
    {/if}
  </div>
</div>