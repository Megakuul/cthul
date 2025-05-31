<script>
  import { Palette } from "$lib/color/color.svelte";
  import { cn } from "$lib/utils";
  import { fade } from "svelte/transition";
  import { flip } from "svelte/animate";

  let {
    value = $bindable(),
    loader,
    placeholder,
    class: classNames,
  } = $props()

  let selected = $state(false)

  /** @type {string} */
  let input = $state("")

  /** @type {{[key: string]: number}} */
  let items = $derived.by(() => {
    if (selected) {
      return loader()
    }
  })

  /** @type {string[]} */
  let sortedItems = $derived.by(() => {
    return Object.keys(items).filter((item) => {
      return item.toUpperCase().includes(input.toUpperCase())
    })
  })

  $effect(() => {
    if (items) {
      value = items[input]
    }
  })
</script>

<div class={cn("relative cursor-pointer transition-all", classNames)} onfocusout={(/** @type {any} */ e) => {
    if (!e.currentTarget.contains(e.relatedTarget)) {
      selected = false
    }
  }}>
  <input placeholder={placeholder} bind:value={input}
    onfocus={() => selected = true}
    class="text-xl w-full p-1 rounded-md bg-slate-50/10 focus:bg-slate-50/20 focus:outline-0 transition-all overflow-hidden" />
  {#if selected}
    <div transition:fade style="background-color: {Palette().bgPrimary()};"
      class="absolute top-10 left-0 right-0 z-10 max-h-24 flex flex-col gap-1 rounded-lg overflow-scroll-hidden">
      {#each sortedItems as item (item)}
        <button animate:flip onfocus={() => selected = true} onclick={() => {
          input = item
          selected = false
        }} class="w-full text-start p-2 cursor-pointer focus:bg-slate-950/10 hover:bg-slate-950/10 focus:outline-0 transition-all">{item}</button>
      {/each}
    </div>
  {/if}
</div>