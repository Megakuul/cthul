<script lang="ts">
  import { create } from '@bufbuild/protobuf';
  import { SetException } from "$lib/exception/exception.svelte";
  import Button from "$lib/component/Button/Button.svelte";
  import Link from "$lib/component/Link/Link.svelte";
  import { flip } from "svelte/animate";
  import { NodeClient } from "$lib/client/client.svelte";
  import { ListRequestSchema, type Node } from "$lib/sdk/types/wave/v1/node/message_pb";
  import { NodeState } from '$lib/sdk/types/wave/v1/node/config_pb';

  let nodes: {[key: string]: Node} = $state({})

  let filteredNodes: {[key: string]: Node} = $state({})

  let search: string = $state("");

  async function listNodes() {
    try {
      const request = create(ListRequestSchema, {});

      const response = await NodeClient().list(request)
      nodes = response.nodes
    } catch (err: any) {
      SetException({title: "LIST NODES", desc: err.message})
    }
  }

  $effect.root(() => {
    listNodes();
  })

  $effect(() => {
    const newNodes: {[key: string]: Node} = {};
    for (const [k, v] of Object.entries(nodes)) {
      if (k.includes(search))
        newNodes[k] = v
    }
    filteredNodes = newNodes
  })
</script>

<div class="w-11/12 flex flex-col gap-4 p-2 mt-20">
  <div class="flex flex-row items-center justify-between">
    <Button onclick={() => listNodes()} scale={0.8} class="flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/40">
      <span>Refresh</span>
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="40" stroke-dashoffset="40" d="M17 15.33c2.41 -0.72 4 -1.94 4 -3.33c0 -2.21 -4.03 -4 -9 -4c-4.97 0 -9 1.79 -9 4c0 2.06 3.5 3.75 8 3.98"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.5s" values="40;0"/></path><path fill="currentColor" d="M12.25 16l0 0l0 0z" opacity="0"><animate fill="freeze" attributeName="d" begin="0.5s" dur="0.2s" values="M12.25 16l0 0l0 0z;M12.25 16L9.5 13.25L9.5 18.75z"/><set fill="freeze" attributeName="opacity" begin="0.5s" to="1"/></path></g></svg>
    </Button>
  </div>
  <input bind:value={search} placeholder="Search..." class="w-full rounded-lg p-2 bg-slate-50/30 focus:outline-none" />

  <div class="flex flex-col gap-4 p-2 h-[700px] overflow-scroll-hidden">
    {#each Object.entries(filteredNodes) as [id, node] (id)}
      <div animate:flip={{ duration: 250 }}>
        <Link href="/wave/node/{id}" class="w-full relative flex flex-col items-start p-4 rounded-lg bg-slate-50/30">
          <h1>
            <span class="text-xl font-bold">{id}</span>
          </h1>
          <div class="flex flex-row gap-2 mt-2 opacity-70">
            <p class="text-xs">
              cpu: <span class="font-bold">
                {node.config?.availableCpu.toFixed(1)} 
                / {node.config?.allocatedCpu.toFixed(1)}
              </span>
            </p>
            <p class="text-xs">
              memory: <span class="font-bold">
                {node.config?.availableMemory ?? 0 / (1000 * 1000 * 1000)}
                / {node.config?.allocatedMemory ?? 0 / (1000 * 1000 * 1000)} GB
              </span>
            </p>
            <p class="text-xs">
              state: <span class="font-bold">
                {#if node.config?.state === NodeState.HEALTHY}
                  healthy
                {:else if node.config?.state === NodeState.MAINTENANCE}
                  maintenance
                {:else if node.config?.state === NodeState.DEGRADED}
                  degraded
                {/if}
              </span>
            </p>
            <p class="text-xs">
              affinity: <span class="font-bold">[{node.config?.affinity.join(", ")}]</span>
            </p>
            {#if node.error}
              <p class="text-xs">
                error: <span class="font-bold text-red-900">{node.error}</span>
              </p>
            {/if}
          </div>
        </Link>
      </div>
    {/each}
  </div>
</div>

<center class="font-bold text-5xl mt-10 pb-20 select-none">***</center>