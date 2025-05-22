<script>
  import { Play } from "$lib/sound/sound.svelte";
  import click from "$lib/assets/click.mp3";
  import { Palette } from "$lib/color/color.svelte";
  import { cn } from "$lib/utils";

  let {
    href,
    class: classNames,
    scale = 1.0,
    children,
    ...restProps
  } = $props()
</script>

<a href="{href}" onclick="{async () => Play(click)}" 
  style="border-color: {Palette().fgPrimary()}; --shadow: {Palette().btnShadow()} {scale * 5}px {scale * 5}px 0px;" 
  class={cn("btn cursor-pointer transition-all", classNames)} {...restProps}>
  {@render children()}
</a>

<style>
  .btn:hover {
    box-shadow: var(--shadow);
    transform: translateX(-3px) translateY(-1px);
  }
</style>