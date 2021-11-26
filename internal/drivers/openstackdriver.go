package drivers

import (
	log "github.com/sirupsen/logrus"
)

type OpenStackDriver struct {
	// Global params
	IdentityEndpoint string
	Username         string
	Password         string
	TenantID         string
	DomainID         string
	// openstack
	// provider       *gophercloud.ProviderClient
	// identityClient *gophercloud.ServiceClient
	// networkClient  *gophercloud.ServiceClient
	// computeClient  *gophercloud.ServiceClient
}

func NewOpenStackDriver(identityEndpoint string, username string, password string, tenantID string, domainID string) *OpenStackDriver {
	return &OpenStackDriver{
		IdentityEndpoint: identityEndpoint,
		Username:         username,
		Password:         password,
		TenantID:         tenantID,
		DomainID:         domainID,
	}
}

// Authenticate function towards OpenStack, it has to be executed before all the other methods
func (obj *OpenStackDriver) Authenticate() {
	log.Info("Authenticating...")
}

// CreateGatewayConnectivity
func (obj *OpenStackDriver) CreateGatewayConnectivity(sliceId string, subnet string) (string, string, string, error) {
	// Create private_network, router with gateway, interface and returns the VIM-id
	log.Info("CreateGatewayConnectivity...")
	return "", "", "", nil
}

// DeleteGatewayConnectivity
func (obj *OpenStackDriver) DeleteGatewayConnectivity(networkId string, subnetId string, routerId string) error {
	// Delete all the gateway resources
	log.Info("DeleteGatewayConnectivity...")
	return nil
}

// GetGatewayConnectivity
func (obj *OpenStackDriver) GetGatewayConnectivity(networkId string, subnetId string, routerId string) {
	log.Info("GetGatewayConnectivity...")
}

// Revoke
func (obj *OpenStackDriver) Revoke() {
	log.Info("Revoking token...")
}
