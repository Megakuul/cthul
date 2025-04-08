<script>
  import rune from "$lib/assets/rune.svg";
  import wave from "$lib/assets/wave.svg";
  import granit from "$lib/assets/granit.svg";
  import proton from "$lib/assets/proton.svg";

  import {
    SetPalette,
    NewRunePalette,
    NewWavePalette,
    NewGranitPalette,
    NewProtonPalette,
    NewDefaultPalette,
  } from "$lib/color/color.svelte";

  /** @type {"rune" | "wave" | "granit" | "proton" | undefined} */
  let selected = $state(undefined);

  $effect(() => {
    switch (selected) {
      case "rune":
        SetPalette(NewRunePalette())
        break;
      case "wave":
        SetPalette(NewWavePalette())
        break;
      case "granit":
        SetPalette(NewGranitPalette())
        break;
      case "proton":
        SetPalette(NewProtonPalette())
        break;
      default:
        SetPalette(NewDefaultPalette())
    }
  });
</script>

<div class="flex flex-row justify-around">
  <button onclick={() => selected = "rune"}>
    <img alt="rune" src={rune}>
  </button>
  <button onclick={() => selected = "wave"}>
    <img alt="wave" src={wave}>
  </button>
  <button onclick={() => selected = "granit"}>
    <img alt="granit" src={granit}>
  </button>
  <button onclick={() => selected = "proton"}>
    <img alt="proton" src={proton}>
  </button>
</div>

<div class="m-4 grid grid-cols-3">
  <div>

  </div>
</div>
<a href="/domain">domain</a>
<p>
  Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the
  documentation
</p>
