package openstackdriver

import (
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	log "github.com/sirupsen/logrus"
)

// Functions used in case of StaticGW, after adding the static GW info in the SAP entity
// the GW needs to be attached on the internal_net of the network service, the one that
// wil be exposed on the VPN connection

func (client *OpenStackDriver) CreateInterfacePort(serverId string, networkId string) (string, error) {
	log.Info("Creating interface Port from ServerId on OpenStack...")
	var portSecurity = false

	// 1. First create the port
	attachOpts := attachinterfaces.CreateOpts{
		NetworkID: networkId,
	}
	interfacePort, err := attachinterfaces.Create(client.computeClient, serverId, attachOpts).Extract()
	if err != nil {
		log.Error("Impossible to create ad interface on Server Instance ID ", serverId, " to network ID ", networkId)
		return "", nil
	}
	// 2. Then, disable port security
	var portWithPortSecurityExtensions struct {
		ports.Port
		portsecurity.PortSecurityExt
	}

	portUpdateOpts := ports.UpdateOpts{}
	updateOpts := portsecurity.PortUpdateOptsExt{
		UpdateOptsBuilder:   portUpdateOpts,
		PortSecurityEnabled: &portSecurity,
	}

	// NOTE: non blocking error
	err = ports.Update(client.networkClient, interfacePort.PortID, updateOpts).ExtractInto(&portWithPortSecurityExtensions)
	if err != nil {
		log.Error("Impossible to disable port security on interface created on GW-VM")
		return "", nil
	}

	return interfacePort.PortID, nil
}

func (client *OpenStackDriver) DeleteInterfacePort(serverId string, portId string) error {
	log.Info("Deleting interface Port from ServerId on OpenStack...")
	err := attachinterfaces.Delete(client.computeClient, serverId, portId).ExtractErr()
	if err != nil {
		log.Error("Impossible to delete ad interface on Server Instance ID ", serverId, " interface Port with ID ", portId)
		return err
	}
	return nil
}
