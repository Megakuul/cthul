<script>
  import "../app.css";
  import { Palette } from "$lib/color/color.svelte";
  import sisyphos from "$lib/assets/sisyphos.svg";
</script>


<center
  style="background-color: {Palette().bgCthul()};"
  class="relative overflow-hidden"
>
  <img class="h-48" alt="" src={sisyphos} />
</center>
