package openstackdriver

import (
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	log "github.com/sirupsen/logrus"
)

func (client *OpenStackDriver) CreateNetwork(name string, cidr string) (string, string, string, error) {
	log.Info("Creating Network...")
	// network global params
	var sharedNetworks bool = false
	var availabilityZoneHints = []string{"nova"}
	var adminStateUp = true

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
		return "", "", "", err
	}
	// Create subnet
	subnetName := name + "_subnet"
	log.Info("Creating Subnet " + subnetName)
	subnet, err := client.createSubnet(subnetName, network.ID, cidr)
	if err != nil {
		log.Error("Error creating Subnet " + subnetName)
		return network.ID, "", "", err
	}

	return network.ID, subnet.ID, subnetName, nil
}

// Function to create a Subnet
func (client *OpenStackDriver) createSubnet(name string, networkID string, cidr string) (*subnets.Subnet, error) {
	var enableDHCP = true
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

func (client *OpenStackDriver) RetrieveNetwork(id string) {
	log.Info("Retrieving Network...")
	log.Info("Authentication Token String used: ", client.provider.Token())

	network, err := networks.Get(client.networkClient, id).Extract()
	if err != nil {
		log.Error("Retriving Network by ID: ", id, err)
	}
	log.Info(network)
}

func (client *OpenStackDriver) DeleteNetwork(networkID string, subnetID string) error {
	log.Info("Deleting Network with ID ", networkID)
	err := networks.Delete(client.networkClient, networkID).ExtractErr()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Network with ID" + networkID + " deleted ")
	return nil
}
