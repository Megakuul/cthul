<script>
  import "../app.css";
  import chainVertical from "$lib/assets/chain-v.svg";
  import chainHorizontal from "$lib/assets/chain-h.svg";
  import chainVerticalScattered from "$lib/assets/chain-v-scattered.svg";
  import chainHorizontalScattered from "$lib/assets/chain-h-scattered.svg";
  import { Palette } from "$lib/color/color.svelte";
  import Cthul from "$lib/components/Cthul.svelte";
</script>

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
    style="background-image: url({chainVertical});"
    class="chain vertical absolute -top-20 left-40 rotate-45 w-3 h-96 opacity-20 peer-hover:opacity-80 scale-[1.25] peer-hover:scale-[1.8]"
  ></div>
  <div
    style="background-image: url({chainVertical});"
    class="chain vertical absolute -top-20 left-0 -rotate-12 w-3 h-96 opacity-15 peer-hover:opacity-60 scale-[1] peer-hover:scale-[1.1]"
  ></div>
  <div
    style="background-image: url({chainVertical});"
    class="chain vertical absolute top-0 left-30 -rotate-45 w-3 h-96 opacity-25 peer-hover:opacity-75 peer-hover:scale-[1.6]"
  ></div>
  <div
    style="background-image: url({chainVerticalScattered});"
    class="chain vertical absolute top-0 left-40 -rotate-12 w-3 h-96 opacity-10 peer-hover:opacity-50 scale-[1.5] peer-hover:scale-[2.2]"
  ></div>
  <div
    style="background-image: url({chainVertical});"
    class="chain vertical absolute top-[600px] peer-hover:top-0 left-60 w-3 h-96 scale-[2.5] peer-hover:scale-[3]"
  ></div>
  <div
    style="background-image: url({chainHorizontalScattered});"
    class="chain horizontal absolute top-55 left-0 w-0 peer-hover:w-screen h-3 rotate-[12deg] scale-[2.5]"
  ></div>

  <div
    style="background-image: url({chainVerticalScattered});"
    class="chain vertical absolute -top-20 right-40 rotate-45 w-3 h-96 opacity-20 peer-hover:opacity-50 scale-[1.25] peer-hover:scale-[1.8]"
  ></div>
  <div
    style="background-image: url({chainVerticalScattered});"
    class="chain vertical absolute -top-20 right-0 -rotate-12 w-3 h-96 opacity-15 peer-hover:opacity-30 scale-[1] peer-hover:scale-[1.1]"
  ></div>
  <div
    style="background-image: url({chainVertical});"
    class="chain vertical absolute -top-20 right-30 -rotate-45 w-3 h-96 opacity-15 peer-hover:opacity-85 peer-hover:scale-[1.6]"
  ></div>
  <div
    style="background-image: url({chainVertical});"
    class="chain vertical absolute top-0 right-40 -rotate-12 w-3 h-96 opacity-10 peer-hover:opacity-90 scale-[1.5] peer-hover:scale-[2.2]"
  ></div>
  <div
    style="background-image: url({chainVertical});"
    class="chain vertical absolute top-[600px] peer-hover:top-0 right-60 w-3 h-96 scale-[2.5] peer-hover:scale-[3]"
  ></div>
  <div
    style="background-image: url({chainHorizontal});"
    class="chain horizontal absolute top-55 right-0 w-0 peer-hover:w-screen h-3 -rotate-[10deg] scale-[2.5]"
  ></div>
</center>

<style>
  .chain {
    @apply visible;
    @media (max-width: 64rem) {
      @apply invisible;
    }
  }

  .chain.vertical {
    @apply bg-repeat-y bg-[length:100%_auto] bg-[position:0_1px] peer-hover:bg-[position:0_-150px] transition-all duration-1000;
  }

  .chain.horizontal {
    @apply bg-repeat-x bg-[length:auto_100%] bg-[position:1_0px] peer-hover:bg-[position:150px_0] transition-all duration-1000;
  }
</style>
