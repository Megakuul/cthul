<script lang="ts">
  import { type Video } from "$lib/sdk/types/wave/v1/video/message_pb";
  import { page } from "$app/state";
  import { DomainClient } from "$lib/client/client.svelte";
  import { ListRequestSchema, type Domain } from "$lib/sdk/types/wave/v1/domain/message_pb";
  import { create } from "@bufbuild/protobuf";
  import { SetException } from "$lib/exception/exception.svelte";
  import Link from "$lib/component/Link/Link.svelte";
  import { fade } from "svelte/transition";

  import waveDomain from "$lib/assets/wave-domain.svg"

  let {
    video = $bindable(),
  }: {
    video: Video;
  } = $props();

  let attached: {id: string, domain: Domain} | undefined = $state(undefined)

  async function checkAttached() {
    try {
      const request = create(ListRequestSchema, {});

      const response = await DomainClient().list(request)
      for (const [id, domain] of Object.entries(response.domains)) {
        for (const device of domain.config?.videoAdapters ?? []) {
          if (device.deviceId === page.params.id) {
            attached = {id: id, domain: domain}
            return
          }
        }
      }
    } catch (err: any) {
      SetException({title: "CHECK ATTACHED DOMAINS", desc: err.message})
    }
  }

  $effect.root(() => {
    if (page.params.id !== "new") checkAttached();
  })
</script>

<div class="w-full h-[500px] flex flex-row justify-between rounded-xl bg-slate-950/20">
  <div class="w-full flex flex-col items-start p-4">
    <h1 class="text-2xl">Name: <span class="font-bold">{video.config?.name}</span></h1>
    <p class="opacity-70">{page.params.id}</p>
    <hr class="w-full my-4">
    <div class="w-full flex flex-row gap-8 justify-start">
      <h2 class="flex flex-col items-start">
        <span class="text-sm font-bold">Node</span>
        <span class="text-2xl font-bold opacity-50">{video.node !== "" ? video.node : "<none>"}</span>
      </h2>
      <h2 class="flex flex-col items-start">
        <span class="text-sm font-bold">Reqnode</span>
        <span class="text-2xl font-bold opacity-50">{video.reqnode !== "" ? video.reqnode : "<none>"}</span>
      </h2>
    </div>
  </div>
  <div class="w-full flex flex-col gap-4 justify-center items-center p-4">
    {#if attached}
      <Link href="/wave/domain/{attached.id}" scale={1} class="w-2/3 bg-slate-50/20 rounded-xl">
        <img transition:fade width="60%" alt="" src={waveDomain}>
      </Link>
      <h1 transition:fade class="text-xl font-bold p-1 opacity-80 w-2/3 bg-slate-50/10 rounded-lg">device is attached</h1>
    {:else}
      <div class="w-2/3 bg-slate-50/20 rounded-xl">
        <img transition:fade width="60%" class="opacity-80 grayscale-100" alt="" src={waveDomain}>
      </div>
      <h1 transition:fade class="text-xl font-bold p-1 opacity-80 w-2/3 bg-slate-50/10 rounded-lg">device is detached</h1>
    {/if}
  </div>
</div>