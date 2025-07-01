# Cthul
---

![overview](/assets/overview.jpg)


## Development setup
---

### Linux / WSL

#### Services

Prerequisites:
- `go`
- `etcd` (config located at `/configs/etcd/config.yaml`)
- `qemu-system`
- `virt-manager`

Launch all services you need with:
```bash
go run cmd/<service>/<service>.go --config ../../configs/<service>/config.toml`
```

#### Web

Prerequisites:
- `npm`
- `nodejs`

Launch the sisyphos webservice with:
```bash
cd web/sisyphos
npm i
npm run dev
```


When developing the sisyphos web interface, you don't need to run live backend services (especially in early phases where services may be subject to heavy refactorings).

Instead, you can rely on the stable central protobuf API. For testing purposes, use the provided types at `/web/sisyphos/lib/sdk/types` to generate mock data. e.g.:

```javascript
import { create } from '@bufbuild/protobuf';
import { type Domain, DomainSchema } from "$lib/sdk/types/wave/v1/domain/message_pb";

let domain: Domain = $state(create(DomainSchema, {node: "", reqnode: "", config: {
  name: "test",
  description: "blub",
  affinity: [],
  resourceConfig: {vcpus: BigInt(2), memory: BigInt(4 * (1000 * 1000 * 1000))},
  firmwareConfig: {firmware: Firmware.SEABIOS, 
    loaderDeviceId: "", tmplDeviceId: "", nvramDeviceId: "", secureBoot: false,
  },
}}))
```

## Cthul api
---

Cthul components implement a gRPC (ConnectRPC) api for internal and external communication.

`proto` definitions are located in `api/<component>/<version>/`, the compiled go outputs are located in `pkg/api/<component>/<version>/`. 

Go outputs are generated using the `buf` compiler (`buf.yaml` & `buf.gen.yaml`) in combination with the `protoc-gen-connect-go` plugin.


The api is documented with inline comments in the proto definitions.


### Cthul etcd configuration

Cthul uses an etcd database on every node. To enhance security and lower complexity the etcd client is only exposed on a local unix socket. The etcd peer interface (used for cluster replication) is exposed on port `2380` and uses mTLS configured to only accept public certs with a SAN that matches `*.etcd.cthul.io`.


The etcd configuration and init files can be found in `/configs/etcd` and `/init/etcd`.

Setting up a etcd cluster requires configuration of users and permissions:
bash
```
export ETCDCTL_ENDPOINTS="unix:///var/run/cthul/etcd/etcd.sock"
etcdctl user add root && export ETCDCTL_USER="root"
etcdctl user grant-role root root
etcdctl auth enable

# TODO: Document component specific user and role creation.
```


### Concepts

Cthul components all use a `controller` and `operator` architecture. Components have internal operators, usually located at `internal/<comp>/<operator>`, which handle events (e.g. creation of a device).
Every operator also exposes a controller, usually located at `pkg/<comp>/<controller>`, which is used by other components to emit events to the operators (e.g. request device).


This concept is consistent for cthul components but is also used for external components, for example domains provide a DomainController in `pkg/domain` which emits events to the operator which in this case is `libvirtd`.



External components are abstracted with a go interface that defines simple and replaceable functions.
The goal is to only interface what exactly is required, this makes it simpler to replace the external component
behind it.
