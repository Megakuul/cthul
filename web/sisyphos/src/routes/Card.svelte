<script>
  import { Palette } from "$lib/color/color.svelte";
  import chainVerticalScattered from "$lib/assets/chain-v-scattered.svg";
  import { Gimmick } from "$lib/gimmick/gimmick.svelte";

  let {
    icon,
    name,
    selected = $bindable()
  } = $props();
</script>

<button aria-label="{name}" class:selected={selected === name} style="background-color: {Palette().bgCthul()}; --main-bg: {Palette().bgPrimary()};"
  class="w-full relative select-none rounded-4xl cursor-pointer overflow-hidden shadow-inner shadow-slate-100/40"
  onclick={() => {
    if (selected === name) {
      selected = undefined;
    } else {
      selected = name;
    }
  }}
  >
  <img alt="{name}" src={icon} class="transition-all duration-500" />

  {#if Gimmick()}
    <div
      style="background-image: url({chainVerticalScattered});"
      class="chain absolute top-0 left-1/2 -rotate-45 w-3 opacity-80 scale-[7]"
    ></div>
    <div
      style="background-image: url({chainVerticalScattered});"
      class="chain absolute top-0 left-1/2 rotate-45 w-3 opacity-80 scale-[7]"
    ></div>
  {/if}
</button>

<style>
  button img {
    filter: drop-shadow(0px 0px 0px rgba(255, 255,255, 0.3)) brightness(60%);
  }

  button.selected img {
    filter: drop-shadow(0px 0px 10px rgba(255, 255,255, 0.3)) brightness(120%);
  }

  .chain {
    visibility: hidden;
    height: 100%;
    background-repeat: no-repeat;
    background-size: 100% 50%;
    animation: unleash 2s linear forwards;
  }

  @keyframes unleash {
    0% {
      visibility: visible;
      background-position: 0 60%;
    }
    99% {
      visibility: visible;
      background-position: 0 -20%;
    }
    100% {
      visibility: hidden;
      background-position: 0 0;
    }
  }
</style>
