<script lang="ts">
  import { SerialClient } from "$lib/client/client.svelte";
  import Dropdown from "$lib/component/Dropdown/Dropdown.svelte";
  import Input from "$lib/component/Input/Input.svelte";
  import VDropdown from "$lib/component/VDropdown/VDropdown.svelte";
  import { SetException } from "$lib/exception/exception.svelte";
  import { SerialBus, type SerialDevice } from "$lib/sdk/types/wave/v1/domain/config_pb";
  import { ListRequestSchema, type Serial } from "$lib/sdk/types/wave/v1/serial/message_pb";
  import { create } from "@bufbuild/protobuf";

  import waveSerial from "$lib/assets/wave-serial.svg";

  let {
    id,
    device = $bindable()
  }: {
    id: number,
    device: SerialDevice
  } = $props();

  let serials: {[key: string]: Serial} = $state({})

  async function listSerial() {
    try {
      const request = create(ListRequestSchema, {});

      const response = await SerialClient().list(request)
      if (response.serials) {
        serials = response.serials
      }
    } catch (err: any) {
      SetException({title: "LIST SERIALS", desc: err.message})
    }
  }
</script>

<div class="w-full flex flex-col items-start gap-4 p-4 rounded-xl bg-slate-950/20 overflow-scroll-hidden">
  <h1 class="w-full flex flex-row justify-between items-center gap-2">
    <span class="text-2xl font-medium">Serial device {id}</span>
    <img width="40" alt="" src={waveSerial} />
  </h1>
  <div class="w-full flex flex-row gap-2">
    <Dropdown title="Device" bind:value={device.deviceId} loader={async () => {
      await listSerial()
      return serials
    }} class=""></Dropdown>
    <VDropdown title="Bus" bind:value={device.serialBus} items={{
      "ISA": SerialBus.ISA, "VirtIO": SerialBus.VIRTIO,
    }} class=""></VDropdown>
    <Input title="Port" type="number" bind:value={device.port} class=""></Input>
  </div>
</div>