package openstackclient

import (
	// "errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	log "github.com/sirupsen/logrus"
)

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

	client.provider, err = openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Authentication Token: " + client.provider.TokenID)
	client.identityClient, err = openstack.NewIdentityV3(client.provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Identity Endpoint: " + client.identityClient.IdentityEndpoint)

	client.networkClient, err = openstack.NewNetworkV2(client.provider, gophercloud.EndpointOpts{
		Name:   "neutron",
		Region: "RegionOne",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Network Endpoint " + client.networkClient.Endpoint)
}

// Function to "close the connection with OS" which is
// revoking the token created in the initial phase
func (client *OpenStackClient) Close() {
	// Revoke token
	log.Info("Revoking Token...")
	token, err := tokens.Revoke(client.identityClient, client.provider.TokenID).Extract()
	if err != nil {
		log.Error(err)
	}
	log.Info(*token)

	// TBT: token seems still valid with a GET /networks/
}
