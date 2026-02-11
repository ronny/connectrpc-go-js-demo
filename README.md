# ConnectRPC Demo: Go backend + JS frontend

## Getting started

```
go run ./cmd/backend-api

# or

make
./bin/backend-api
```

You can use grpcurl or grpcui to test the API via gRPC, e.g.:

```
grpcurl -plaintext localhost:8080 list myorg.demo.v1.DemoAPI

grpcurl -plaintext localhost:8080 myorg.demo.v1.DemoAPI.Login

grpcurl -plaintext -d '{"email": "void@example.com", "password": "dogecoin"}' localhost:8080 myorg.demo.v1.DemoAPI.Login

grpcurl -plaintext -H "Demo-Auth-Token: xyz" localhost:8080 myorg.demo.v1.DemoAPI.GetBalance

grpcurl -plaintext -H "Demo-Auth-Token: xyz" localhost:8080 myorg.demo.v1.DemoAPI.CreateTransfer
```

or use curl:

```
curl -X POST -H "Content-Type: application/json" -d '{"email": "void@example.com", "password": "dogecoin"}' http://localhost:8080/myorg.demo.v1.DemoAPI/Login

curl -X GET -H "Demo-Auth-Token: xyz" http://localhost:8080/myorg.demo.v1.DemoAPI/GetBalance
```

For development you'll need these:

```
go install github.com/bufbuild/buf/cmd/buf@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
```

and golangci-lint.

Frontend: TODO, but basically use https://www.npmjs.com/package/@bufbuild/protoc-gen-es
to generate the client SDK from the proto files.
