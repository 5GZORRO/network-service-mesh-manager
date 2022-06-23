# Network Service Mesh Manager (NSMM)
## Introduction
Network Service Mesh Manager is a service to establish secure connection between slices/network services. It can:
- create network resource on the VIM needed by the network service
- configure the gateway, which is a VNF in the _Network Service_ that contains the [VPN as a Service](https://github.com/5GZORRO/VPNaaS) module.

### High-Level Architecture
![](docs/architecture.png)
- The `Network Manager` module manages the creation/removal of requested networks and subnets on the VIM
- The `Gateway Manager` module manages the creation/removal of the resources need to create a SAP (allocate a floatingIP), interacting with the Network Manager for networks and subnets
- The `Gateway Config` module handles the configuration of the gateways (VPNaaS) of each network service, starting the VPN, creating/removing connections

### NorthBound Interface
The NorthBound Interface of NSMM is described [here](api/README.md)


## Prerequisities
### System requirements
Minimum requirements:
- 2vCPU
- 4GB
- 10GB
- Deployment/Virtualization technology: VMs/Containers

### Software dependencies
Network Service Mesh Manager is implemented using Go (1.17) and uses:
- GIN (v1.7.4) as a HTTP server 
- Postgres Driver (v1.2.3) as DB driver
- github.com/gophercloud/gophercloud (v0.23) as SBI for openstack VIM
all the dependecies are listed in the `go.mod` file

### 5GZORRO Module dependencies
NSMM uses as dependencies:
- [VPN as a Service](https://github.com/5GZORRO/VPNaaS) module, which should be deployed in the GW VNF of the Network Services.
- [Identity and Permission Manager](https://github.com/5GZORRO/identity) module, which should be deployed in the platform and reachable from the NSMM and the GWs of the Network Services.

## Installation
NSMM can be deployed in 2 different ways:
1. running the executable
2. using Docker with docker-compose
3. using Helm (and K8s)

As a first step, download the repo:
```
git clone https://github.com/5GZORRO/network-service-mesh-manager.git
```

### 1. Running executable
If you want to locally run the executable, first you need to:
- Install Go: https://golang.org/doc/install
- Download all the dependecies using the `go.mod`
- Install and start Postgres (use `/test/db/docker-compose.yaml` to start a Postgres DB locally in a container plus phppgadmin)

To run NSMM without build it separately, you can type the following command that loads the config file in the root directory:
```
go run cmd/nsmm/main.go
```
Otherwise if you need to specify a different config file:
```
go run cmd/nsmm/main.go -config=<filename>
```

To build and run the executable file and run it:
```
cd /cmd/nsmm
go build
./nsmm
```

Terminate the program with <CTRL+C> to revoke the token

### 2. docker-compose deployment
The `docker-compose.yaml` in the `deployment/docker` directory, contains a complete environment. Check the variables defined in the `.env` before starting containers.

1. Check values in `deployment/docker/config.yaml`, which is mapped as volume in the nsmm container, so used at startup
2. Start docker-compose 
```
docker-compose up -d
```
3. Terminate with
```
docker-compose down
```

### 3. Helm deployment
The Helm chart is defined in the `deployment/helm` directory. Before starting the chart, check the values in the `values.yaml` file.

1. To start the chat
```
cd deployment/helm
helm install nsmm nsmm-chart/
```
2. Terminate with
```
helm uninstall nsmm
```

## Configuration
The main program (NSMM) loads some configuration parameters from the `config.yaml`. For example, the DB credential and the VIM's info.

### Standalone configuration
If you want to deploy NSMM as standalone module (for testing purposes, without interacting with the other 5GZORRO modules) you need to edit in the `config.yaml` the `vpnaas` section and set the `environment` as __local__. In this way the module will use local keys without contacting the ID&P

```
vpnaas:
  vpnaasPort: 8080
  environment: local
```

## Maintainers
Elena Bucchianeri - Developer and Designer - e.bucchianeri@nextworks.it

Pietro G. Giardina - Developer and Designer - p.giardina@nextworks.it

## License
This module is distributed under [Apache 2.0 LICENSE](LICENSE) terms.
