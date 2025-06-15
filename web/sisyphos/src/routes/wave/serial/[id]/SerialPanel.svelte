<script lang="ts">
  import { type Serial } from "$lib/sdk/types/wave/v1/serial/message_pb";
  import Input from "$lib/component/Input/Input.svelte";

  export type PanelType = "serial"|"video"|"adapter"|"storage"|"network"|undefined

  let {
    serial = $bindable(),
  }: {
    serial: Serial;
  } = $props();
</script>

<div class="w-full flex flex-row p-2 rounded-xl bg-slate-950/20 overflow-scroll-hidden">
  <div class="w-1/3 flex flex-col items-start gap-4 p-4">
    <span class="flex flex-row gap-2">
      <Input title="Name" type="text" bind:value={serial.config!.name} class="w-full"></Input>
    </span>
  </div>
</div>