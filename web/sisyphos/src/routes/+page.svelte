<script>
  import { Palette } from "$lib/color/color.svelte";
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
        SetPalette(NewRunePalette());
        break;
      case "wave":
        SetPalette(NewWavePalette());
        break;
      case "granit":
        SetPalette(NewGranitPalette());
        break;
      case "proton":
        SetPalette(NewProtonPalette());
        break;
      default:
        SetPalette(NewDefaultPalette());
    }
  });

  /** @param {"rune" | "wave" | "granit" | "proton" | undefined} newState */
  function select(newState) {
    if (selected === newState) {
      newState = undefined;
    }
    selected = newState;
  }
</script>

<div class="flex flex-row justify-around">
  <button
    id="rune"
    class:selected={selected === "rune"}
    onclick={() => select("rune")}
  >
    <img alt="rune" src={rune} />
  </button>
  <button
    id="wave"
    class:selected={selected === "wave"}
    onclick={() => select("wave")}
  >
    <img alt="wave" src={wave} />
  </button>
  <button
    id="granit"
    class:selected={selected === "granit"}
    onclick={() => select("granit")}
  >
    <img alt="granit" src={granit} />
  </button>
  <button
    id="proton"
    class:selected={selected === "proton"}
    onclick={() => select("proton")}
  >
    <img alt="proton" src={proton} />
  </button>
</div>

<div class="m-4 grid grid-cols-3">
  <div></div>
</div>
<a href="/domain">domain</a>
<p>
  Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the
  documentation
</p>

<style>
  #rune {
    cursor: pointer;
    transition: all ease .25s;
  }

  #rune:hover,
  #rune.selected {
    filter: drop-shadow(6px 6px 10px #5e4a11);
  }

  #wave {
    transition: all ease .25s;
  }

  #wave:hover,
  #wave.selected {
    filter: drop-shadow(6px 6px 10px  #0D65A4);
  }

  #granit {
    transition: all ease .25s;
  }

  #granit:hover,
  #granit.selected {
    filter: drop-shadow(6px 6px 10px  #042F0B);
  }

  #proton {
    transition: all ease .25s;
  }

  #proton:hover,
  #proton.selected {
    filter: drop-shadow(6px 6px 10px #57056C);
  }
</style>
