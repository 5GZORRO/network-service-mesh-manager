# network-service-mesh-manager (NSMM)
Network Service Mesh Manager is a service to establish secure connection between slices/network services. It can create network resource needed by the network service, configure the gateway, which is a VNF in the network service that contains the VPNaaS module.

Network Service Mesh Manager is implemented using Go and uses:
- GIN as a HTTP server 
- github.com/gophercloud/gophercloud v0.23 as SBI for openstack VIM

etc.. all the dependecies are listed in the `go.mod`

# High-Level Architecture
![](docs/architecture.png)
- The `Network Manager` module manages the creation/removal of requested networks and subnets on the VIM
- The `Gateway Manager` module manages the creation/removal of the resources need to create a SAP (allocate a floatingIP), interacting with the Network Manager for networks and subnets
- The `Gateway Config` module handles the configuration of the gateways (VPNaaS) of each network service, starting the VPN, creating/removing connections

# Project Structure
Following the basic layout for Go application projects (https://github.com/golang-standards/project-layout), the project structure is described below:
```
.
├── api/                # NBI generated code + API spec + Postman collection
├── cmd/                # Executables
    └── nsmm            # Main application (NSMM)
        └── main.go
    └── test            # Test executables (for VPNaaS or other tests)
        └── test.go
├── docs/               # Docs/images
├── internal/           # Internal packages
    ├── config          # NSMM config package
    ├── gateway-config  # VPNaaS client to configure the gateway and VPN
    ├── nsm             # NSMM core package
    ├── openstackdriver # Driver for OpenstackVIM
    ├── stubdriver      # Stub driver for testing purposes
    └── vim             # General interface
├── sbi/                # SBI realized as a Postman Collection (test)
├── test/
├── config.yaml         # Config file
├── go.mod
├── go.sum
└── README.md
```
- The functionalities of the Network Manager and Gateway Manager are implemented in `internal/nsm/`
- The functionalities of the Gateway Config are implemented in `internal/gateway-config`

# NBI API
It exposes API to:
- create network resources (networks and sap) 
- configure the gateway 
- manage VPN connections, configure the gateway and create/delete secure connections.

The API is defined in `api/nsmm.json`. Examples of the NBI are in the Postman collection `api/NSMM.postman_collection.json`


## Generate go server NBI from JSON API
NBI of the GIN server is generated using [oapi-codegen](https://github.com/deepmap/oapi-codegen), using the following commands:
```
go run cmd/oapi-codegen/oapi-codegen.go -generate spec api/nsmm-api.json > api/spec.go
go run cmd/oapi-codegen/oapi-codegen.go -generate gin api/nsmm-api.json > api/server.go
go run cmd/oapi-codegen/oapi-codegen.go -generate types api/nsmm-api.json > api/types.go
```

# Configuration
The main program (NSMM) loads some configuration parameters from the `config.yaml`, for example, the DB credential and the VIM's info.

# Run

### Dependencies
- Install Go: https://golang.org/doc/install
- All the dependencies are listed in the `go.mod`
- Start the DB

### Execute NSMM
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

### Terminate
Terminate the program with <CTRL+C> to revoke the token

# SBI
The SBI, which is OpenStack API used to create needed network resources, are definedin in a Postman Collection:
[Readme](sbi/README.md)
