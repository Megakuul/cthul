<script lang="ts">
  import { VideoClient } from "$lib/client/client.svelte";
  import Dropdown from "$lib/component/Dropdown/Dropdown.svelte";
  import { SetException } from "$lib/exception/exception.svelte";
  import { type VideoAdapter } from "$lib/sdk/types/wave/v1/domain/config_pb";
  import { create } from "@bufbuild/protobuf";
  import { ListRequestSchema, type Video } from "$lib/sdk/types/wave/v1/video/message_pb";

  import waveVideo from "$lib/assets/wave-video.svg";
  
  let {
    id,
    device = $bindable()
  }: {
    id: number,
    device: VideoAdapter
  } = $props();

  let videos: {[key: string]: Video} = $state({})

  async function listVideo() {
    try {
      const request = create(ListRequestSchema, {});

      const response = await VideoClient().list(request)
      if (response.videos) {
        videos = response.videos
      }
    } catch (err: any) {
      SetException({title: "LIST VIDEOS", desc: err.message})
    }
  }
</script>

<div class="w-full flex flex-col items-start gap-4 p-4 rounded-xl bg-slate-950/20 overflow-scroll-hidden">
  <h1 class="w-full flex flex-row justify-between items-center gap-2">
    <span class="text-2xl font-medium">Video adapter {id}</span>
    <img width="40" alt="" src={waveVideo} />
  </h1>
  <div class="w-full flex flex-row gap-2">
    <Dropdown title="Device" bind:value={device.deviceId} loader={async () => {
      await listVideo()
      return videos
    }} class=""></Dropdown>
  </div>
</div>