package openstackclient

import (
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	log "github.com/sirupsen/logrus"
)

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

// RetrieveSubnet retrieves the subnet associated to the id
// passed as params
func (client *OpenStackClient) RetrieveSubnetById(id string) (*subnets.Subnet, error) {
	subnet, err := subnets.Get(client.networkClient, id).Extract()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return subnet, nil
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
