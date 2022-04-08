package stubdriver

import (
	"nextworks/nsm/internal/config"

	log "github.com/sirupsen/logrus"
)

// VIM Driver for Testing, it implements VimDriver
type StubDriver struct {
	Username            string
	Password            string
	StaticGateway       *config.StaticGateway // Fixed static Gateway for this VIM
	FloatingNetworkName string
	FloatingNetworkID   string
}

func NewStubDriver(username string, password string, floatingID string, floatingnet string, staticgw *config.StaticGateway) *StubDriver {
	return &StubDriver{
		Username:            username,
		Password:            password,
		FloatingNetworkID:   floatingID,
		FloatingNetworkName: floatingnet,
		StaticGateway:       staticgw,
	}
}

func (client *StubDriver) GetStaticGatewayInfo() (string, string, string, string, string, error) {
	log.Info("Retrieving Static Gateway Info from Stub...")
	if client.StaticGateway == nil {
		log.Error("No Static Gateway specified for VIM Stub...")
		return "", "", "", "", "", ErrStaticGatewayNotFound
	} else {
		return client.StaticGateway.NetworkName, client.StaticGateway.SubnetCidr, client.StaticGateway.ExternalInterfaceName, client.StaticGateway.FloatingIP, client.StaticGateway.InstanceID, nil
	}

}

// Authenticate function
func (client *StubDriver) Authenticate() {
	log.Info("Authenticating to Stub...")
}

func (client *StubDriver) AllocateFloatingIP(networkname string) (string, string, string, string, error) {
	log.Info("Allocating FIP using Stub...")
	return "portID", "ens5", "fipID", "10.30.6.6", nil
}

func (client *StubDriver) DeallocateFloatingIP(portID string, fipID string) error {
	log.Info("Deallocating FIP using Stub...")
	return nil
}

// RetrieveFloatingNetworkID function
func (client *StubDriver) RetrieveFloatingNetworkID() string {
	log.Info("Retrieve FloatingNetworkID for Stub...")
	return client.FloatingNetworkID
}

// RetrieveFloatingNetworkName function
func (client *StubDriver) RetrieveFloatingNetworkName() string {
	log.Info("Retrieve FloatingNetworkName for Stub...")
	return client.FloatingNetworkName
}

// Revoke token
func (client *StubDriver) Revoke() {
	log.Info("Close connection to Stub...")
}

func (obj *StubDriver) CreateNetwork(networkName string, cidr string, gateway bool) (string, string, string, error) {
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
	log.Info("Creating SAP with network name ", networkName, " on Stub... network + router + interface to floating ", floatingNetName)
	return networkID, subnetID, subnetName, routerID, routerName, portID, nil
}

func (obj *StubDriver) DeleteSAP(networkID string, subnetID string, routerID string, portID string) error {
	log.Info("Deleting SAP on Stub...")
	return nil
}

func (client *StubDriver) RetrieveNetwork(id string) {
	log.Info("Retrieving Network on Stub...")
}

// Functions used in case of StaticGws
func (client *StubDriver) CreateInterfacePort(serverId string, networkId string) (string, error) {
	log.Info("Creating interface Port from ServerId on Stub...")
	return "portID", nil
}

func (client *StubDriver) DeleteInterfacePort(serverId string, portId string) error {
	log.Info("Deleting interface Port from ServerId on Stub...")
	return nil
}
