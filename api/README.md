# Network Service Mesh Manager NBI
It exposes API to:
- create network resources (networks and sap) 
- configure the gateway 
- manage VPN connections, configure the gateway and create/delete secure connections.

The API is defined in `nsmm.json`. Examples of the NBI are in the Postman collection `NSMM.postman_collection.json`

## NSMM API
The NorthBound Interface of NSMM is detailed [here](https://5gzorro.github.io/network-service-mesh-manager/ "NSMM API")


## Generate go server NBI from JSON API
NBI of the GIN server is generated using [oapi-codegen](https://github.com/deepmap/oapi-codegen), using the following commands:
```
go run cmd/oapi-codegen/oapi-codegen.go -generate spec api/nsmm.json > api/spec.gen.go
go run cmd/oapi-codegen/oapi-codegen.go -generate gin api/nsmm.json > api/gin-server.gen.go
go run cmd/oapi-codegen/oapi-codegen.go -generate types api/nsmm.json > api/types.gen.go
```