package openstackclient

import (
	"errors"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	log "github.com/sirupsen/logrus"
)

// RetriveGatewayFloatingIP retrieves the floating IP associated to a VM (intended as the gateway VM) on a network.
// It takes as parameters the VMID and the network NAME (not ID) to retrieve the floating IP on this network
// NOTE: the VM is deployed by OSM connected with 2 networks: the private and the ones to be able to assign FIP
// the floating IP is assumed allocated by OSM
func (client *OpenStackClient) RetrieveGatewayFloatingIP(vmId string, networkName string) (string, error) {

	// We need the UUID in string form
	server, err := servers.Get(client.computeClient, vmId).Extract()
	if err != nil {
		log.Error(err)
		return "", err
	}

	list, exists := server.Addresses[networkName]
	if !exists {
		log.Error(err)
		return "", err
	}
	log.Info(list)
	// check type assertion
	test, ok := list.([]interface{})
	if !ok {
		log.Error("list is not type []interface{}")
		return "", errors.New("error retrieving FloatingIP")
	}

	// search for a floating IP
	for i := range test {
		networkInfo := test[i].(map[string]interface{})
		ipType := networkInfo["OS-EXT-IPS:type"]
		log.Info(ipType)
		if ipType == "floating" {
			floatingIp, ok := networkInfo["addr"]
			if !ok {
				log.Error("floating IP address not found on a floatingIP")
				return "", errors.New("floating IP address not found on a floatingIP")
			}
			log.Info("Found floatingIP ", floatingIp)
			return floatingIp.(string), nil
		}
	}
	return "", errors.New("No Floating IP found on network " + networkName)
}
