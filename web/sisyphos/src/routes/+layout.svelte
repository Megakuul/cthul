<script>
  import "../app.css";
  import { Palette } from "$lib/color/color.svelte"
  import Footer from "./Footer.svelte";

  let { children } = $props();
</script>

<center
  style="color: {Palette().fgPrimary()}; background-color: {Palette().bgPrimary()};"
  class="h-full min-h-screen box-border transition-all duration-700"
>
  {@render children()}
</center>

<Footer></Footer>
