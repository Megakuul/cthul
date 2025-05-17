<script>
  import "../app.css";
  import { NewSandstormPalette, NewSlatePalette, Palette, SetPalette } from "$lib/color/color.svelte"
  import Footer from "./Footer.svelte";
  import { browser } from "$app/environment";
  import { Mute, Unmute } from "$lib/sound/sound.svelte";

  let { children } = $props();

  let Ready = $state(false)

  let Dark = $state(false)

  let Muted = $state(false)

  $effect.root(() => {
    if (browser) {
      Dark = localStorage.getItem("dark")?.toLowerCase() === "true";
      Dark ? SetPalette(NewSlatePalette()) : SetPalette(NewSandstormPalette())
    }
    Ready = true; // avoid FOUC
  })

  $effect(() => {
    localStorage.setItem("dark", Dark.toString())
    Dark ? SetPalette(NewSlatePalette()) : SetPalette(NewSandstormPalette())
  })

  $effect.root(() => {
    if (browser) {
      Muted = localStorage.getItem("muted")?.toLowerCase() === "true";
      Muted ? Mute() : Unmute()
    }
  })

  $effect(() => {
    localStorage.setItem("muted", Muted.toString())
    Muted ? Mute() : Unmute()
  })
</script>

{#if Ready}
  <center
  style="color: {Palette().fgPrimary()}; background-color: {Palette().bgPrimary()};"
  class="min-h-screen box-border transition-all duration-700"
  >
  <div class="w-full h-8 px-4 py-1 flex justify-start gap-2">
    <button onclick={() => Muted = !Muted} title="{Muted ? "unmute" : "mute"}"
      class="flex justify-center items-center cursor-pointer w-12 rounded-lg bg-slate-50/40">
      {#if Muted}
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path fill="currentColor" fill-opacity="0" stroke-dasharray="32" stroke-dashoffset="32" d="M4 10h3.5l3.5 -3.5v10.5l-3.5 -3.5h-3.5Z"><animate fill="freeze" attributeName="fill-opacity" begin="0.8s" dur="0.15s" values="0;0.3"/><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.4s" values="32;0"/></path><path stroke-dasharray="8" stroke-dashoffset="8" d="M16 10l4 4"><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.4s" dur="0.2s" values="8;0"/></path><path stroke-dasharray="8" stroke-dashoffset="8" d="M20 10l-4 4"><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.6s" dur="0.2s" values="8;0"/></path></g></svg>
      {:else}
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M16 12c0 0 0 0 0 0c0 0 0 0 0 0Z"><animate fill="freeze" attributeName="d" begin="0.4s" dur="0.2s" values="M16 12c0 0 0 0 0 0c0 0 0 0 0 0Z;M16 16c1.5 -0.71 2.5 -2.24 2.5 -4c0 -1.77 -1 -3.26 -2.5 -4Z"/></path><path fill="currentColor" fill-opacity="0" stroke="currentColor" stroke-dasharray="32" stroke-dashoffset="32" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 10h3.5l3.5 -3.5v10.5l-3.5 -3.5h-3.5Z"><animate fill="freeze" attributeName="fill-opacity" begin="0.6s" dur="0.15s" values="0;0.3"/><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.4s" values="32;0"/></path></svg>
      {/if}
    </button>
    <button onclick={() => Dark = !Dark} title="change theme"
      class="flex justify-center items-center cursor-pointer w-12 rounded-lg bg-slate-50/40">
      {#if Dark}
        <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24"><g fill="currentColor" fill-opacity="0"><path d="M15.22 6.03l2.53-1.94L14.56 4L13.5 1l-1.06 3l-3.19.09l2.53 1.94l-.91 3.06l2.63-1.81l2.63 1.81z"><animate id="lineMdMoonRisingLoop0" fill="freeze" attributeName="fill-opacity" begin="0.7s;lineMdMoonRisingLoop0.begin+6s" dur="0.4s" values="0;1"/><animate fill="freeze" attributeName="fill-opacity" begin="lineMdMoonRisingLoop0.begin+2.2s" dur="0.4s" values="1;0"/></path><path d="M13.61 5.25L15.25 4l-2.06-.05L12.5 2l-.69 1.95L9.75 4l1.64 1.25l-.59 1.98l1.7-1.17l1.7 1.17z"><animate fill="freeze" attributeName="fill-opacity" begin="lineMdMoonRisingLoop0.begin+3s" dur="0.4s" values="0;1"/><animate fill="freeze" attributeName="fill-opacity" begin="lineMdMoonRisingLoop0.begin+5.2s" dur="0.4s" values="1;0"/></path><path d="M19.61 12.25L21.25 11l-2.06-.05L18.5 9l-.69 1.95l-2.06.05l1.64 1.25l-.59 1.98l1.7-1.17l1.7 1.17z"><animate fill="freeze" attributeName="fill-opacity" begin="lineMdMoonRisingLoop0.begin+0.4s" dur="0.4s" values="0;1"/><animate fill="freeze" attributeName="fill-opacity" begin="lineMdMoonRisingLoop0.begin+2.8s" dur="0.4s" values="1;0"/></path><path d="M20.828 9.731l1.876-1.439l-2.366-.067L19.552 6l-.786 2.225l-2.366.067l1.876 1.439L17.601 12l1.951-1.342L21.503 12z"><animate fill="freeze" attributeName="fill-opacity" begin="lineMdMoonRisingLoop0.begin+3.4s" dur="0.4s" values="0;1"/><animate fill="freeze" attributeName="fill-opacity" begin="lineMdMoonRisingLoop0.begin+5.6s" dur="0.4s" values="1;0"/></path></g><path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" d="M7 6 C7 12.08 11.92 17 18 17 C18.53 17 19.05 16.96 19.56 16.89 C17.95 19.36 15.17 21 12 21 C7.03 21 3 16.97 3 12 C3 8.83 4.64 6.05 7.11 4.44 C7.04 4.95 7 5.47 7 6 Z" transform="translate(0 22)" stroke-width="1"><animateMotion fill="freeze" calcMode="linear" dur="0.6s" path="M0 0v-22"/></path></svg>
      {:else}
        <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="2" stroke-dashoffset="2" d="M12 19v1M19 12h1M12 5v-1M5 12h-1"><animate fill="freeze" attributeName="d" begin="0.6s" dur="0.2s" values="M12 19v1M19 12h1M12 5v-1M5 12h-1;M12 21v1M21 12h1M12 3v-1M3 12h-1"/><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.6s" dur="0.2s" values="2;0"/></path><path stroke-dasharray="2" stroke-dashoffset="2" d="M17 17l0.5 0.5M17 7l0.5 -0.5M7 7l-0.5 -0.5M7 17l-0.5 0.5"><animate fill="freeze" attributeName="d" begin="0.8s" dur="0.2s" values="M17 17l0.5 0.5M17 7l0.5 -0.5M7 7l-0.5 -0.5M7 17l-0.5 0.5;M18.5 18.5l0.5 0.5M18.5 5.5l0.5 -0.5M5.5 5.5l-0.5 -0.5M5.5 18.5l-0.5 0.5"/><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.8s" dur="0.2s" values="2;0"/></path><animateTransform attributeName="transform" dur="30s" repeatCount="indefinite" type="rotate" values="0 12 12;360 12 12"/></g><mask id="lineMdMoonAltToSunnyOutlineLoopTransition0"><circle cx="12" cy="12" r="12" fill="#fff"/><circle cx="12" cy="12" r="8"><animate fill="freeze" attributeName="r" dur="0.4s" values="8;4"/></circle><circle cx="18" cy="6" r="12" fill="#fff"><animate fill="freeze" attributeName="cx" dur="0.4s" values="18;22"/><animate fill="freeze" attributeName="cy" dur="0.4s" values="6;2"/><animate fill="freeze" attributeName="r" dur="0.4s" values="12;3"/></circle><circle cx="18" cy="6" r="10"><animate fill="freeze" attributeName="cx" dur="0.4s" values="18;22"/><animate fill="freeze" attributeName="cy" dur="0.4s" values="6;2"/><animate fill="freeze" attributeName="r" dur="0.4s" values="10;1"/></circle></mask><circle cx="12" cy="12" r="10" mask="url(#lineMdMoonAltToSunnyOutlineLoopTransition0)" fill="currentColor"><animate fill="freeze" attributeName="r" dur="0.4s" values="10;6"/></circle></svg>
      {/if}
    </button>
    <a href="https://github.com/megakuul/cthul" class="flex justify-center items-center ml-auto cursor-pointer hover:underline">
      <span class="text-xs font-bold">Cthul Sisyphos ©</span>
    </a>
  </div>
  {@render children()}
  </center>

  <Footer></Footer>
{/if}
