<script>
  import rune from "$lib/assets/rune.svg";
  import wave from "$lib/assets/wave.svg";
  import granit from "$lib/assets/granit.svg";
  import proton from "$lib/assets/proton.svg";
  import Card from "./Card.svelte";

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
        services = [];
        break;
      case "wave":
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
        services = [
          {
            title: "Disk",
            desc: "Manage replicated disks",
            route: "/granit/disk",
          },
        ];
        break;
      case "proton":
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
    }
  });

</script>

<div class="flex flex-row justify-around gap-6 p-4">
  <Card icon={wave} name="wave" bind:selected={selected}></Card>
  <Card icon={granit} name="granit" bind:selected={selected}></Card>
  <Card icon={proton} name="proton" bind:selected={selected}></Card>
  <Card icon={rune} name="rune" bind:selected={selected}></Card>
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
