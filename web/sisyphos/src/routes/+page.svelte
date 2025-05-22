<script>
  import rune from "$lib/assets/rune.svg";
  import wave from "$lib/assets/wave.svg";
  import granit from "$lib/assets/granit.svg";
  import proton from "$lib/assets/proton.svg";
  import waveDomain from "$lib/assets/wave-domain.svg";
  import waveSerial from "$lib/assets/wave-serial.svg";
  import waveVideo from "$lib/assets/wave-video.svg";
  import granitDisk from "$lib/assets/granit-disk.svg";
  import click from "$lib/assets/click.mp3";
  import Card from "./Card.svelte";
  import { Palette } from "$lib/color/color.svelte";
  import { Play } from "$lib/sound/sound.svelte";
    import Link from "$lib/component/Link/Link.svelte";

  /** @type {"rune" | "wave" | "granit" | "proton" | undefined} */
  let selected = $state(undefined);

  /** @typedef service
   * @property {string} title
   * @property {string} icon
   * @property {string} route
   */

  /** @type {service[] | undefined} */
  let services = $state(undefined);

  $effect(() => {
    switch (selected) {
      case "rune":
        services = [];
        break;
      case "wave":
        services = [
          {
            title: "Domain",
            icon: waveDomain,
            route: "/wave/domain",
          },
          {
            title: "Serial",
            icon: waveSerial,
            route: "/wave/serial",
          },
          {
            title: "Video",
            icon: waveVideo,
            route: "/wave/video",
          },
        ];
        break;
      case "granit":
        services = [
          {
            title: "Disk",
            icon: granitDisk,
            route: "/granit/disk",
          },
        ];
        break;
      case "proton":
        services = [
          {
            title: "Inter",
            icon: "",
            route: "/proton/inter",
          },
        ];
        break;
      default:
        services = undefined;
    }
  });

</script>

<div class="flex flex-row justify-around gap-5 sm:gap-20 p-4">
  <Card icon={wave} name="wave" bind:selected={selected}></Card>
  <Card icon={granit} name="granit" bind:selected={selected}></Card>
  <Card icon={proton} name="proton" bind:selected={selected}></Card>
  <Card icon={rune} name="rune" bind:selected={selected}></Card>
</div>

