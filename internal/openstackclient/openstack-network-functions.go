package openstackclient

import (
	"errors"
	// "github.com/gophercloud/gophercloud"
	// "github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	log "github.com/sirupsen/logrus"
)

// network global params
var sharedNetworks bool = false
var availabilityZoneHints = []string{"nova"}
var adminStateUp = true

// subnet global params
var enableDHCP = true

// Function to create a Network and its subnet
func (client *OpenStackClient) CreateNetwork(name string, cidr string) (*networks.Network, *subnets.Subnet, error) {
	createOpts := networks.CreateOpts{
		Name:                  name,
		AdminStateUp:          &adminStateUp,
		Shared:                &sharedNetworks,
		TenantID:              client.TenantID,
		AvailabilityZoneHints: availabilityZoneHints,
	}

	network, err := networks.Create(client.networkClient, createOpts).Extract()
	if err != nil {
		log.Error("Error creating Network " + name)
		return nil, nil, err
	}
	// Create subnet
	subnetName := name + "_subnet"
	log.Info("Creating Subnet " + subnetName)
	subnet, err := client.createSubnet(subnetName, network.ID, cidr)
	if err != nil {
		log.Error("Error creating Subnet " + subnetName)
		return nil, nil, err
	}
	return network, subnet, nil
}

// Function to create a Subnet
func (client *OpenStackClient) createSubnet(name string, networkID string, cidr string) (*subnets.Subnet, error) {
	createOpts := subnets.CreateOpts{
		NetworkID:  networkID,
		Name:       name,
		TenantID:   client.TenantID,
		EnableDHCP: &enableDHCP,
		IPVersion:  4,
		CIDR:       cidr,
	}

	subnet, err := subnets.Create(client.networkClient, createOpts).Extract()
	if err != nil {
		log.Error("Error creating Subnet " + name)
		return nil, err
	}
	return subnet, nil

}

// Function to retrieve a Network by name
func (client *OpenStackClient) RetrieveNetwork(name string) (*networks.Network, error) {
	sharedNetworks := false
	listOpts := networks.ListOpts{
		TenantID: client.TenantID,
		Name:     name,
		Shared:   &sharedNetworks,
	}

	allPages, err := networks.List(client.networkClient, listOpts).AllPages()
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
func (client *OpenStackClient) DeleteNetwork(name string) error {
	network, err := client.RetrieveNetwork(name)
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
		err = client.deleteSubnet(network.Subnets[0])
		if err != nil {
			log.Error("Error deleting subnet")
			return err
		}
	} else {
		log.Error("Network "+name+" has %d subnets", numSubnets)
		return errors.New("expected exactly one subnet in the network")
	}
	err = networks.Delete(client.networkClient, network.ID).ExtractErr()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Network " + name + " deleted ")
	return nil
}

func (client *OpenStackClient) deleteSubnet(id string) error {
	err := subnets.Delete(client.networkClient, id).ExtractErr()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("SubnetID " + id + " deleted ")
	return nil
}

// TODO
func (client *OpenStackClient) CreateRouter() error {
	log.Info("Creating router")
	return nil
}

// TODO
func (client *OpenStackClient) DeleteRouter() error {
	log.Info("Creating router")
	return nil
}
