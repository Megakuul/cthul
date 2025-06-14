<script lang="ts">
  import Dropdown from "$lib/component/Dropdown/Dropdown.svelte";
  import Input from "$lib/component/Input/Input.svelte";
  import VDropdown from "$lib/component/VDropdown/VDropdown.svelte";
  import { NetworkBus, type NetworkDevice } from "$lib/sdk/types/wave/v1/domain/config_pb";

  import protonInterface from "$lib/assets/proton-interface.svg";

  let {
    id,
    device = $bindable()
  }: {
    id: number,
    device: NetworkDevice
  } = $props();
</script>

<div class="w-full flex flex-col items-start gap-4 p-4 rounded-xl bg-slate-950/20 overflow-scroll-hidden">
  <h1 class="w-full flex flex-row justify-between items-center gap-2">
    <span class="text-2xl font-medium">Network interface {id}</span>
    <img width="40" alt="" src={protonInterface} />
  </h1>
  <div class="w-full flex flex-row gap-2">
    <Dropdown title="Device" bind:value={device.deviceId} loader={async () => {
      return {} // TODO
    }} class=""></Dropdown>
    <VDropdown title="Bus" bind:value={device.networkBus} items={{
      "E1000": NetworkBus.E1000, "VirtIO": NetworkBus.VIRTIO,
    }} class=""></VDropdown>
    <Input title="Bootpriority" type="number" bind:value={device.bootPriority} class=""></Input>
  </div>
</div>