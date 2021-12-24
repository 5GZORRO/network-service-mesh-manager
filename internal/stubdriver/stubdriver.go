package stubdriver

import (
	log "github.com/sirupsen/logrus"
)

// VIM Driver for OpenStack, it implements VimDriver
type StubDriver struct {
	Username string
	Password string
}

func NewStubDriver(username string, password string) *StubDriver {
	return &StubDriver{
		Username: username,
		Password: password,
	}
}

// Authenticate function
func (client *StubDriver) Authenticate() {
	log.Info("Authenticating to Stub...")
}

// Revoke token
func (client *StubDriver) Revoke() {
	log.Info("Close connection to Stub...")
}

func (obj *StubDriver) CreateNetwork(networkName string, cidr string) (string, string, string, error) {
	log.Info("Creating Network network name ", networkName, " on Stub...")
	networkID := "test"
	subnetID := "test"
	subnetName := networkName + "_subnet"

	return networkID, subnetID, subnetName, nil
}

func (obj *StubDriver) DeleteNetwork(networkID string, subnetID string) error {
	log.Info("Deleting Network on Stub...")
	return nil
}

func (obj *StubDriver) CreateSAP(floatingNetName string, networkName string, cidr string) (string, string, string, string, string, string, error) {
	networkID := "test"
	subnetID := "test"
	subnetName := networkName + "_subnet"
	routerID := "routertest"
	routerName := "routername"
	portID := "porttest"
	// floatingID := "floatingID"
	log.Info("Creating SAP with network name ", networkName, " on Stub... network + router + interface to floating")
	return networkID, subnetID, subnetName, routerID, routerName, portID, nil
}

func (obj *StubDriver) DeleteSAP(networkID string, subnetID string, routerID string, portID string) error {
	log.Info("Deleting SAP on Stub...")
	return nil
}

func (client *StubDriver) RetrieveNetwork(id string) {
	log.Info("Retrieving Network on Stub...")
}
