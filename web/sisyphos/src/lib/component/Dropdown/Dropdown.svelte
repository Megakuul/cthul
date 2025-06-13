<!-- Dynamic unvalidated dropdown menu -->
<script>
  import { Palette } from "$lib/color/color.svelte";
  import { cn } from "$lib/utils";
  import { fade } from "svelte/transition";
  import { flip } from "svelte/animate";

  let {
    value = $bindable(),
    loader,
    title,
    class: classNames,
  } = $props()

  /** @type {boolean} */
  let selected = $state(false)

  let items = $derived.by(async () => {
    console.log(selected)
    if (selected) {
      return await loader()
    }
    return []
  })

  /** @type {string[]} */
  let sortedItems = $derived.by(() => {
    return Object.keys(items).filter((item) => {
      return item.toUpperCase().includes(value.toUpperCase())
    })
  })
</script>

<div class={cn("relative cursor-pointer transition-all", classNames)} onfocusout={(/** @type {any} */ e) => {
    if (!e.currentTarget.contains(e.relatedTarget)) {
      selected = false
    }
  }}>
  <span class="absolute -top-2 left-1 text-xs font-bold">{title}</span>
  <input placeholder={title} bind:value={value}
    onfocus={() => selected = true}
    class="text-xl w-full p-1 rounded-md focus:outline-0 transition-all overflow-hidden bg-slate-50/10 focus:bg-slate-50/20"
  />
  {#if selected}
    <div style="background-color: {Palette().bgPrimary()};"
      class="absolute top-10 left-0 right-0 z-10 max-h-24 flex flex-col gap-1 rounded-lg overflow-scroll-hidden">
      {#each sortedItems as item (item)}
        <button animate:flip onfocus={() => selected = true} onclick={() => {
          value = item
          selected = false
        }} class="w-full text-start p-2 cursor-pointer focus:bg-slate-950/10 hover:bg-slate-950/10 focus:outline-0 transition-all">{item}</button>
      {/each}
    </div>
  {/if}
</div>