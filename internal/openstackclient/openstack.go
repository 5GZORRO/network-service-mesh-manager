package openstackclient

import (
	"github.com/gophercloud/gophercloud"
)

type OpenStackClient struct {
	// Global params
	IdentityEndpoint string
	Username         string
	Password         string
	TenantID         string
	DomainID         string
	// Clients
	provider       *gophercloud.ProviderClient
	identityClient *gophercloud.ServiceClient
	networkClient  *gophercloud.ServiceClient
	computeClient  *gophercloud.ServiceClient
}

// builder pattern
func NewOpenStackClient(identityEndpoint string, username string, password string, tenantID string, domainID string) *OpenStackClient {
	return &OpenStackClient{
		IdentityEndpoint: identityEndpoint,
		Username:         username,
		Password:         password,
		TenantID:         tenantID,
		DomainID:         domainID,
	}
}
