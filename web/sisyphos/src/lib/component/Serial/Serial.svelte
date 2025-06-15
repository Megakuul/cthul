<script lang="ts">
  import { create } from '@bufbuild/protobuf';
  import { DomainClient, SerialClient } from "$lib/client/client.svelte";
  import { SetException } from "$lib/exception/exception.svelte";
  import { ConnectRequestSchema, type ConnectRequest } from "$lib/sdk/types/wave/v1/serial/message_pb";
  import Button from '../Button/Button.svelte';
  import type { SerialDevice } from '$lib/sdk/types/wave/v1/domain/config_pb';

  import { Xterm, XtermAddon } from '@battlefieldduck/xterm-svelte';
	import type {
		ITerminalOptions,
		ITerminalInitOnlyOptions,
		Terminal,
    FitAddon
	} from '@battlefieldduck/xterm-svelte';

  let {
    devices,
  }: {
    devices: SerialDevice[];
  } = $props()

  let terminal: Terminal;

  let container: HTMLDivElement;

  let options: ITerminalOptions & ITerminalInitOnlyOptions = {
    fontFamily: "Consolas",
  };

  let deviceIndex: number = $state(0);

  class StreamController<T> implements AsyncIterable<T> {
    private queue: T[] = [];
    private isClosed = false;
    private waitingPromise: Promise<void> | null = null;
    private notify: (() => void) | null = null;

    public write(data: T): void {
      if (this.isClosed) return;
      this.queue.push(data);
      if (this.notify) {
        this.notify();
        this.notify = null;
        this.waitingPromise = null;
      }
    }

    public close(): void {
      this.isClosed = true;
      if (this.notify) {
        this.notify();
        this.notify = null;
        this.waitingPromise = null;
      }
    }

    public async *[Symbol.asyncIterator](): AsyncGenerator<T, void, undefined> {
      while (true) {
        while (this.queue.length > 0) {
          yield this.queue.shift()!;
        }
        if (this.isClosed) {
          return;
        }
        if (!this.waitingPromise) {
            this.waitingPromise = new Promise((resolve) => {
              this.notify = resolve;
            });
        }
        await this.waitingPromise;
      }
    }
  }

  let controller: StreamController<ConnectRequest>;

  let fitAddon: FitAddon;

  async function connect() {
    try {
      const responseStream = SerialClient().connect(controller)
      for await (const response of responseStream) {
        terminal.write(response.output)
      }
    } catch (err: any) {
      SetException({title: "CONNECT SERIAL", desc: err.message})
    }
  }

  async function onLoad() {
    fitAddon = new(await XtermAddon.FitAddon()).FitAddon();
    terminal.loadAddon(fitAddon)
    fitAddon.fit()

    const resizeObserver = new ResizeObserver(() => {
      fitAddon.fit();
    });
    resizeObserver.observe(container);

    controller = new StreamController<ConnectRequest>();
    await connect()
  }

  const encoder = new TextEncoder();

  function onData(data: string) {
    controller?.write(create(ConnectRequestSchema, { input: encoder.encode(data) }))
  }
</script>

<div bind:this={container} class="flex flex-col w-full h-full">
  <Xterm bind:terminal {options} {onLoad} {onData} class="h-full" />
  <div class="flex flex-row items-center justify-start gap-2 py-1 px-2">
    <Button scale={0.4} class="p-1 rounded-lg bg-slate-50/10" onclick={async () => {await connect()}}>
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path stroke-dasharray="40" stroke-dashoffset="40" d="M17 15.33c2.41 -0.72 4 -1.94 4 -3.33c0 -2.21 -4.03 -4 -9 -4c-4.97 0 -9 1.79 -9 4c0 2.06 3.5 3.75 8 3.98"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.5s" values="40;0"/></path><path fill="currentColor" d="M12.25 16l0 0l0 0z" opacity="0"><animate fill="freeze" attributeName="d" begin="0.5s" dur="0.2s" values="M12.25 16l0 0l0 0z;M12.25 16L9.5 13.25L9.5 18.75z"/><set fill="freeze" attributeName="opacity" begin="0.5s" to="1"/></path></g></svg>
    </Button>
    <Button scale={0.4} class="p-1 rounded-lg bg-slate-50/10" onclick={() => {
      if (deviceIndex + 1 >= devices.length) deviceIndex = 0
      else deviceIndex = deviceIndex + 1
    }}>
      <p class="w-[24px] h-[24px]">{deviceIndex}</p>
    </Button>
    <p class="ml-auto text-sm">status: <span class="font-bold">connected</span></p>
  </div>
</div>
