<script lang="ts">
  import { createClient } from "@connectrpc/connect";
  import { createConnectTransport } from "@connectrpc/connect-web";
  import {DomainService} from "$lib/sdk/types/wave/v1/domain/service_pb"
  import { create } from '@bufbuild/protobuf';
  import { type Domain, UpdateRequestSchema, CreateRequestSchema, GetRequestSchema, DomainSchema, StatRequestSchema } from "$lib/sdk/types/wave/v1/domain/message_pb";
  import { SetException } from "$lib/exception/exception.svelte";
  import { Arch, Chipset, DomainState, Firmware } from "$lib/sdk/types/wave/v1/domain/config_pb";
  import Button from "$lib/component/Button/Button.svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import Spice from "$lib/component/Spice/Spice.svelte";
  import Serial from "$lib/component/Serial/Serial.svelte";
  import Radio from "$lib/component/Radio/Radio.svelte";
  import { DomainPowerState, type DomainStats, DomainStatsSchema } from "$lib/sdk/types/wave/v1/domain/stat_pb";
  import { flip } from "svelte/animate";
  import VDropdown from "$lib/component/VDropdown/VDropdown.svelte";
  import Dropdown from "$lib/component/Dropdown/Dropdown.svelte";
  import { type Disk, DiskSchema, ListRequestSchema } from "$lib/sdk/types/granit/v1/disk/message_pb";
  import { DiskService } from "$lib/sdk/types/granit/v1/disk/service_pb";
  import Input from "$lib/component/Input/Input.svelte";
    import { DiskClient, DomainClient } from "$lib/client/client.svelte";

  let domain: Domain = $state(create(DomainSchema, {config: {}}))
  let stats: DomainStats = $state(create(DomainStatsSchema, {}))

  let disks: {[key: string]: Disk} = $state({})

  let mode: "serial" | "spice" = $state("serial");

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

  async function listDisks() {
    try {
      // TODO: Remove, currently test data becuase not implemented
      disks = {
        "c056d400-7c56-4c84-827a-d038a4c24b01": create(DiskSchema, {config: {size: 10000}}),
      }
      return
      // TODO: Remove, currently test data becuase not implemented

      const request = create(ListRequestSchema, {});

      const response = await DiskClient().list(request)
      if (response.disks) {
        disks = response.disks
      }
    } catch (err: any) {
      SetException({title: "LIST DISKS", desc: err.message})
    }
  }

  $effect.root(() => {
    if (page.params.id !== "new") {
      getDomain(page.params.id)
      if (domain.node !== "") {
        statDomain(page.params.id)
      }
    }
  })
</script>

