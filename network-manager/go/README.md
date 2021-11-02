# Network Manager Module
First version of __NetworkManager__ module implemented using Go and github.com/gophercloud/gophercloud v0.22.

It exposes 3 API to interact with OpenStack and create/retrieve/delete networks and subnets:
```
GET network?name={name}
POST network?name={name} + body
DELETE network?name={name}
```



These API are described in the Postman collection `NSMM - NetworkManager`

## Run
Install Go: https://golang.org/doc/install
All the dependencies are listed in the `go.mod`

Run the code:
```
cd go
go run .
```
or:
```
cd go
go install .
./{GOPATH}/network-manager
```

## Terminate
Terminate the program with <CTRL+C> to revoke the token