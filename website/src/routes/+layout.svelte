<script>
  import "../app.css"
  import Footer from "./Footer.svelte";

  let { children } = $props();
</script>

<main class="overflow-hidden flex flex-col items-center">
  {@render children()}
</main>

<Footer></Footer>