<!-- TODO: style + load values from oauth auth source -->
<script>
    import { SetTransport } from "$lib/client/client.svelte";

  import Button from "$lib/component/Button/Button.svelte";
  import Input from "$lib/component/Input/Input.svelte";

  let apiEndpoint = "http://127.0.0.1"
</script>

<div class="w-10/12 my-10 bg-slate-50/20 rounded-xl flex flex-col gap-3 p-4">
  <h1 class="text-2xl font-bold">
    Preferences
  </h1>
  <Input title="api endpoint" type="text" bind:value={apiEndpoint} class="" />

  <Button class="p-2 rounded-lg" scale={0.8} onclick={() => {
    SetTransport(apiEndpoint)
  }}>
    Submit
  </Button>
</div>

