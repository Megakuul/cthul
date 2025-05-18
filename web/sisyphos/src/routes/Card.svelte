<script>
  import stoneForward from "$lib/assets/stone_forward.mp3";
  import stoneBackward from "$lib/assets/stone_backward.mp3";
  import { Play } from "$lib/sound/sound.svelte";

  let {
    icon,
    name,
    selected = $bindable()
  } = $props();
</script>

<button aria-label="{name}" class:selected={selected === name}
  class="w-full relative select-none rounded-4xl bg-black cursor-pointer overflow-hidden"
  onclick={() => {
    if (selected === name) {
      Play(stoneBackward)
      selected = undefined;
    } else {
      Play(stoneForward)
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