<div class="w-11/12 flex flex-col gap-4 p-2 mt-20">
  <div class="w-full h-[500px] flex flex-row justify-between rounded-xl bg-slate-950/20">
    <div class="relative w-full flex flex-col items-start p-4">
      <div class="absolute top-4 right-4 rounded-lg p-2 bg-slate-950/10">
        {#if stats.state === DomainPowerState.DOMAIN_RUNNING}
          <span class="font-bold text-green-700/60">running</span>
        {:else if stats.state === DomainPowerState.DOMAIN_PAUSED}
          <span class="font-bold text-orange-700/60">paused</span>
        {:else if stats.state === DomainPowerState.DOMAIN_PMSUSPENDED}
          <span class="font-bold text-orange-900/60">suspended</span>
        {:else if stats.state === DomainPowerState.DOMAIN_BLOCKED}
          <span class="font-bold text-amber-900/60">blocked</span>
        {:else if stats.state === DomainPowerState.DOMAIN_CRASHED}
          <span class="font-bold text-red-900/80">crashed</span>
        {:else if stats.state === DomainPowerState.DOMAIN_SHUTDOWN}
          <span class="font-bold text-red-800/60">shutdown</span>
        {:else if stats.state === DomainPowerState.DOMAIN_SHUTOFF}
          <span class="font-bold text-red-800/60">shutoff</span>
        {:else if stats.state === DomainPowerState.DOMAIN_NOSTATE}
          <span class="font-bold text-slate-100/60">nostate</span>
        {/if}
      </div>
      
      <h1 class="text-2xl">Name: <span class="font-bold">{domain.config?.name}</span></h1>
      <p class="opacity-70">{page.params.id}</p>
      <hr class="w-full my-4">
      <div class="w-full flex flex-row gap-8 justify-start">
        <h2 class="flex flex-col items-start">
          <span class="text-sm font-bold">Node</span>
          <span class="text-2xl font-bold opacity-50">{domain.node !== "" ? domain.node : "<none>"}</span>
        </h2>
        <h2 class="flex flex-col items-start">
          <span class="text-sm font-bold">Reqnode</span>
          <span class="text-2xl font-bold opacity-50">{domain.reqnode !== "" ? domain.reqnode : "<none>"}</span>
        </h2>
      </div>
      <hr>
    </div>
    <div class="w-full flex flex-col items-center bg-slate-600/10">
      <div class="w-full flex flex-row justify-around gap-8 p-3">
        <Radio value="serial" bind:group={mode} 
          class="flex flex-row gap-2 justify-center w-full p-2 rounded-lg bg-slate-200/20" selectedClass="bg-slate-50/30">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none"><path fill="currentColor" fill-opacity="0.16" d="M20.6 4H3.4A2.4 2.4 0 0 0 1 6.4v11.2A2.4 2.4 0 0 0 3.4 20h17.2a2.4 2.4 0 0 0 2.4-2.4V6.4A2.4 2.4 0 0 0 20.6 4"/><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-miterlimit="10" stroke-width="1.5" d="m5 16l4-4l-4-4m6 8h8M3.4 4h17.2A2.4 2.4 0 0 1 23 6.4v11.2a2.4 2.4 0 0 1-2.4 2.4H3.4A2.4 2.4 0 0 1 1 17.6V6.4A2.4 2.4 0 0 1 3.4 4"/></g></svg>
        </Radio>
        <Radio value="spice" bind:group={mode} 
          class="flex flex-row gap-2 justify-center w-full p-2 rounded-lg bg-slate-200/20" selectedClass="bg-slate-50/30">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M4 17h16V5H4zm11-2.5h2.5V12H19v4h-4zM5 6h4v1.5H6.5V10H5z" opacity="0.3"/><path fill="currentColor" d="M20 3H4c-1.11 0-2 .89-2 2v12a2 2 0 0 0 2 2h4v2h8v-2h4c1.1 0 2-.9 2-2V5a2 2 0 0 0-2-2m0 14H4V5h16z"/><path fill="currentColor" d="M6.5 7.5H9V6H5v4h1.5zM19 12h-1.5v2.5H15V16h4z"/></svg>
        </Radio>
      </div>
      {#if mode === "serial"}
        <Serial></Serial>
      {:else if mode === "spice"}
        <Spice></Spice>
      {/if}
    </div>
  </div>

  <div class="w-full h-[400px] flex flex-row p-2 rounded-xl bg-slate-950/20 overflow-scroll-hidden">
    <div class="w-1/3 flex flex-col items-start gap-4 p-4">
      <input placeholder="Name" bind:value={domain.config!.name} 
        class="text-xl w-full p-1 rounded-md bg-slate-50/10 focus:bg-slate-50/20 focus:outline-0 transition-all overflow-hidden" />
      <input placeholder="Description" bind:value={domain.config!.description} 
        class="text-sm w-full p-1 rounded-md bg-slate-50/10 focus:bg-slate-50/20 focus:outline-0 transition-all overflow-hidden" />
      <div class="flex flex-wrap gap-2 justify-stretch">
        <input placeholder="Tag" onkeyup={(e: any) => {
          if (e.key === "Enter" && e.target?.value) {
            domain.config!.affinity.push(e.target.value)
            domain = domain
            e.target.value = ""
          }
        }} class="w-24 p-1 bg-slate-50/10 rounded-lg focus:outline-0" />
        {#each domain.config!.affinity as tag, i (i)}
          <button animate:flip onclick={() => {
            domain.config!.affinity.splice(i, 1)
            domain = domain
          }} class="p-1 flex flex-row gap-1 justify-center items-center cursor-pointer bg-slate-50/10 rounded-lg">
            <span class="max-w-24 overflow-hidden">{tag}</span>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-dasharray="24" stroke-dashoffset="24" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path d="M5 5l14 14"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.4s" values="24;0"/></path><path d="M19 5l-14 14"><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.4s" dur="0.4s" values="24;0"/></path></g></svg>
          </button>
        {/each}
      </div>
    </div>
    <span class="h-full w-0 border-1 rounded-full"></span>
    <div class="w-1/3 flex flex-col items-start gap-4 p-4">
      {#if domain.config!.resourceConfig}
        <div class="w-full flex flex-row gap-2">
          <Input title="vCPUs" type="number" class="flex-4/12" bind:value={domain.config!.resourceConfig!.vcpus} />
          <Input title="Memory (bytes)" type="number" class="flex-8/12" bind:value={domain.config!.resourceConfig!.memory} />
          <span class="w-40 p-1 rounded-md text-xl text-nowrap overflow-hidden">
            ~ {(Number(domain.config!.resourceConfig!.memory) / (1000 * 1000 * 1000)).toFixed(2)} GB
          </span>
        </div>
      {/if}
      {#if domain.config!.firmwareConfig}
        <VDropdown placeholder="Firmware" bind:value={domain.config!.firmwareConfig!.firmware} items={{
          "OVMF": Firmware.OVMF, "SEABIOS": Firmware.SEABIOS
        }} class="w-full"></VDropdown>
        <Dropdown placeholder="Loader Device" bind:value={domain.config!.firmwareConfig!.loaderDeviceId} loader={async () => {
          await listDisks()
          return disks
        }} class="w-full"></Dropdown>

        <button onclick={() => {
          domain.config!.firmwareConfig!.secureBoot = !domain.config!.firmwareConfig!.secureBoot
        }} class="w-full p-1 rounded-md cursor-pointer transition-all {domain.config!.firmwareConfig!.secureBoot ? "font-bold bg-slate-50/20" : "bg-slate-50/10"}">
          Secureboot
        </button>
      {/if}
    </div>
    <span class="h-full w-0 border-1 rounded-full"></span>
    <div class="w-1/3 flex flex-col items-start gap-4 p-4">
    </div>
  </div>

  <div class="w-full flex flex-row justify-between">
    {#if page.params.id === "new"}
      <Button onclick={() => domain = create(DomainSchema, {config: {}})} scale={0.8} class="flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/40">
        <span>Reset</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="40" stroke-dashoffset="40" d="M17 15.33c2.41 -0.72 4 -1.94 4 -3.33c0 -2.21 -4.03 -4 -9 -4c-4.97 0 -9 1.79 -9 4c0 2.06 3.5 3.75 8 3.98"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.5s" values="40;0"/></path><path fill="currentColor" d="M12.25 16l0 0l0 0z" opacity="0"><animate fill="freeze" attributeName="d" begin="0.5s" dur="0.2s" values="M12.25 16l0 0l0 0z;M12.25 16L9.5 13.25L9.5 18.75z"/><set fill="freeze" attributeName="opacity" begin="0.5s" to="1"/></path></g></svg>      </Button>
      <Button onclick={() => createDomain()} scale={0.8} class="flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/40">
        <span>Create</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-dasharray="16" stroke-dashoffset="16" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path d="M5 12h14"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.4s" values="16;0"/></path><path d="M12 5v14"><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.4s" dur="0.4s" values="16;0"/></path></g></svg>
      </Button>
    {:else}
      <Button onclick={() => getDomain(page.params.id)} scale={0.8} class="flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/40">
        <span>Reset</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="40" stroke-dashoffset="40" d="M17 15.33c2.41 -0.72 4 -1.94 4 -3.33c0 -2.21 -4.03 -4 -9 -4c-4.97 0 -9 1.79 -9 4c0 2.06 3.5 3.75 8 3.98"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.5s" values="40;0"/></path><path fill="currentColor" d="M12.25 16l0 0l0 0z" opacity="0"><animate fill="freeze" attributeName="d" begin="0.5s" dur="0.2s" values="M12.25 16l0 0l0 0z;M12.25 16L9.5 13.25L9.5 18.75z"/><set fill="freeze" attributeName="opacity" begin="0.5s" to="1"/></path></g></svg>
      </Button>
      <Button onclick={() => updateDomain(page.params.id)} scale={0.8} class="flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/40">
        <span>Update</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="64" stroke-dashoffset="64" d="M3 12c0 -4.97 4.03 -9 9 -9c4.97 0 9 4.03 9 9c0 4.97 -4.03 9 -9 9c-4.97 0 -9 -4.03 -9 -9Z"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.6s" values="64;0"/></path><path stroke-dasharray="14" stroke-dashoffset="14" d="M8 12l3 3l5 -5"><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.6s" dur="0.2s" values="14;0"/></path></g></svg>
      </Button>
    {/if}
  </div>
</div>

