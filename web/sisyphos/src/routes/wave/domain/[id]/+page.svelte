<script lang="ts">
  import { create } from '@bufbuild/protobuf';
  import { type Domain, UpdateRequestSchema, CreateRequestSchema, GetRequestSchema, DomainSchema, StatRequestSchema, AttachRequestSchema, DetachRequestSchema, DeleteRequestSchema } from "$lib/sdk/types/wave/v1/domain/message_pb";
  import { SetException } from "$lib/exception/exception.svelte";
  import Button from "$lib/component/Button/Button.svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import { type DomainStats, DomainStatsSchema } from "$lib/sdk/types/wave/v1/domain/stat_pb";
  import { DomainClient } from "$lib/client/client.svelte";

  import Overview from "./Overview.svelte";
  import DomainPanel, { type PanelType } from "./DomainPanel.svelte";
  import SerialPanel from './SerialPanel.svelte';
  import VideoPanel from './VideoPanel.svelte';
  import AdapterPanel from './AdapterPanel.svelte';
  import StoragePanel from './StoragePanel.svelte';
  import NetworkPanel from './NetworkPanel.svelte';
  import { Firmware } from '$lib/sdk/types/wave/v1/domain/config_pb';

  let domain: Domain = $state(create(DomainSchema, {node: "", reqnode: "", config: {
    name: "",
    description: "",
    affinity: [],
    resourceConfig: {vcpus: BigInt(2), memory: BigInt(4 * (1000 * 1000 * 1000))},
    firmwareConfig: {firmware: Firmware.SEABIOS, 
      loaderDeviceId: "", tmplDeviceId: "", nvramDeviceId: "", secureBoot: false,
    },
    serialDevices: [],
    videoDevices: [],
    videoAdapters: [],
    inputDevices: [],
    storageDevices: [],
    networkDevices: [],
  }}))
  let stats: DomainStats = $state(create(DomainStatsSchema, {}))

  let panel: {type: PanelType, id: number} = $state({type: undefined, id: 0});

  async function getDomain(id: string) {
    try {
      const request = create(GetRequestSchema, {
        id: id,
      });

      const response = await DomainClient().get(request)
      if (response.domain) {
        domain = response.domain
      }
    } catch (err: any) {
      SetException({title: "RETRIEVE DOMAIN", desc: err.message})
    }
  }

  async function statDomain(id: string) {
    try {
      const request = create(StatRequestSchema, {
        id: id,
      });

      const response = await DomainClient().stat(request)
      if (response.stats) {
        stats = response.stats
      }
    } catch (err: any) {
      SetException({title: "STAT DOMAIN", desc: err.message})
    }
  }

  async function createDomain() {
    try {
      const request = create(CreateRequestSchema, {
        config: domain.config,
      })

      const response = await DomainClient().create(request)
      goto(`/wave/domain/${response.id}`);
    } catch (err: any) {
      SetException({title: "CREATE DOMAIN", desc: err.message})
    }
  }

  async function updateDomain(id: string) {
    try {
      const request = create(UpdateRequestSchema, {
        id: id,
        config: domain.config,
      })

      await DomainClient().update(request)
    } catch (err: any) {
      SetException({title: "UPDATE DOMAIN", desc: err.message})
    }
  }

  async function attachDomain(id: string, node: string) {
    try {
      const request = create(AttachRequestSchema, {
        id: id,
        node: node,
      })

      await DomainClient().attach(request)
    } catch (err: any) {
      SetException({title: "ATTACH DOMAIN", desc: err.message})
    }
  }

  async function detachDomain(id: string) {
    try {
      const request = create(DetachRequestSchema, {
        id: id,
      })

      await DomainClient().detach(request)
    } catch (err: any) {
      SetException({title: "DETACH DOMAIN", desc: err.message})
    }
  }

  async function deleteDomain(id: string) {
    try {
      const request = create(DeleteRequestSchema, {
        id: id,
      })

      await DomainClient().delete(request)
      goto(`/wave/domain`);
    } catch (err: any) {
      SetException({title: "DELETE DOMAIN", desc: err.message})
    }
  }

  $effect.root(() => {
    if (page.params.id !== "new") {
      getDomain(page.params.id)
      setInterval(() => {
        if (domain.node !== "") {
          statDomain(page.params.id)
        }
      }, 1000)
    }
  })
