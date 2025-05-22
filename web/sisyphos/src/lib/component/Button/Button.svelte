<script>
  import { Play } from "$lib/sound/sound.svelte";
  import click from "$lib/assets/click.mp3";
  import { Palette } from "$lib/color/color.svelte";
  import { cn } from "$lib/utils";

  let {
    onclick,
    class: classNames,
    scale = 1.0,
    children
  } = $props()
</script>

<button onclick={async () => {
    Play(click);
    onclick();
  }} style="border-color: {Palette().fgPrimary()}; --shadow: {Palette().btnShadow()} {scale * 5}px {scale * 5}px 0px;" 
  class={cn("btn cursor-pointer transition-all", classNames)}>
  {@render children()}
</button>

<style>
  .btn:hover {
    box-shadow: var(--shadow);
    transform: translateX(-3px) translateY(-1px);
  }
</style>