package openstackdriver

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	log "github.com/sirupsen/logrus"
)

// VIM Driver for OpenStack, it implements VimDriver
type OpenStackDriver struct {
	// Global params
	IdentityEndpoint string
	Username         string
	Password         string
	TenantID         string
	DomainID         string
	// openstack
	provider       *gophercloud.ProviderClient
	identityClient *gophercloud.ServiceClient
	networkClient  *gophercloud.ServiceClient
	computeClient  *gophercloud.ServiceClient
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
func (client *OpenStackDriver) Authenticate() {
	log.Info("Authenticating to OpenStack...")
	var err error

	// Without token and scope
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: client.IdentityEndpoint,
		Username:         client.Username,
		Password:         client.Password,
		TenantID:         client.TenantID,
		DomainID:         client.DomainID,
		AllowReauth:      true,
	}

	// OpenStack providerClient
	client.provider, err = openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Authentication Token provided: ", client.provider.Token())

	// retrieve IdentityClient as a ServiceClient
	client.identityClient, err = openstack.NewIdentityV3(client.provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Identity Endpoint: " + client.identityClient.IdentityEndpoint)
	// token, err := tokens.Get(client.identityClient, client.provider.Token()).Extract()
	// if err != nil {
	// 	log.Error("Retriving Authentication token error: ", err)
	// }

	// retrieve NetworkClient as a ServiceClient
	client.networkClient, err = openstack.NewNetworkV2(client.provider, gophercloud.EndpointOpts{
		Name:   "neutron",
		Region: "RegionOne",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Network Endpoint " + client.networkClient.Endpoint)

	// retrieve ComputeClient as a ServiceClient
	client.computeClient, err = openstack.NewComputeV2(client.provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Compute Endpoint " + client.computeClient.Endpoint)
}

// Revoke
func (client *OpenStackDriver) Revoke() {
	log.Info("Revoking token...")
	// token, err := tokens.Revoke(client.identityClient, client.provider.Token()).Extract()
	// if err != nil {
	// 	log.Error(err)
	// }
	// log.Info(*token)
}
