<script lang="ts">
  import { createClient } from "@connectrpc/connect";
  import { createConnectTransport } from "@connectrpc/connect-web";
  import {DomainService} from "$lib/types/wave/v1/domain/service_pb"
  import { create } from '@bufbuild/protobuf';
  import { type Domain, ListRequestSchema, CreateRequestSchema } from "$lib/types/wave/v1/domain/message_pb";
  import { SetException } from "$lib/exception/exception.svelte";
  import { Arch, Chipset, DomainState, Firmware } from "$lib/types/wave/v1/domain/config_pb";
  import Button from "$lib/component/Button/Button.svelte";
  import Link from "$lib/component/Link/Link.svelte";
    import { flip } from "svelte/animate";

  const transport = createConnectTransport({
    baseUrl: "http://127.0.0.1:1870",
  })

  const client = createClient(DomainService, transport)

  let domains: {[key: string]: Domain} = $state({})

  let filteredDomains: {[key: string]: Domain} = $state({})

  let search: string = $state("");

  async function listDomains() {
    try {
      const request = create(ListRequestSchema, {});

      const response = await client.list(request)
      domains = response.domains
    } catch (err: any) {
      SetException({title: "Domain List", desc: err.message})
    }
  }

  async function createDomain() {
    try {
      const request = create(CreateRequestSchema, {
        config: {
          name: "test",
          description: "this is a test",
          affinity: [""],
          firmwareConfig: {
            firmware: Firmware.OVMF,
            loaderDeviceId: "asdf",
            nvramDeviceId: "asdf",
            tmplDeviceId: "asdf",
            secureBoot: false,
          },
          title: "whatistitle",
          state: DomainState.DOWN,
          systemConfig: {
            architecture: Arch.AMD64,
            chipset: Chipset.Q35,
          },
          resourceConfig: {
            memory: BigInt(2000000000),
            vcpus: BigInt(2),
          },
        }
      })

      const response = await client.create(request)
      console.log(response.id)
    } catch (err: any) {
      SetException({title: "Domain List", desc: err.message})
    }
  }

  $effect.root(() => {
    listDomains();
  })

  $effect(() => {
    const newDomains: {[key: string]: Domain} = {};
    for (const [k, v] of Object.entries(domains)) {
      if (k.includes(search) || v.config?.name.includes(search) || v.config?.description.includes(search))
      newDomains[k] = v
    }
    filteredDomains = newDomains
  })
</script>

