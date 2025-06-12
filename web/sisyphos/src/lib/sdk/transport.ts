import { Code, ConnectError, type Interceptor, type Transport } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

export enum Port {
  WAVE = "1870",
  GRANIT = "1970",
  PROTON = "2070",
  RUNE = "2170"
}

export function NewTransport(url: string, port: Port): Transport {
  return createConnectTransport({
    baseUrl: `${url}:${port}`,
    interceptors: [authInterceptor, redirectInterceptor],
  })
}

const redirectInterceptor: Interceptor = (next) => async (req) => {
  try {
    return await next(req);
  } catch (err) {
    if (err instanceof ConnectError && err.code === Code.NotFound) {
      const redirect = err.metadata.get("location");
      if (redirect) return await next({...req, url: redirect});
    }
    throw err
  }
}

const authInterceptor: Interceptor = (next) => async (req) => {
  // req.header.set("authorization", "Bearer ey")
  return await next(req)
}