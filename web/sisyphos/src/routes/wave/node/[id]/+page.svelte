<script lang="ts">
  import { create } from '@bufbuild/protobuf';
  import { type Node, GetRequestSchema, NodeSchema } from "$lib/sdk/types/wave/v1/node/message_pb";
  import { SetException } from "$lib/exception/exception.svelte";
  import { page } from "$app/state";
  import { NodeClient } from "$lib/client/client.svelte";

  import Overview from "./Overview.svelte";

  let node: Node = $state(create(NodeSchema, {}))

  async function getNode(id: string) {
    try {
      const request = create(GetRequestSchema, {
        id: id,
      });

      const response = await NodeClient().get(request)
      if (response.node) {
        node = response.node
      }
    } catch (err: any) {
      SetException({title: "RETRIEVE NODE", desc: err.message})
    }
  }

  $effect.root(() => {
    if (page.params.id !== "new") {
      getNode(page.params.id)
    }
  })
</script>

<div class="w-11/12 flex flex-col gap-4 p-2 mt-20">
  <Overview bind:node={node}></Overview>
</div>

