# network-service-mesh-manager (NSMM)
First version of Network Service Mesh Manager implemented using Go and github.com/gophercloud/gophercloud v0.23.

## Architecture
![](docs/architecture.png)


## Project Structure
Following the basic layout for Go application projects (https://github.com/golang-standards/project-layout), the project structure is described below:
```
.
├── api/                # Postman collection and tests with OpenApi code generators
├── cmd/                # Main applications
    └── nsmm
        └── main.go
    └── test         # Test executable (for VPNaaS or other tests)
        └── test.go
├── docs/               # Docs/images
├── internal/           # Internal packages
    ├── config          # NSMM config package
    ├── gateway-config  # VPNaaS client to configure the gateway and VPN
    ├── nsm             # NSMM core package
    ├── openstackdriver
    ├── stubdriver
    └── vim
├── sbi/                # SBI realized as a Postman Collection (test)
├── test/
├── config.yaml         # Config file
├── go.mod
├── go.sum
└── README.md
```
- The functionalities of the Network Manager and Gateway Manager are implemented in `internal/nsm/`
- The functionalities of the Gateway Config are implemented in `internal/gateway-config`

## NBI API
It exposes API to create networks and saps, configure the gateway and create/delete secure connections.
These API are described in the Postman collection `NSMM.postman_collection.json`

## Configuration
The main program (NSMM) loads some configuration parameters from the `config.yaml`, for example, the DB credential, the VIM's info.

## Run
Install Go: https://golang.org/doc/install

All the dependencies are listed in the `go.mod`

Run the code:
```
go run cmd/nsmm/main.go
```
or creating an executable file:
```
cd /cmd/nsmm
go build
./nsmm
```

## Terminate
Terminate the program with <CTRL+C> to revoke the token

## Generate API
NBI is generated using [oapi-codegen](https://github.com/deepmap/oapi-codegen), using the following commands:
```
go run cmd/oapi-codegen/oapi-codegen.go -generate spec api/nsmm-api.json > api/spec.go
go run cmd/oapi-codegen/oapi-codegen.go -generate gin api/nsmm-api.json > api/server.go
go run cmd/oapi-codegen/oapi-codegen.go -generate types api/nsmm-api.json > api/types.go
```

# SBI
The SBI, which is OpenStack API used to create needed network resources, are definedin in a Postman Collection:
[Readme](sbi/README.md)
