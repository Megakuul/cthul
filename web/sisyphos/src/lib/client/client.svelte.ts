import { createClient, type Transport } from "@connectrpc/connect";
import {DomainService} from "$lib/sdk/types/wave/v1/domain/service_pb"
import { DiskService } from "$lib/sdk/types/granit/v1/disk/service_pb";
import { VideoService } from "$lib/sdk/types/wave/v1/video/service_pb";
import { SerialService } from "$lib/sdk/types/wave/v1/serial/service_pb";
import { NewTransport, Port } from "$lib/sdk/transport";

let waveTransport: Transport = $state(NewTransport("https://cthul.io", Port.WAVE))
let granitTransport: Transport = $state(NewTransport("https://cthul.io", Port.GRANIT))

let domainClient = $derived(createClient(DomainService, waveTransport));
let videoClient = $derived(createClient(VideoService, waveTransport));
let serialClient = $derived(createClient(SerialService, waveTransport));
let diskClient = $derived(createClient(DiskService, granitTransport));

export function SetTransport(url: string) {
  waveTransport = NewTransport(url, Port.WAVE)
  granitTransport = NewTransport(url, Port.GRANIT)
}

export const DomainClient = () => domainClient;
export const VideoClient = () => videoClient;
export const SerialClient = () => serialClient;
export const DiskClient = () => diskClient;