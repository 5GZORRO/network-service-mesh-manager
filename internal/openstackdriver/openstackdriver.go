package openstackdriver

import (
	"nextworks/nsm/internal/config"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	log "github.com/sirupsen/logrus"
)

// VIM Driver for OpenStack, it implements VimDriver
type OpenStackDriver struct {
	// Global params
	IdentityEndpoint    string
	Username            string
	Password            string
	TenantID            string
	DomainID            string
	FloatingNetworkName string
	FloatingNetworkID   string
	// Fixed static Gateway for this VIM
	StaticGateway *config.StaticGateway
	// other config params
	AvailabilityZone string
	// openstack
	provider       *gophercloud.ProviderClient
	identityClient *gophercloud.ServiceClient
	networkClient  *gophercloud.ServiceClient
	computeClient  *gophercloud.ServiceClient
}

func NewOpenStackDriver(identityEndpoint string, username string, password string, tenantID string, domainID string,
	floatinNetID string, floatNetName string, availZone string, staticgw *config.StaticGateway) *OpenStackDriver {
	return &OpenStackDriver{
		IdentityEndpoint:    identityEndpoint,
		Username:            username,
		Password:            password,
		TenantID:            tenantID,
		DomainID:            domainID,
		FloatingNetworkID:   floatinNetID,
		FloatingNetworkName: floatNetName,
		AvailabilityZone:    availZone,
		StaticGateway:       staticgw,
	}
}

// Retrive Static GW info from this object, using the interface
func (client *OpenStackDriver) GetStaticGatewayInfo() (string, string, string, string, string, error) {
	log.Info("Retrieving Static Gateway Info from OpenStack...")
	if client.StaticGateway == nil {
		log.Error("No Static Gateway specified for VIM OpenStack...")
		return "", "", "", "", "", ErrStaticGatewayNotFound
	} else {
		return client.StaticGateway.NetworkName, client.StaticGateway.SubnetCidr, client.StaticGateway.ExternalInterfaceName, client.StaticGateway.FloatingIP, client.StaticGateway.InstanceID, nil
	}

}

// RetrieveFloatingNetworkID function
func (client *OpenStackDriver) RetrieveFloatingNetworkID() string {
	log.Info("Retrieve FloatingNetworkID for OpenStack...")
	return client.FloatingNetworkID
}

// RetrieveFloatingNetworkName function
func (client *OpenStackDriver) RetrieveFloatingNetworkName() string {
	log.Info("Retrieve FloatingNetworkName for OpenStack...")
	return client.FloatingNetworkName
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
	// TODO add when authenticated, otherwise SIGSEGV
	log.Info("Revoking token...")
	token, err := tokens.Revoke(client.identityClient, client.provider.Token()).Extract()
	if err != nil {
		log.Error(err)
	}
	log.Info(*token)
}
