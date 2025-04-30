<script>
    import { Palette } from "$lib/color/color.svelte";
    import chainVertical from "$lib/assets/chain-v.svg";
  import chainHorizontal from "$lib/assets/chain-h.svg";
  import chainVerticalScattered from "$lib/assets/chain-v-scattered.svg";
  import chainHorizontalScattered from "$lib/assets/chain-h-scattered.svg";
    import { Gimmick } from "$lib/gimmick/gimmick.svelte";

  let {
    icon,
    name,
    selected = $bindable()
  } = $props();
</script>

<button aria-label="{name}" class:selected={selected === name} style="background-color: {Palette().bgCthul()}; --main-bg: {Palette().bgPrimary()};"
  class="w-full cursor-pointer overflow-hidden shadow-inner shadow-slate-100/40"
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
      class="chain absolute top-0 left-1/2 -rotate-45 w-3 h-96 opacity-80"
    ></div>
    <div
      style="background-image: url({chainVerticalScattered});"
      class="chain absolute top-0 left-1/2 rotate-45 w-3 h-96 opacity-80"
    ></div>
  {/if}
</button>

<style>
  button {
    position: relative;
  }
  button:before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 2;
    animation: open 1s linear forwards;
  }

  @keyframes open {
    0% {
      box-shadow: inset 0 0 0 200px var(--main-bg);
    }
    100% {
      box-shadow: inset 0 0 0 0 var(--main-bg);
    }
  }

  button img {
    filter: drop-shadow(0px 0px 0px rgba(255, 255,255, 0.3)) brightness(60%);
  }

  button.selected img {
    filter: drop-shadow(0px 0px 10px rgba(255, 255,255, 0.3)) brightness(120%);
  }

  .chain {
    position: absolute;
    top: 0px;
    
    visibility: hidden;
    background-repeat: repeat-y;
    background-size: 100% 50%;
    animation: unleash 2s linear forwards;
  }

  @keyframes unleash {
    0% {
      visibility: visible;
      scale: 5;
      background-position: 0 50%;
    }
    100% {
      visibility: visible;
      scale: 7;
      background-position: 0 0;
    }
  }
</style>