</script>

<div class="w-11/12 flex flex-col gap-4 p-2 mt-20">
  <Overview bind:domain={domain} bind:stats={stats}></Overview>

  <DomainPanel bind:domain={domain} bind:panel={panel}></DomainPanel>

  {#if panel.type === "serial"}
    <SerialPanel id={panel.id} bind:device={domain.config!.serialDevices[panel.id]}></SerialPanel>
  {:else if panel.type === "video"}
    <VideoPanel id={panel.id} bind:device={domain.config!.videoDevices[panel.id]}></VideoPanel>
  {:else if panel.type === "adapter"}
    <AdapterPanel id={panel.id} bind:device={domain.config!.videoAdapters[panel.id]}></AdapterPanel>
  {:else if panel.type === "storage"}
    <StoragePanel id={panel.id} bind:device={domain.config!.storageDevices[panel.id]}></StoragePanel>
  {:else if panel.type === "network"}
    <NetworkPanel id={panel.id} bind:device={domain.config!.networkDevices[panel.id]}></NetworkPanel>
  {/if}

  <div class="w-full flex flex-row justify-start gap-3 my-2">
    {#if page.params.id === "new"}
      <Button onclick={() => domain = create(DomainSchema, {config: {}})} scale={0.8} class="flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/30">
        <span>Reset</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="40" stroke-dashoffset="40" d="M17 15.33c2.41 -0.72 4 -1.94 4 -3.33c0 -2.21 -4.03 -4 -9 -4c-4.97 0 -9 1.79 -9 4c0 2.06 3.5 3.75 8 3.98"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.5s" values="40;0"/></path><path fill="currentColor" d="M12.25 16l0 0l0 0z" opacity="0"><animate fill="freeze" attributeName="d" begin="0.5s" dur="0.2s" values="M12.25 16l0 0l0 0z;M12.25 16L9.5 13.25L9.5 18.75z"/><set fill="freeze" attributeName="opacity" begin="0.5s" to="1"/></path></g></svg>      
      </Button>
      <Button onclick={() => createDomain()} scale={0.8} class="flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/30">
        <span>Create</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-dasharray="16" stroke-dashoffset="16" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path d="M5 12h14"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.4s" values="16;0"/></path><path d="M12 5v14"><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.4s" dur="0.4s" values="16;0"/></path></g></svg>
      </Button>
    {:else}
      {#if domain.reqnode !== ""}
        <Button onclick={() => detachDomain(page.params.id)} scale={0.8} class="flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/30">
          <span>Detach</span>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><mask id="lineMdEngineOff0"><g fill="none" stroke="#fff" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="48" stroke-dashoffset="48" d="M11 9h6v10h-6.5l-2 -2h-2.5v-6.5l1.5 -1.5Z"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.6s" values="48;0"/></path><path fill="#fff" fill-opacity="0" d="M17 13h0v-3h0v8h0v-3h0z" opacity="0"><animate fill="freeze" attributeName="d" begin="0.6s" dur="0.2s" values="M17 13h0v-3h0v8h0v-3h0z;M17 13h4v-3h1v8h-1v-3h-4z"/><set fill="freeze" attributeName="fill-opacity" begin="0.8s" to="1"/><set fill="freeze" attributeName="opacity" begin="0.6s" to="1"/></path><path d="M6 14h0M6 11v6" opacity="0"><animate fill="freeze" attributeName="d" begin="0.8s" dur="0.2s" values="M6 14h0M6 11v6;M6 14h-4M2 11v6"/><set fill="freeze" attributeName="opacity" begin="0.8s" to="1"/></path><path d="M11 9v0M8 9h6" opacity="0"><animate fill="freeze" attributeName="d" begin="1s" dur="0.2s" values="M11 9v0M8 9h6;M11 9v-4M8 5h6"/><set fill="freeze" attributeName="opacity" begin="1s" to="1"/></path><path stroke="#000" stroke-dasharray="28" stroke-dashoffset="28" d="M0 11h26" transform="rotate(45 12 12)"><animate fill="freeze" attributeName="stroke-dashoffset" begin="1.3s" dur="0.4s" values="28;0"/></path><path stroke-dasharray="28" stroke-dashoffset="28" d="M0 13h26" transform="rotate(45 12 12)"><animate fill="freeze" attributeName="stroke-dashoffset" begin="1.3s" dur="0.4s" values="28;0"/></path></g></mask><rect width="24" height="24" fill="currentColor" mask="url(#lineMdEngineOff0)"/></svg>
        </Button>
      {:else}
        <Button onclick={() => attachDomain(page.params.id, "")} scale={0.8} class="flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/30">
          <span>Attach</span>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="48" stroke-dashoffset="48" d="M11 9h6v10h-6.5l-2 -2h-2.5v-6.5l1.5 -1.5Z"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.6s" values="48;0"/></path><path fill="currentColor" fill-opacity="0" d="M17 13h0v-3h0v8h0v-3h0z" opacity="0"><animate fill="freeze" attributeName="d" begin="0.6s" dur="0.2s" values="M17 13h0v-3h0v8h0v-3h0z;M17 13h4v-3h1v8h-1v-3h-4z"/><set fill="freeze" attributeName="fill-opacity" begin="0.8s" to="1"/><set fill="freeze" attributeName="opacity" begin="0.6s" to="1"/></path><path d="M6 14h0M6 11v6" opacity="0"><animate fill="freeze" attributeName="d" begin="0.8s" dur="0.2s" values="M6 14h0M6 11v6;M6 14h-4M2 11v6"/><set fill="freeze" attributeName="opacity" begin="0.8s" to="1"/></path><path d="M11 9v0M8 9h6" opacity="0"><animate fill="freeze" attributeName="d" begin="1s" dur="0.2s" values="M11 9v0M8 9h6;M11 9v-4M8 5h6"/><set fill="freeze" attributeName="opacity" begin="1s" to="1"/></path></g></svg>
        </Button>
      {/if}

      <Button onclick={() => deleteDomain(page.params.id)} scale={0.8} class="ml-auto flex flex-row gap-2 justify-center w-32 p-2 rounded-lg text-slate-50/80 bg-red-900/80">
        <span>Delete</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="64" stroke-dashoffset="64" d="M12 3c4.97 0 9 4.03 9 9c0 4.97 -4.03 9 -9 9c-4.97 0 -9 -4.03 -9 -9c0 -4.97 4.03 -9 9 -9Z"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.6s" values="64;0"/></path><path stroke-dasharray="8" stroke-dashoffset="8" d="M12 12l4 4M12 12l-4 -4M12 12l-4 4M12 12l4 -4"><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.6s" dur="0.2s" values="8;0"/></path></g></svg>
      </Button>
      <Button onclick={() => getDomain(page.params.id)} scale={0.8} class="flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/30">
        <span>Reset</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="40" stroke-dashoffset="40" d="M17 15.33c2.41 -0.72 4 -1.94 4 -3.33c0 -2.21 -4.03 -4 -9 -4c-4.97 0 -9 1.79 -9 4c0 2.06 3.5 3.75 8 3.98"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.5s" values="40;0"/></path><path fill="currentColor" d="M12.25 16l0 0l0 0z" opacity="0"><animate fill="freeze" attributeName="d" begin="0.5s" dur="0.2s" values="M12.25 16l0 0l0 0z;M12.25 16L9.5 13.25L9.5 18.75z"/><set fill="freeze" attributeName="opacity" begin="0.5s" to="1"/></path></g></svg>
      </Button>
      <Button onclick={() => updateDomain(page.params.id)} scale={0.8} class="flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/30">
        <span>Update</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="64" stroke-dashoffset="64" d="M3 12c0 -4.97 4.03 -9 9 -9c4.97 0 9 4.03 9 9c0 4.97 -4.03 9 -9 9c-4.97 0 -9 -4.03 -9 -9Z"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.6s" values="64;0"/></path><path stroke-dasharray="14" stroke-dashoffset="14" d="M8 12l3 3l5 -5"><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.6s" dur="0.2s" values="14;0"/></path></g></svg>
      </Button>
    {/if}
  </div>
</div>

