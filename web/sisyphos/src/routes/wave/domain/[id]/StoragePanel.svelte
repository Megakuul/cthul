<script lang="ts">
  import { DiskClient } from "$lib/client/client.svelte";
  import Dropdown from "$lib/component/Dropdown/Dropdown.svelte";
  import Input from "$lib/component/Input/Input.svelte";
  import VDropdown from "$lib/component/VDropdown/VDropdown.svelte";
  import { SetException } from "$lib/exception/exception.svelte";
  import { StorageBus, StorageType, type StorageDevice } from "$lib/sdk/types/wave/v1/domain/config_pb";
  import { create } from "@bufbuild/protobuf";
  import { ListRequestSchema, type Disk } from "$lib/sdk/types/granit/v1/disk/message_pb";

  import granitDisk from "$lib/assets/granit-disk.svg";
    

  let {
    id,
    device = $bindable()
  }: {
    id: number,
    device: StorageDevice
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
</script>

<div class="w-full flex flex-col items-start gap-4 p-4 rounded-xl bg-slate-950/20 overflow-scroll-hidden">
  <h1 class="w-full flex flex-row justify-between items-center gap-2">
    <span class="text-2xl font-medium">Storage devices {id}</span>
    <img width="40" alt="" src={granitDisk} />
  </h1>
  <div class="w-full flex flex-row gap-2">
    <Dropdown title="Device" bind:value={device.deviceId} loader={async () => {
      await listDisks()
      return disks
    }} class=""></Dropdown>
    <VDropdown title="Type" bind:value={device.storageType} items={{
      "DISK": StorageType.DISK, "CDROM": StorageType.CDROM,
    }} class=""></VDropdown>
    <VDropdown title="Bus" bind:value={device.storageBus} items={{
      "IDE": StorageBus.IDE, "SATA": StorageBus.SATA, "VirtIO": StorageBus.VIRTIO,
    }} class=""></VDropdown>
    <Input title="Bootpriority" type="number" bind:value={device.bootPriority} class=""></Input>
  </div>
</div>