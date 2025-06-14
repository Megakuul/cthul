<script lang="ts">
  import Input from "$lib/component/Input/Input.svelte";
  import VDropdown from "$lib/component/VDropdown/VDropdown.svelte";
  import { Video, type VideoDevice } from "$lib/sdk/types/wave/v1/domain/config_pb";

  import waveVideo from "$lib/assets/wave-video.svg";

  let {
    id,
    device = $bindable()
  }: {
    id: number,
    device: VideoDevice
  } = $props();
</script>

<div class="w-full flex flex-col items-start gap-4 p-4 rounded-xl bg-slate-950/20 overflow-scroll-hidden">
  <h1 class="w-full flex flex-row justify-between items-center gap-2">
    <span class="text-2xl font-medium">Video buffer device {id}</span>
    <img width="40" alt="" src={waveVideo} />
  </h1>
  <div class="w-full flex flex-row gap-2">
    <VDropdown title="Type" bind:value={device.video} items={{
      "QXL": Video.QXL, "VGA": Video.VGA, "HOST": Video.HOST,
    }} class=""></VDropdown>
    <Input title="Framebuffer size (bytes)" type="number" bind:value={device.framebufferSize} class=""></Input>
    <Input title="Commandbuffer size (bytes)" type="number" bind:value={device.commandbufferSize} class=""></Input>
    <Input title="Videobuffer size (bytes)" type="number" bind:value={device.videobufferSize} class=""></Input>
  </div>
</div>