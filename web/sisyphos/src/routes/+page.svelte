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

  /** @typedef service
   * @property {string} title
   * @property {string} desc
   * @property {string} route
   */

  /** @type {service[] | undefined} */
  let services = $state(undefined);

  $effect(() => {
    switch (selected) {
      case "rune":
        SetPalette(NewRunePalette());
        services = [];
        break;
      case "wave":
        SetPalette(NewWavePalette());
        services = [
          {
            title: "Domain",
            desc: "Manage virtual machines",
            route: "/wave/domain",
          },
          {
            title: "Serial",
            desc: "Configure serial devices",
            route: "/wave/serial",
          },
          {
            title: "Video",
            desc: "Adjust video devices",
            route: "/wave/video",
          },
        ];
        break;
      case "granit":
        SetPalette(NewGranitPalette());
        services = [
          {
            title: "Disk",
            desc: "Manage replicated disks",
            route: "/granit/disk",
          },
        ];
        break;
      case "proton":
        SetPalette(NewProtonPalette());
        services = [
          {
            title: "Inter",
            desc: "Manage network interfaces",
            route: "/proton/inter",
          },
        ];
        break;
      default:
        services = undefined;
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

<div class="flex flex-row justify-around p-4">
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

{#if services}
  <div class="grid grid-cols-3 gap-2 p-10">
    {#each services as service}
      <a href="{service.route}" class="w-64 rounded-2xl border-2">
        <h1>{service.title}</h1>
        <p>{service.desc}</p>
      </a>  
    {/each}
  </div>
{:else}
  <div>Here is some general content</div>
{/if}

<style>
  #rune {
    cursor: pointer;
    transition: all ease 0.25s;
  }

  #rune:hover,
  #rune.selected {
    filter: drop-shadow(6px 6px 10px #5e4a11);
  }

  #wave {
    cursor: pointer;
    transition: all ease 0.25s;
  }

  #wave:hover,
  #wave.selected {
    filter: drop-shadow(6px 6px 10px #0d65a4);
  }

  #granit {
    cursor: pointer;
    transition: all ease 0.25s;
  }

  #granit:hover,
  #granit.selected {
    filter: drop-shadow(6px 6px 10px #042f0b);
  }

  #proton {
    cursor: pointer;
    transition: all ease 0.25s;
  }

  #proton:hover,
  #proton.selected {
    filter: drop-shadow(6px 6px 10px #57056c);
  }
</style>
