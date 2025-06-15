<script lang="ts">
  import { type Domain } from "$lib/sdk/types/wave/v1/domain/message_pb";
  import { page } from "$app/state";
  import Spice from "$lib/component/Spice/Spice.svelte";
  import Serial from "$lib/component/Serial/Serial.svelte";
  import Radio from "$lib/component/Radio/Radio.svelte";
  import { DomainPowerState, type DomainStats } from "$lib/sdk/types/wave/v1/domain/stat_pb";

  let {
    domain = $bindable(),
    stats = $bindable(),
  }: {
    domain: Domain;
    stats: DomainStats;
  } = $props();

  let mode: "serial" | "spice" = $state("serial");
</script>

<div class="w-full h-[500px] flex flex-row justify-between rounded-xl bg-slate-950/20 overflow-hidden">
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
    <hr class="w-full my-4">

    <!-- TODO: analyze result data and then visualize with https://github.com/vnau/svelte-gauge -->
    <div class="w-full flex flex-row gap-8 justify-start">
      <h2 class="flex flex-col items-start">
        <span class="text-sm font-bold">CPU</span>
        <span class="text-2xl font-bold opacity-50">vcpus {stats.cpu?.vcpus}</span>
        <span class="text-2xl font-bold opacity-50">cpu time {stats.cpu?.cpuTime}</span>
        <span class="text-2xl font-bold opacity-50">kernel time {stats.cpu?.kernelTime}</span>
        <span class="text-2xl font-bold opacity-50">user time {stats.cpu?.userTime}</span>
      </h2>
      <h2 class="flex flex-col items-start">
        <span class="text-sm font-bold">Memory</span>
        <span class="text-2xl font-bold opacity-50">available {stats.memory?.available ?? 0 / (1000 * 1000 * 1000)} GB</span>
      </h2>
    </div>
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
      <Serial devices={domain.config?.serialDevices ?? []}></Serial>
    {:else if mode === "spice"}
      <Spice></Spice>
    {/if}
  </div>
</div>