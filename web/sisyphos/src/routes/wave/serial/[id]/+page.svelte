<script lang="ts">
  import { create } from '@bufbuild/protobuf';
  import { type Serial, UpdateRequestSchema, CreateRequestSchema, GetRequestSchema, SerialSchema, DeleteRequestSchema } from "$lib/sdk/types/wave/v1/serial/message_pb";
  import { SetException } from "$lib/exception/exception.svelte";
  import Button from "$lib/component/Button/Button.svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import { SerialClient } from "$lib/client/client.svelte";

  import Overview from "./Overview.svelte";
  import SerialPanel from "./SerialPanel.svelte";

  let serial: Serial = $state(create(SerialSchema, {node: "", reqnode: "", config: {
    name: "",
  }}))

  async function getSerial(id: string) {
    try {
      const request = create(GetRequestSchema, {
        id: id,
      });

      const response = await SerialClient().get(request)
      if (response.serial) {
        serial = response.serial
      }
    } catch (err: any) {
      SetException({title: "RETRIEVE SERIAL", desc: err.message})
    }
  }

  async function createSerial() {
    try {
      const request = create(CreateRequestSchema, {
        config: serial.config,
      })

      const response = await SerialClient().create(request)
      goto(`/wave/serial/${response.id}`);
    } catch (err: any) {
      SetException({title: "CREATE SERIAL", desc: err.message})
    }
  }

  async function updateSerial(id: string) {
    try {
      const request = create(UpdateRequestSchema, {
        id: id,
        config: serial.config,
      })

      await SerialClient().update(request)
    } catch (err: any) {
      SetException({title: "UPDATE SERIAL", desc: err.message})
    }
  }

  async function deleteSerial(id: string) {
    try {
      const request = create(DeleteRequestSchema, {
        id: id,
      })

      await SerialClient().delete(request)
      goto(`/wave/serial`);
    } catch (err: any) {
      SetException({title: "DELETE SERIAL", desc: err.message})
    }
  }

  $effect.root(() => {
    if (page.params.id !== "new") {
      getSerial(page.params.id)
    }
  })
</script>

<div class="w-11/12 flex flex-col gap-4 p-2 mt-20">
  <Overview bind:serial={serial}></Overview>

  <SerialPanel bind:serial={serial}></SerialPanel>

  <div class="w-full flex flex-row justify-start gap-3 my-2">
    {#if page.params.id === "new"}
      <Button onclick={() => serial = create(SerialSchema, {config: {}})} scale={0.8} class="ml-auto flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/30">
        <span>Reset</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="40" stroke-dashoffset="40" d="M17 15.33c2.41 -0.72 4 -1.94 4 -3.33c0 -2.21 -4.03 -4 -9 -4c-4.97 0 -9 1.79 -9 4c0 2.06 3.5 3.75 8 3.98"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.5s" values="40;0"/></path><path fill="currentColor" d="M12.25 16l0 0l0 0z" opacity="0"><animate fill="freeze" attributeName="d" begin="0.5s" dur="0.2s" values="M12.25 16l0 0l0 0z;M12.25 16L9.5 13.25L9.5 18.75z"/><set fill="freeze" attributeName="opacity" begin="0.5s" to="1"/></path></g></svg>      
      </Button>
      <Button onclick={() => createSerial()} scale={0.8} class="flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/30">
        <span>Create</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-dasharray="16" stroke-dashoffset="16" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path d="M5 12h14"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.4s" values="16;0"/></path><path d="M12 5v14"><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.4s" dur="0.4s" values="16;0"/></path></g></svg>
      </Button>
    {:else}
      <Button onclick={() => deleteSerial(page.params.id)} scale={0.8} class="ml-auto flex flex-row gap-2 justify-center w-32 p-2 rounded-lg text-slate-50/80 bg-red-900/80">
        <span>Delete</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="64" stroke-dashoffset="64" d="M12 3c4.97 0 9 4.03 9 9c0 4.97 -4.03 9 -9 9c-4.97 0 -9 -4.03 -9 -9c0 -4.97 4.03 -9 9 -9Z"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.6s" values="64;0"/></path><path stroke-dasharray="8" stroke-dashoffset="8" d="M12 12l4 4M12 12l-4 -4M12 12l-4 4M12 12l4 -4"><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.6s" dur="0.2s" values="8;0"/></path></g></svg>
      </Button>
      <Button onclick={() => getSerial(page.params.id)} scale={0.8} class="flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/30">
        <span>Reset</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="40" stroke-dashoffset="40" d="M17 15.33c2.41 -0.72 4 -1.94 4 -3.33c0 -2.21 -4.03 -4 -9 -4c-4.97 0 -9 1.79 -9 4c0 2.06 3.5 3.75 8 3.98"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.5s" values="40;0"/></path><path fill="currentColor" d="M12.25 16l0 0l0 0z" opacity="0"><animate fill="freeze" attributeName="d" begin="0.5s" dur="0.2s" values="M12.25 16l0 0l0 0z;M12.25 16L9.5 13.25L9.5 18.75z"/><set fill="freeze" attributeName="opacity" begin="0.5s" to="1"/></path></g></svg>
      </Button>
      <Button onclick={() => updateSerial(page.params.id)} scale={0.8} class="flex flex-row gap-2 justify-center w-32 p-2 rounded-lg bg-slate-50/30">
        <span>Update</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="64" stroke-dashoffset="64" d="M3 12c0 -4.97 4.03 -9 9 -9c4.97 0 9 4.03 9 9c0 4.97 -4.03 9 -9 9c-4.97 0 -9 -4.03 -9 -9Z"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.6s" values="64;0"/></path><path stroke-dasharray="14" stroke-dashoffset="14" d="M8 12l3 3l5 -5"><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.6s" dur="0.2s" values="14;0"/></path></g></svg>
      </Button>
    {/if}
  </div>
</div>

