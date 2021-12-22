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

func (obj *StubDriver) CreateNetwork(name string, cidr string) (string, string, string, error) {
	log.Info("Creating Network on Stub...")
	return "", "", "", nil
}

func (client *StubDriver) RetrieveNetwork(id string) {
	log.Info("Retrieving Network on Stub...")
}

func (obj *StubDriver) DeleteNetwork(id string) {
	log.Info("Deleting Network on Stub...")
}

func (obj *StubDriver) CreateSAP() {
	log.Info("Creating SAP on Stub...")
}

func (obj *StubDriver) DeleteSAP() {
	log.Info("Deleting SAP on Stub...")
}
