import { createClient, type Transport } from "@connectrpc/connect";
import {DomainService} from "$lib/sdk/types/wave/v1/domain/service_pb"
import { DiskService } from "$lib/sdk/types/granit/v1/disk/service_pb";
import { VideoService } from "$lib/sdk/types/wave/v1/video/service_pb";
import { SerialService } from "$lib/sdk/types/wave/v1/serial/service_pb";
import { NewTransport, Port } from "$lib/sdk/transport";
import { NodeService } from "$lib/sdk/types/wave/v1/node/service_pb";

let waveTransport: Transport = $state(NewTransport("https://cthul.io", Port.WAVE))
let granitTransport: Transport = $state(NewTransport("https://cthul.io", Port.GRANIT))

let domainClient = $derived(createClient(DomainService, waveTransport));
let nodeClient = $derived(createClient(NodeService, waveTransport));
let serialClient = $derived(createClient(SerialService, waveTransport));
let videoClient = $derived(createClient(VideoService, waveTransport));
let diskClient = $derived(createClient(DiskService, granitTransport));

export function SetTransport(url: string) {
  waveTransport = NewTransport(url, Port.WAVE)
  granitTransport = NewTransport(url, Port.GRANIT)
}

export const DomainClient = () => domainClient;
export const NodeClient = () => nodeClient;
export const SerialClient = () => serialClient;
export const VideoClient = () => videoClient;
export const DiskClient = () => diskClient;