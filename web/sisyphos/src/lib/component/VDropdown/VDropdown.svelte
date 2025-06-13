<!-- Static validated dropdown menu -->
<script>
  import { Palette } from "$lib/color/color.svelte";
  import { cn } from "$lib/utils";
  import { fade } from "svelte/transition";
  import { flip } from "svelte/animate";

  let {
    value = $bindable(),
    items,
    title,
    class: classNames,
  } = $props()

  const baseValidationClass = ""
  let validationClass = $state(baseValidationClass)

  /** @type {boolean} */
  let selected = $state(false)

  /** @type {string} */
  let input = $state(getOption(value))

  /** @type {string[]} */
  let sortedItems = $state([])

  /**
   * @param {any} optionValue 
   * @returns {string} the option string for the current value 
  */
  function getOption(optionValue) {
    return Object.keys(items).filter((k) => {
      if (items[k] === optionValue) return true;
    })[0] || ""
  }

  /** @param {string} newInput  */
  function handleInput(newInput) {
    if (items) {
      sortedItems = Object.keys(items).filter((item) => {
        return item.toUpperCase().includes(newInput.toUpperCase())
      })
    }
    validationClass = baseValidationClass
  }

  /** @param {string} newValue  */
  function update(newValue) {
    if (items && items.hasOwnProperty(newValue)) {
      value = items[newValue]
      input = newValue
      validationClass = "";
    } else {
      validationClass = "text-red-800/90";
    }
  }
</script>

<div class={cn("relative cursor-pointer transition-all", classNames)} onfocusout={(/** @type {any} */ e) => {
    if (!e.currentTarget.contains(e.relatedTarget)) {
      selected = false
      update(input)
    }
  }}>
  <span class="absolute -top-2 left-1 text-xs font-bold">{title}</span>
  <input placeholder={title} bind:value={input}
    onfocus={() => selected = true}
    oninput={(/** @type {any} */ e) => handleInput(e.target?.value)}
    class={cn("text-xl w-full p-1 rounded-md focus:outline-0 transition-all overflow-hidden bg-slate-50/10 focus:bg-slate-50/20", validationClass)} 
  />
  {#if selected}
    <div style="background-color: {Palette().bgPrimary()};"
      class="absolute top-11 left-0 right-0 z-10 max-h-24 flex flex-col gap-1 rounded-lg overflow-scroll-hidden">
      {#each sortedItems as item (item)}
        <button animate:flip onfocus={() => selected = true} onclick={() => {
          update(item)
          selected = false
        }} class="w-full text-start p-2 cursor-pointer focus:bg-slate-950/10 hover:bg-slate-950/10 focus:outline-0 transition-all">{item}</button>
      {/each}
    </div>
  {/if}
</div>