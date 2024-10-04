

### Cthul API
---

Cthul components implement a gRPC (ConnectRPC) api for internal and external communication.

`proto` definitions are located in `api/<component>/<version>/`, the compiled go outputs are located in `pkg/api/<component>/<version>/`. 

Go outputs are generated using the `buf` compiler (`buf.yaml` & `buf.gen.yaml`) in combination with the `protoc-gen-connect-go` plugin.


The api is documented with inline comments in the proto definitions.