<div class="w-11/12 flex flex-col gap-4 p-2 mt-20">
  <div class="flex flex-row items-center justify-between">
    <Button onclick={() => listDomains()} scale={0.8} class="flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/40">
      <span>Refresh</span>
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="40" stroke-dashoffset="40" d="M17 15.33c2.41 -0.72 4 -1.94 4 -3.33c0 -2.21 -4.03 -4 -9 -4c-4.97 0 -9 1.79 -9 4c0 2.06 3.5 3.75 8 3.98"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.5s" values="40;0"/></path><path fill="currentColor" d="M12.25 16l0 0l0 0z" opacity="0"><animate fill="freeze" attributeName="d" begin="0.5s" dur="0.2s" values="M12.25 16l0 0l0 0z;M12.25 16L9.5 13.25L9.5 18.75z"/><set fill="freeze" attributeName="opacity" begin="0.5s" to="1"/></path></g></svg>
    </Button>
    
    <Button onclick={() => createDomain()} scale={0.8} class="flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/40">
      <span>New</span>
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-dasharray="16" stroke-dashoffset="16" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path d="M5 12h14"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.4s" values="16;0"/></path><path d="M12 5v14"><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.4s" dur="0.4s" values="16;0"/></path></g></svg>
    </Button>
  </div>
  <input bind:value={search} placeholder="Search..." class="w-full rounded-lg p-2 bg-slate-50/30 focus:outline-none" />

  <div class="flex flex-col gap-4 p-2 h-[700px] overflow-scroll-hidden">
    {#each Object.entries(filteredDomains) as [id, domain] (id)}
      <div animate:flip={{ duration: 250 }}>
        <Link href="/wave/domain/{id}" class="w-full relative flex flex-col items-start p-4 rounded-lg bg-slate-50/30">
          <h1>
            <span class="text-xl font-bold">{domain.config?.name.slice(0, 16)}</span>
            <span class="text-lg font-medium text-slate-800/60">#{id.slice(0, 6)}</span>
          </h1>
          <p class="text-sm">{domain.config?.description}</p>
          <div class="flex flex-row gap-2 mt-2 opacity-70">
            <p class="text-xs">
              reqnode: <span class="font-bold">{domain.reqnode ? domain.reqnode : "none"}</span>
            </p>
            <p class="text-xs">
              node: <span class="font-bold">{domain.node ? domain.node : "none"}</span>
            </p>
            <p class="text-xs">
              affinity: <span class="font-bold">[{domain.config?.affinity.join(", ")}]</span>
            </p>
          </div>
          {#if domain.config?.state === DomainState.UP}
            <div title="up" class="absolute right-10 top-1/2 -translate-y-1/2 text-green-800/80">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" fill-opacity="0" stroke="currentColor" stroke-dasharray="64" stroke-dashoffset="64" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12c0 -4.97 4.03 -9 9 -9c4.97 0 9 4.03 9 9c0 4.97 -4.03 9 -9 9c-4.97 0 -9 -4.03 -9 -9Z"><animate fill="freeze" attributeName="fill-opacity" begin="0.6s" dur="0.15s" values="0;0.3"/><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.6s" values="64;0"/></path></svg>
            </div>
          {:else if domain.config?.state === DomainState.PAUSE}
            <div title="paused" class="absolute right-10 top-1/2 -translate-y-1/2 text-orange-800/60">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" fill-opacity="0" stroke="currentColor" stroke-dasharray="64" stroke-dashoffset="64" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12c0 -4.97 4.03 -9 9 -9c4.97 0 9 4.03 9 9c0 4.97 -4.03 9 -9 9c-4.97 0 -9 -4.03 -9 -9Z"><animate fill="freeze" attributeName="fill-opacity" begin="0.6s" dur="0.15s" values="0;0.3"/><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.6s" values="64;0"/></path></svg>
            </div>
          {:else if domain.config?.state === DomainState.DOWN}
            <div title="down" class="absolute right-10 top-1/2 -translate-y-1/2 text-stone-800/60">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" fill-opacity="0" stroke="currentColor" stroke-dasharray="64" stroke-dashoffset="64" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12c0 -4.97 4.03 -9 9 -9c4.97 0 9 4.03 9 9c0 4.97 -4.03 9 -9 9c-4.97 0 -9 -4.03 -9 -9Z"><animate fill="freeze" attributeName="fill-opacity" begin="0.6s" dur="0.15s" values="0;0.3"/><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.6s" values="64;0"/></path></svg>
            </div>
          {:else if domain.config?.state === DomainState.FORCED_DOWN}
            <div title="forced down" class="absolute right-10 top-1/2 -translate-y-1/2 text-stone-950/70">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" fill-opacity="0" stroke="currentColor" stroke-dasharray="64" stroke-dashoffset="64" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12c0 -4.97 4.03 -9 9 -9c4.97 0 9 4.03 9 9c0 4.97 -4.03 9 -9 9c-4.97 0 -9 -4.03 -9 -9Z"><animate fill="freeze" attributeName="fill-opacity" begin="0.6s" dur="0.15s" values="0;0.3"/><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.6s" values="64;0"/></path></svg>
            </div>
          {/if}
        </Link>
      </div>
    {/each}
  </div>
</div>

<center class="font-bold text-5xl mt-10 pb-20 select-none">***</center>