# NSMM Test environment
NSMM Test environment is composed of
- a set of Go tests (To be finished)
- a MySQL DB

# Database
## How to start the test environment
```
cd db
docker build -t mysql-nsmm-db .
docker-compose up -d
```
The docker-compose file starts 2 containers:
- mysql with a nsmm DB initialized with a gateways table corresponding to a GatewayConnectivity. The port exposed on localhost is the default one.
- a phpmyadmin exposed on localhost 8080 port.

## How to stop
```
docker-compose down
```

# Go Tests (TO BE FINISHED)
## How to execute go tests
Execute all the tests files
```
go test ./...
```

Execute only one test file
```
go test test/provisioning_test.go
```

Execute a function
```
go test packageName -run NameOfTest
go test ./... -run "GatewayConnectivityRetrieve" -v
```
