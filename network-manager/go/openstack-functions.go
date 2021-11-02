package main

import (
	"errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	log "github.com/sirupsen/logrus"
)

// Global variables to access OpenStack API
var Provider *gophercloud.ProviderClient
var IdentityClient *gophercloud.ServiceClient
var NetworkClient *gophercloud.ServiceClient

// network global params
var sharedNetworks bool = false
var availabilityZoneHints = []string{"nova"}
var adminStateUp = true

// subnet global params
var enableDHCP = true

func Init() {
	log.Info("Init function")
	var err error

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://10.30.7.10:5000/v3",
		Username:         "timeo",
		Password:         "nextworks",
		TenantID:         TenantID,
		DomainID:         "default",
	}

	Provider, err = openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Authentication Token: " + Provider.TokenID)
	IdentityClient, err = openstack.NewIdentityV3(Provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Identity Endpoint: " + IdentityClient.IdentityEndpoint)

	NetworkClient, err = openstack.NewNetworkV2(Provider, gophercloud.EndpointOpts{
		Name:   "neutron",
		Region: "RegionOne",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Network Endpoint " + NetworkClient.Endpoint)
}

// Function to create a Network and its subnet
func CreateNetwork(name string, cidr string) (*networks.Network, *subnets.Subnet, error) {
	createOpts := networks.CreateOpts{
		Name:                  name,
		AdminStateUp:          &adminStateUp,
		Shared:                &sharedNetworks,
		TenantID:              TenantID,
		AvailabilityZoneHints: availabilityZoneHints,
	}

	network, err := networks.Create(NetworkClient, createOpts).Extract()
	if err != nil {
		log.Error("Error creating Network " + name)
		return nil, nil, err
	}
	// Create subnet
	subnetName := name + "_subnet"
	log.Info("Creating Subnet " + subnetName)
	subnet, err := CreateSubnet(subnetName, network.ID, cidr)
	if err != nil {
		log.Error("Error creating Subnet " + subnetName)
		return nil, nil, err
	}
	return network, subnet, nil
}

// Function to create a Subnet
func CreateSubnet(name string, networkID string, cidr string) (*subnets.Subnet, error) {
	createOpts := subnets.CreateOpts{
		NetworkID:  networkID,
		Name:       name,
		TenantID:   TenantID,
		EnableDHCP: &enableDHCP,
		IPVersion:  4,
		CIDR:       cidr,
	}

	subnet, err := subnets.Create(NetworkClient, createOpts).Extract()
	if err != nil {
		log.Error("Error creating Subnet " + name)
		return nil, err
	}
	return subnet, nil

}

// Function to retrieve a Network by name
func RetrieveNetwork(name string) (*networks.Network, error) {
	sharedNetworks := false
	listOpts := networks.ListOpts{
		TenantID: TenantID,
		Name:     name,
		Shared:   &sharedNetworks,
	}

	allPages, err := networks.List(NetworkClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	pages, _ := allPages.IsEmpty()
	if !pages {
		allNetworks, err := networks.ExtractNetworks(allPages)
		if err != nil {
			panic(err)
		}
		log.Info("Networks: ", len(allNetworks))

		if len(allNetworks) > 1 {
			log.Error("More than one Network with name " + name)
			return nil, errors.New("More than one network with name " + name)
		}
		// return first network
		return &allNetworks[0], nil
	} else {
		// Network not found
		log.Error("Network " + name + " not found")
		return nil, errors.New("Network " + name + " not found")
	}
}

// Function to delete a network by name (not ID) and its subnet,
// assuming only one subnet
func DeleteNetwork(name string) error {
	network, err := RetrieveNetwork(name)
	if err != nil {
		return err
	}
	log.Info("Network " + name + " has ID " + network.ID)
	//
	numSubnets := len(network.Subnets)
	if numSubnets == 0 {
		log.Info("No subnets to be deleted")
	} else if numSubnets == 1 {
		log.Info("1 Subnet to be deleted: " + network.Subnets[0])
		err = DeleteSubnet(network.Subnets[0])
		if err != nil {
			log.Error("Error deleting subnet")
			return err
		}
	} else {
		log.Error("Network "+name+" has %d subnets", numSubnets)
		return errors.New("expected exactly one subnet in the network")
	}
	err = networks.Delete(NetworkClient, network.ID).ExtractErr()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Network " + name + " deleted ")
	return nil
}

func DeleteSubnet(id string) error {
	err := subnets.Delete(NetworkClient, id).ExtractErr()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("SubnetID " + id + " deleted ")
	return nil
}

// Function to "close the connection with OS" which is
// revoking the token created in the initial phase
func Close() {
	// Revoke token
	log.Info("Revoking Token...")
	token, err := tokens.Revoke(IdentityClient, Provider.TokenID).Extract()
	if err != nil {
		log.Error(err)
	}
	log.Info(*token)

	// TBT: token seems still valid with a GET /networks/
}
