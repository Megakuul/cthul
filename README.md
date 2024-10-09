

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
