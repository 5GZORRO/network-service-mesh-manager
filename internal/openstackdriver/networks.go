package openstackdriver

import (
	"errors"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	log "github.com/sirupsen/logrus"
)

var networkWithPortSecurityExt struct {
	networks.Network
	portsecurity.PortSecurityExt
}

func (client *OpenStackDriver) CreateNetwork(name string, cidr string, gateway bool) (string, string, string, error) {
	log.Trace("Creating Network with name ", name)
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
	log.Debug("Created network with name: ", network.Name, " and ID: ", network.ID)

	// Create subnet
	subnetName := name + "_subnet"
	log.Trace("Creating Subnet " + subnetName)
	subnet, err := client.createSubnet(subnetName, network.ID, cidr, gateway)
	if err != nil {
		log.Error("Error creating Subnet " + subnetName)
		return network.ID, "", "", err
	}
	log.Debug("Created subnet with name: ", subnet.Name, " and ID: ", subnet.ID)

	// TODO not all the networks require the port security disabled (?)
	// only the one exposed

	// disable port security on an existing network
	log.Debug("Disabling port security on network ", network.Name)
	var portSecurity bool = false
	networkUpdateOpts := networks.UpdateOpts{}
	updateOpts := portsecurity.NetworkUpdateOptsExt{
		UpdateOptsBuilder:   networkUpdateOpts,
		PortSecurityEnabled: &portSecurity,
	}

	err = networks.Update(client.networkClient, network.ID, updateOpts).ExtractInto(&networkWithPortSecurityExt)
	if err != nil {
		log.Error("Error disabling port security on network " + network.Name)
		return network.ID, subnet.ID, subnet.Name, err
	}

	return network.ID, subnet.ID, subnet.Name, nil
	// return "", "", "", nil
}

// Function to create a Subnet
func (client *OpenStackDriver) createSubnet(name string, networkID string, cidr string, gateway bool) (*subnets.Subnet, error) {
	var enableDHCP = true
	var gatewayIP *string = nil
	if !gateway {
		// case where we dont want a default gateway
		str := ""
		gatewayIP = &str
	}
	createOpts := subnets.CreateOpts{
		NetworkID:  networkID,
		Name:       name,
		TenantID:   client.TenantID,
		GatewayIP:  gatewayIP, // handle default gateway
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

// RetrieveNetwork retrieves a Network by its name
func (client *OpenStackDriver) RetrieveFloatingNetworkByName(name string) (string, error) {
	sharedNetworks := true
	listOpts := networks.ListOpts{
		Name:   name,
		Shared: &sharedNetworks,
	}

	allPages, err := networks.List(client.networkClient, listOpts).AllPages()
	if err != nil {
		log.Error(err)
		return "", err
	}

	pages, _ := allPages.IsEmpty()
	if !pages {
		allNetworks, err := networks.ExtractNetworks(allPages)
		if err != nil {
			log.Error(err)
			return "", err
		}
		if len(allNetworks) > 1 {
			log.Error("More than one Network with name " + name)
			return "", errors.New("More than one network with name " + name)
		}
		// return first network
		return allNetworks[0].ID, nil
	} else {
		// Network not found
		log.Error("Network " + name + " not found")
		return "", errors.New("Network " + name + " not found")
	}
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
	if networkID != "" {
		log.Trace("Deleting Network with ID ", networkID)
		err := networks.Delete(client.networkClient, networkID).ExtractErr()
		if err != nil {
			log.Error(err)
			return ErrNetworkRemoval
		}
		log.Debug("Deleted network with ID: " + networkID)
		return nil
	} else {
		log.Info("Empty networkID parameter, impossible to delete the network")
		return nil
	}
	// return nil
}
