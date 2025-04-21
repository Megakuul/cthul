# Cthul

Work in progress - The stuff below is irrelevant, just notes so my sievebrain doesn't forget to document it.


Components:
- wave -> domain manager
- granit -> storage manager
- proton -> network manager
- rune -> trust + authentication manager (iam)
- sisyphos -> webinterface
- flow -> log + metric collector/distributor

### Cthul api
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
