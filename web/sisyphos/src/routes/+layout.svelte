<script>
  import "../app.css";
  import { Palette } from "$lib/color/color.svelte"
  import Footer from "./Footer.svelte";
  import { browser } from "$app/environment";
  import { Mute, Unmute } from "$lib/sound/sound.svelte";

  let { children } = $props();

  let Muted = $state(false)

  $effect.root(() => {
    if (browser) {
      Muted = Boolean(localStorage.getItem("muted") ?? false);
    }
  })

  $effect(() => {
    localStorage.setItem("muted", Muted.toString())

    if (Muted) {
      Mute()
    } else {
      Unmute()
    }
  })
</script>

<center
  style="color: {Palette().fgPrimary()}; background-color: {Palette().bgPrimary()};"
  class="min-h-screen box-border transition-all duration-700"
>
  <button onclick={() => Muted = !Muted}>
    {#if Muted}
      Muted
    {:else}
      Unmuted
    {/if}
  </button>
  {@render children()}
</center>

<Footer></Footer>