package openstackdriver

import (
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	log "github.com/sirupsen/logrus"
)

func (obj *OpenStackDriver) CreateNetwork(name string, cidr string) (string, string, string, error) {
	log.Info("Creating Network...")
	return "", "", "", nil
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

func (obj *OpenStackDriver) DeleteNetwork(id string) {
	log.Info("Deleting Network...")
}
