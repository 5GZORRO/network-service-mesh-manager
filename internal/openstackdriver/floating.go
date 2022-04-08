package openstackdriver

import (
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	log "github.com/sirupsen/logrus"
)

// AllocateFloatingIP allocates a floatingIp identified to a (compute) port found on the SAP network, selected using the prefix network
func (client *OpenStackDriver) AllocateFloatingIP(networkid string) (string, string, string, string, error) {
	log.Trace("Allocating floatingIP for network ID", networkid)

	// 1. Retrieve the OpenStack port associated with the GW-VM interface to which associate a new FIP
	port, err := retrievePortGW(client.networkClient, networkid)
	if err != nil {
		return "", "", "", "", err
	}

	// 2. Request the creation of a new FIP
	createOpts := floatingips.CreateOpts{
		FloatingNetworkID: client.FloatingNetworkID,
		PortID:            port.ID,
	}

	fip, err := floatingips.Create(client.networkClient, createOpts).Extract()
	if err != nil {
		log.Error("Error creating floatingIP: ", err)
		return "", "", "", "", err
	}

	log.Trace("floatingIP  allocated with PortID: ", port.ID, " and floatingID: ", fip.ID)
	return port.ID, port.Name, fip.ID, fip.FloatingIP, nil
}

// DeallocateFloatingIP deallocates a floatingIp identified by the fipID from a port identified by the portID
func (client *OpenStackDriver) DeallocateFloatingIP(portID string, fipID string) error {
	log.Trace("Deallocating FloatingIP with ID: ", fipID, " from PortID: ", portID)

	// Disassociate a FloatingIP with a Port
	updateOpts := floatingips.UpdateOpts{
		PortID: new(string),
	}

	_, err := floatingips.Update(client.networkClient, fipID, updateOpts).Extract()
	if err != nil {
		log.Error("Error deallocating floatingIP: ", err)
		return err
	}

	// Deleting the FloatingIP
	err = floatingips.Delete(client.networkClient, fipID).ExtractErr()
	if err != nil {
		log.Error("Error deleting floatingIP: ", err)
		return err
	}
	log.Trace("FloatingIP with ID: ", fipID, " from PortID: ", portID, "deallocated and deleted correctly")
	return nil
}

func retrievePortGW(client *gophercloud.ServiceClient, networkid string) (*ports.Port, error) {
	log.Trace("Retrieving GW Port ID for networkID ", networkid)
	listOpts := ports.ListOpts{
		NetworkID: networkid,
	}

	allPages, err := ports.List(client, listOpts).AllPages()
	if err != nil {
		log.Error("GW SAP network port not found on OpenStack: error retriving ports")
		return nil, ErrGWPortNotFound
	}

	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		log.Error("GW SAP network port not found on OpenStack: error extracting ports")
		return nil, ErrGWPortNotFound
	}

	if len(allPorts) == 1 {
		log.Error("GW SAP network port not found on OpenStack: no ports retrieved with networkID ", networkid)
		return nil, ErrGWPortNotFound
	}

	for _, port := range allPorts {
		if strings.HasPrefix(port.DeviceOwner, "compute:") { //TODO, fixed string by config file
			log.Debug("GW Port ID found with ID ", port.ID)
			return &port, nil
		}
	}
	log.Error("GW SAP network port not found on OpenStack")
	return nil, ErrGWPortNotFound
}
