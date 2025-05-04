<script>
  import { Palette } from "$lib/color/color.svelte";
  import stone from "$lib/assets/stone.ogg";
  import { Play } from "$lib/sound/sound.svelte";

  let {
    icon,
    name,
    selected = $bindable()
  } = $props();
</script>

<button aria-label="{name}" class:selected={selected === name} style="background-color: {Palette().bgCthul()}; --main-bg: {Palette().bgPrimary()};"
  class="w-full relative select-none rounded-4xl cursor-pointer overflow-hidden shadow-inner shadow-slate-100/40"
  onclick={() => {
    Play(stone)
    if (selected === name) {
      selected = undefined;
    } else {
      selected = name;
    }
  }}
  >
  <img alt="{name}" src={icon} class="transition-all select-none duration-500" />
</button>

<style>
  button img {
    filter: drop-shadow(0px 0px 0px rgba(255, 255,255, 0.3)) brightness(60%);
  }

  button.selected img {
    scale: 1.05;
    filter: drop-shadow(0px 0px 40px rgba(255, 255,255, 0.3)) brightness(120%);
  }
</style>
