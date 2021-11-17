package openstackclient

import (
	// "errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	log "github.com/sirupsen/logrus"
)

// Init function initializes the fields of the OpenStackClient struct,
// related to OpenStack endpoints, which are the ProviderClient and the ServiceClients
// ServiceClients are similar to endpoints for each OpenStack service
func (client *OpenStackClient) Init() {
	log.Info("Init function to authenticato to OpenStack")
	var err error

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: client.IdentityEndpoint,
		Username:         client.Username,
		Password:         client.Password,
		TenantID:         client.TenantID,
		DomainID:         client.DomainID,
	}

	// OpenStack providerClient
	client.provider, err = openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Authentication Token: " + client.provider.TokenID)

	// retrieve IdentityClient as a ServiceClient
	client.identityClient, err = openstack.NewIdentityV3(client.provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Identity Endpoint: " + client.identityClient.IdentityEndpoint)

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

// Close function to "close the connection with OS" which is
// revoking the token created in the initial phase
func (client *OpenStackClient) Close() {
	// Revoke token
	log.Info("Revoking Token...")
	token, err := tokens.Revoke(client.identityClient, client.provider.TokenID).Extract()
	if err != nil {
		log.Error(err)
	}
	log.Info(*token)

	// TODO: token seems still valid with a GET /networks/
}