<div class="w-full flex flex-col gap-6 px-4 py-10">
  {#if services && selected}
    <Link href="/{selected}/settings" class="flex flex-row justify-center items-center gap-2 p-4 rounded-lg text-2xl bg-slate-50/40">
      {#if selected === "wave"}
        <span style="color: {Palette().fgWave()}B3" class="font-bold">Wave</span>
      {:else if selected === "granit"}
        <span style="color: {Palette().fgGranit()}B3" class="font-bold">Granit</span>
      {:else if selected === "proton"}
        <span style="color: {Palette().fgProton()}B3" class="font-bold">Proton</span>
      {:else if selected === "rune"}
        <span style="color: {Palette().fgRune()}B3" class="font-bold">Rune</span>
      {/if}
      <span class="font-bold">Settings</span>
      <svg xmlns="http://www.w3.org/2000/svg" width="30" height="30" viewBox="0 0 24 24"><path fill="currentColor" d="m19.28 8.6l-.7-1.21l-1.27.51l-1.06.43l-.91-.7c-.39-.3-.8-.54-1.23-.71l-1.06-.43l-.16-1.13L12.7 4h-1.4l-.19 1.35l-.16 1.13l-1.06.44c-.41.17-.82.41-1.25.73l-.9.68l-1.05-.42l-1.27-.52l-.7 1.21l1.08.84l.89.7l-.14 1.13c-.03.3-.05.53-.05.73s.02.43.05.73l.14 1.13l-.89.7l-1.08.84l.7 1.21l1.27-.51l1.06-.43l.91.7c.39.3.8.54 1.23.71l1.06.43l.16 1.13l.19 1.36h1.39l.19-1.35l.16-1.13l1.06-.43c.41-.17.82-.41 1.25-.73l.9-.68l1.04.42l1.27.51l.7-1.21l-1.08-.84l-.89-.7l.14-1.13c.04-.31.05-.52.05-.73s-.02-.43-.05-.73l-.14-1.13l.89-.7zM12 16c-2.21 0-4-1.79-4-4s1.79-4 4-4s4 1.79 4 4s-1.79 4-4 4" opacity="0.3"/><path fill="currentColor" d="M19.43 12.98c.04-.32.07-.64.07-.98s-.03-.66-.07-.98l2.11-1.65c.19-.15.24-.42.12-.64l-2-3.46a.5.5 0 0 0-.61-.22l-2.49 1c-.52-.4-1.08-.73-1.69-.98l-.38-2.65A.49.49 0 0 0 14 2h-4c-.25 0-.46.18-.49.42l-.38 2.65c-.61.25-1.17.59-1.69.98l-2.49-1a.6.6 0 0 0-.18-.03c-.17 0-.34.09-.43.25l-2 3.46c-.13.22-.07.49.12.64l2.11 1.65c-.04.32-.07.65-.07.98s.03.66.07.98l-2.11 1.65c-.19.15-.24.42-.12.64l2 3.46a.5.5 0 0 0 .61.22l2.49-1c.52.4 1.08.73 1.69.98l.38 2.65c.03.24.24.42.49.42h4c.25 0 .46-.18.49-.42l.38-2.65c.61-.25 1.17-.59 1.69-.98l2.49 1q.09.03.18.03c.17 0 .34-.09.43-.25l2-3.46c.12-.22.07-.49-.12-.64zm-1.98-1.71c.04.31.05.52.05.73s-.02.43-.05.73l-.14 1.13l.89.7l1.08.84l-.7 1.21l-1.27-.51l-1.04-.42l-.9.68c-.43.32-.84.56-1.25.73l-1.06.43l-.16 1.13l-.2 1.35h-1.4l-.19-1.35l-.16-1.13l-1.06-.43c-.43-.18-.83-.41-1.23-.71l-.91-.7l-1.06.43l-1.27.51l-.7-1.21l1.08-.84l.89-.7l-.14-1.13c-.03-.31-.05-.54-.05-.74s.02-.43.05-.73l.14-1.13l-.89-.7l-1.08-.84l.7-1.21l1.27.51l1.04.42l.9-.68c.43-.32.84-.56 1.25-.73l1.06-.43l.16-1.13l.2-1.35h1.39l.19 1.35l.16 1.13l1.06.43c.43.18.83.41 1.23.71l.91.7l1.06-.43l1.27-.51l.7 1.21l-1.07.85l-.89.7zM12 8c-2.21 0-4 1.79-4 4s1.79 4 4 4s4-1.79 4-4s-1.79-4-4-4m0 6c-1.1 0-2-.9-2-2s.9-2 2-2s2 .9 2 2s-.9 2-2 2"/></svg>
    </Link>
    {#each services as service}
      <Link href={service.route} class="w-full h-24 p-4 flex flex-row justify-between items-center rounded-lg bg-slate-50/30">
        <h1 class="text-2xl font-bold">{service.title}</h1>
        <img alt="{service.title}" src={service.icon} class="h-full" />
      </Link>
    {/each}
  {:else}
    <Link href="/profile" class="flex flex-row justify-center items-center gap-2 p-4 rounded-lg text-2xl bg-slate-50/20">
      <span class="font-bold">Account Preferences</span>
      <svg xmlns="http://www.w3.org/2000/svg" width="30" height="30" viewBox="0 0 24 24"><path fill="currentColor" d="m19.28 8.6l-.7-1.21l-1.27.51l-1.06.43l-.91-.7c-.39-.3-.8-.54-1.23-.71l-1.06-.43l-.16-1.13L12.7 4h-1.4l-.19 1.35l-.16 1.13l-1.06.44c-.41.17-.82.41-1.25.73l-.9.68l-1.05-.42l-1.27-.52l-.7 1.21l1.08.84l.89.7l-.14 1.13c-.03.3-.05.53-.05.73s.02.43.05.73l.14 1.13l-.89.7l-1.08.84l.7 1.21l1.27-.51l1.06-.43l.91.7c.39.3.8.54 1.23.71l1.06.43l.16 1.13l.19 1.36h1.39l.19-1.35l.16-1.13l1.06-.43c.41-.17.82-.41 1.25-.73l.9-.68l1.04.42l1.27.51l.7-1.21l-1.08-.84l-.89-.7l.14-1.13c.04-.31.05-.52.05-.73s-.02-.43-.05-.73l-.14-1.13l.89-.7zM12 16c-2.21 0-4-1.79-4-4s1.79-4 4-4s4 1.79 4 4s-1.79 4-4 4" opacity="0.3"/><path fill="currentColor" d="M19.43 12.98c.04-.32.07-.64.07-.98s-.03-.66-.07-.98l2.11-1.65c.19-.15.24-.42.12-.64l-2-3.46a.5.5 0 0 0-.61-.22l-2.49 1c-.52-.4-1.08-.73-1.69-.98l-.38-2.65A.49.49 0 0 0 14 2h-4c-.25 0-.46.18-.49.42l-.38 2.65c-.61.25-1.17.59-1.69.98l-2.49-1a.6.6 0 0 0-.18-.03c-.17 0-.34.09-.43.25l-2 3.46c-.13.22-.07.49.12.64l2.11 1.65c-.04.32-.07.65-.07.98s.03.66.07.98l-2.11 1.65c-.19.15-.24.42-.12.64l2 3.46a.5.5 0 0 0 .61.22l2.49-1c.52.4 1.08.73 1.69.98l.38 2.65c.03.24.24.42.49.42h4c.25 0 .46-.18.49-.42l.38-2.65c.61-.25 1.17-.59 1.69-.98l2.49 1q.09.03.18.03c.17 0 .34-.09.43-.25l2-3.46c.12-.22.07-.49-.12-.64zm-1.98-1.71c.04.31.05.52.05.73s-.02.43-.05.73l-.14 1.13l.89.7l1.08.84l-.7 1.21l-1.27-.51l-1.04-.42l-.9.68c-.43.32-.84.56-1.25.73l-1.06.43l-.16 1.13l-.2 1.35h-1.4l-.19-1.35l-.16-1.13l-1.06-.43c-.43-.18-.83-.41-1.23-.71l-.91-.7l-1.06.43l-1.27.51l-.7-1.21l1.08-.84l.89-.7l-.14-1.13c-.03-.31-.05-.54-.05-.74s.02-.43.05-.73l.14-1.13l-.89-.7l-1.08-.84l.7-1.21l1.27.51l1.04.42l.9-.68c.43-.32.84-.56 1.25-.73l1.06-.43l.16-1.13l.2-1.35h1.39l.19 1.35l.16 1.13l1.06.43c.43.18.83.41 1.23.71l.91.7l1.06-.43l1.27-.51l.7 1.21l-1.07.85l-.89.7zM12 8c-2.21 0-4 1.79-4 4s1.79 4 4 4s4-1.79 4-4s-1.79-4-4-4m0 6c-1.1 0-2-.9-2-2s.9-2 2-2s2 .9 2 2s-.9 2-2 2"/></svg>
    </Link>
    <Link href="/help" class="flex flex-row justify-center items-center gap-2 p-4 rounded-lg text-2xl bg-slate-50/20">
      <span class="font-bold">The <span class="text-amber-800">Helper</span> Manifesto</span>
      <svg xmlns="http://www.w3.org/2000/svg" width="30" height="30" viewBox="0 0 22 22"><path fill="currentColor" d="M18 20H4v-1H3v-1H2V4h1V3h1V2h14v1h1v1h1v14h-1v1h-1m-6-6v-1h1v-1h1V7h-1V6H9v1H8v3h2V8h2v2h-1v1h-1v2m2 3v-2h-2v2Z"/></svg>
    </Link>
  {/if}
</div>

<center class="font-bold text-5xl pb-20 select-none">***</center>