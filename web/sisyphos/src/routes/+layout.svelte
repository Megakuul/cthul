<script>
  import "../app.css";
  import chain from "$lib/assets/chain.svg";
  import { Palette } from "$lib/color/color.svelte";
  import Cthul from "$lib/components/Cthul.svelte";

  let { children } = $props();
</script>

<center
  style="background-color: {Palette.bgPrimary()};"
  class="h-full min-h-screen box-border"
>
  {@render children()}
</center>
<center
  style="background-color: {Palette.bgSecondary()};"
  class="relative h-48 overflow-hidden"
>
  <Cthul
    eyeClass="animate-breath"
    bodyClass="brightness-20 group-hover:brightness-80 transition-all duration-1000"
    class="peer group cursor-pointer h-full w-max mt-20 hover:mt-10 hover:scale-125 transition-all duration-1000"
  ></Cthul>
  <div
    style="background-image: url({chain})"
    class="absolute top-0 left-20 -rotate-45 w-3 h-full chain"
  ></div>
  <div
    style="background-image: url({chain})"
    class="absolute -top-20 right-60 rotate-45 w-3 h-96 chain"
  ></div>
</center>

<style>
  .chain {
    @apply bg-repeat-y bg-[length:100%_auto] bg-[position:0_1px] peer-hover:bg-[position:0_-150px] transition-all duration-1000;
  }
</style>
