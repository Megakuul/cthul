<script lang="ts">
  import { create } from '@bufbuild/protobuf';
  import { type Domain } from "$lib/sdk/types/wave/v1/domain/message_pb";
  import { SetException } from "$lib/exception/exception.svelte";
  import { Arch, Chipset, DomainState, Firmware, NetworkBus, NetworkDeviceSchema, SerialBus, SerialDeviceSchema, StorageBus, StorageDeviceSchema, StorageType, Video, VideoAdapterSchema, VideoDeviceSchema } from "$lib/sdk/types/wave/v1/domain/config_pb";
  import Button from "$lib/component/Button/Button.svelte";
  import { flip } from "svelte/animate";
  import VDropdown from "$lib/component/VDropdown/VDropdown.svelte";
  import Dropdown from "$lib/component/Dropdown/Dropdown.svelte";
  import { type Disk, ListRequestSchema } from "$lib/sdk/types/granit/v1/disk/message_pb";
  import Input from "$lib/component/Input/Input.svelte";
  import { DiskClient } from "$lib/client/client.svelte";

  import waveSerial from "$lib/assets/wave-serial.svg";
  import waveVideo from "$lib/assets/wave-video.svg";
  import granitDisk from "$lib/assets/granit-disk.svg";
  import protonInterface from "$lib/assets/proton-interface.svg";

  export type PanelType = "serial"|"video"|"adapter"|"storage"|"network"|undefined

  let {
    domain = $bindable(),
    panel = $bindable(),
  }: {
    domain: Domain;
    panel: {type: PanelType, id: number}
  } = $props();

  let disks: {[key: string]: Disk} = $state({})

  async function listDisks() {
    try {
      const request = create(ListRequestSchema, {});

      const response = await DiskClient().list(request)
      if (response.disks) {
        disks = response.disks
      }
    } catch (err: any) {
      SetException({title: "LIST DISKS", desc: err.message})
    }
  }

  let deviceGroups = $derived([
    {title: "Serial devices", type: "serial", icon: waveSerial, new: create(SerialDeviceSchema, {
      deviceId: "",
      serialBus: SerialBus.ISA,
    }), devices: domain.config!.serialDevices},

    {title: "Video devices", type: "video", icon: waveVideo, new: create(VideoDeviceSchema, {
      video: Video.QXL,
      videobufferSize: BigInt(1000 * 1000),
      commandbufferSize: BigInt(1000 * 1000),
      framebufferSize: BigInt(1000 * 1000),
    }), devices: domain.config!.videoDevices},

    {title: "Video adapters", type: "adapter", icon: waveVideo, new: create(VideoAdapterSchema, {
      deviceId: "",
    }), devices: domain.config!.videoAdapters},

    {title: "Storage devices", type: "storage", icon: granitDisk, new: create(StorageDeviceSchema, {
      deviceId: "",
      bootPriority: BigInt(10),
      storageBus: StorageBus.IDE,
      storageType: StorageType.DISK,
    }), devices: domain.config!.storageDevices},

    {title: "Network devices", type: "network", icon: protonInterface, new: create(NetworkDeviceSchema, {
      deviceId: "",
      bootPriority: BigInt(20),
      networkBus: NetworkBus.E1000,
    }), devices: domain.config!.networkDevices},
  ])
</script>

<div class="w-full h-[420px] flex flex-row p-2 rounded-xl bg-slate-950/20 overflow-scroll-hidden">
  <div class="w-1/3 flex flex-col items-start gap-4 p-4">
    <span class="flex flex-row gap-2">
      <Input title="Name" type="text" bind:value={domain.config!.name} class="w-full"></Input>
      <VDropdown title="state" bind:value={domain.config!.state} items={{
        "up": DomainState.UP, "pause": DomainState.PAUSE, "down": DomainState.DOWN, "forced down": DomainState.FORCED_DOWN,
      }} class=""></VDropdown>
    </span>
    <Input title="Description" type="text" bind:value={domain.config!.description} class="w-full text-sm"></Input>

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
    {#if domain.config!.systemConfig}
      <div class="w-full flex flex-row gap-2">
        <VDropdown title="Architecture" bind:value={domain.config!.systemConfig!.architecture} items={{
          "AMD64": Arch.AMD64, "AARCH64": Arch.AARCH64
        }} class=""></VDropdown>
        <VDropdown title="Chipset" bind:value={domain.config!.systemConfig!.chipset} items={{
          "I440FX": Chipset.I440FX, "Q35": Chipset.Q35, "VirtIO": Chipset.VIRT,
        }} class=""></VDropdown>
      </div>
    {/if}
    {#if domain.config!.firmwareConfig}
      <VDropdown title="Firmware" bind:value={domain.config!.firmwareConfig!.firmware} items={{
        "OVMF": Firmware.OVMF, "SEABIOS": Firmware.SEABIOS
      }} class="w-full"></VDropdown>
      <Dropdown title="Firmware Loader" bind:value={domain.config!.firmwareConfig!.loaderDeviceId} loader={async () => {
        await listDisks()
        return disks
      }} class="w-full"></Dropdown>
      <Dropdown title="Firmware Template" bind:value={domain.config!.firmwareConfig!.tmplDeviceId} loader={async () => {
        await listDisks()
        return disks
      }} class="w-full"></Dropdown>
      <Dropdown title="Firmware NVRAM" bind:value={domain.config!.firmwareConfig!.nvramDeviceId} loader={async () => {
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

  <div class="w-1/3 flex flex-col items-start gap-2 p-4">
    {#each deviceGroups as group}
      <div class="w-full flex flex-row justify-between gap-4">
        <h1 class="text-lg font-medium">{group.title}</h1>
        <Button class="p-1 rounded-lg bg-slate-50/20" scale={0.5} onclick={() => {
          group.devices.push(group.new as any)
          domain = domain
          panel = {type: group.type as PanelType, id: group.devices.length - 1}
        }}>
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-dasharray="16" stroke-dashoffset="16" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path d="M5 12h14"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.4s" values="16;0"/></path><path d="M12 5v14"><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.4s" dur="0.4s" values="16;0"/></path></g></svg>
        </Button>
      </div>
      <div class="flex flex-row justify-start gap-4">
        {#each group.devices as device, i (device)}
          <span animate:flip>
            <Button  class="relative w-full p-1 rounded-lg bg-slate-50/20" scale={0.5} onclick={() => {
              panel = {type: group.type as PanelType, id: i}
            }}>
              <img title="{group.type} device {i}" width="26" alt={i.toString()} src={group.icon} />
              <Button class="absolute -top-1.5 -right-1.5 p-1 rounded-lg bg-slate-50/10" scale={0} onclick={(e: Event) => {
                e.stopPropagation();
                panel = {type: undefined, id: 0}
                group.devices.splice(i, 1)
                domain = domain
              }}>
                <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-dasharray="16" stroke-dashoffset="16" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path d="M7 7l10 10"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.4s" values="16;0"/></path><path d="M17 7l-10 10"><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.4s" dur="0.4s" values="16;0"/></path></g></svg>
              </Button>
            </Button>
          </span>
        {/each}
      </div>
    {/each}
  </div>
</div>