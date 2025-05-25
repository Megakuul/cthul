<script>
  import { Play } from "$lib/sound/sound.svelte";
  import click from "$lib/assets/click.mp3";
  import { cn } from "$lib/utils";

  let {
    value,
    group = $bindable(),
    class: classNames,
    selectedClass,
    children,
  } = $props()
</script>

<input onclick={() => Play(click)} type="radio" id={value} bind:group={group} value={value} class="hidden">
<label for={value} class={cn("cursor-pointer transition-all", classNames, group === value ? selectedClass : "")}>
  {@render children()}
</label>