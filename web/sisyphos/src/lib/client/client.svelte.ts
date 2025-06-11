import { Code, ConnectError, createClient, type Transport } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import {DomainService} from "$lib/types/wave/v1/domain/service_pb"
import { DiskService } from "$lib/types/granit/v1/disk/service_pb";
import { VideoService } from "$lib/types/wave/v1/video/service_pb";
import { SerialService } from "$lib/types/wave/v1/serial/service_pb";

let transport: Transport = $state(createConnectTransport({baseUrl: "https://cthul.io"}))

export function SetTransport(url: string) {
  transport = createConnectTransport({
    baseUrl: url,
    interceptors: [(next) => async (req) => {
      try {
        return await next(req);
      } catch (err) {
        if (err instanceof ConnectError && err.code === Code.NotFound) {
          const redirect = err.metadata.get("location");
          if (redirect) return await next({...req, url: redirect});
        }
        throw err
      }
    }]
  })
}

export let DomainClient = $derived(createClient(DomainService, transport));
export let VideoClient = $derived(createClient(VideoService, transport));
export let SerialClient = $derived(createClient(SerialService, transport));
export let DiskClient = $derived(createClient(DiskService, transport));