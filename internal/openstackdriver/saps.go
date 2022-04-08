package openstackdriver

import (
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	log "github.com/sirupsen/logrus"
)

func (client *OpenStackDriver) CreateSAP(floatingNetName string, networkName string, cidr string) (string, string, string, string, string, string, error) {
	log.Info("Creating SAP...")
	var floatingNetID string
	var routerAvailabilityZoneHints = []string{client.AvailabilityZone}
	var routerAdminStateUp = true

	// 1. Retrieve the floatingID of the floatingNetName
	floatingNetID, err := client.RetrieveFloatingNetworkByName(floatingNetName)
	if err != nil {
		log.Error("Error retrieving floating_net_ID")
		return "", "", "", "", "", "", ErrFloatingNotFound
	}
	log.Trace("Creating SAP: Retrieved floatingNetID: ", floatingNetID)
	var routerGatewayInfo = routers.GatewayInfo{
		NetworkID: floatingNetID,
	}

	// // 2. Create the private network with the subnet
	netID, subnetID, subnetName, err := client.CreateNetwork(networkName, cidr, true)
	if err != nil {
		log.Error("Error creating SAP: impossible to create the network")
		return "", "", "", "", "", "", ErrNetworkCreation
	}

	// // 3. Create the router
	routerName := networkName + "_router"
	createOpts := routers.CreateOpts{
		Name:                  routerName,
		AdminStateUp:          &routerAdminStateUp,
		AvailabilityZoneHints: routerAvailabilityZoneHints,
		GatewayInfo:           &routerGatewayInfo,
	}
	log.Trace("Creating SAP: Creating router with name: ", routerName)
	router, err := routers.Create(client.networkClient, createOpts).Extract()
	if err != nil {
		log.Error(err)
		return netID, subnetID, subnetName, "", "", "", ErrRouterCreation
	}
	log.Debug("Creating SAP: Creating router with name: ", routerName, " and ID: ", router.ID)

	// 4. Add interface to the router
	intOpts := routers.AddInterfaceOpts{
		SubnetID: subnetID,
	}

	log.Trace("Creating SAP: adding interface")
	port, err := routers.AddInterface(client.networkClient, router.ID, intOpts).Extract()
	if err != nil {
		log.Error(err)
		return netID, subnetID, subnetName, router.ID, routerName, "", ErrPortCreation
	}
	log.Debug("Creating SAP: added to router interface with ID: ", port.PortID)
	return netID, subnetID, subnetName, router.ID, routerName, port.PortID, nil
}

func (client *OpenStackDriver) DeleteSAP(networkID string, subnetID string, routerID string, portID string) error {
	log.Info("Deleting SAP...")

	// try to delete what exists
	if subnetID != "" && routerID != "" {
		// 1. delete port/interface of the router to the subnet
		if portID != "" {
			// Delete interface
			log.Trace("Deleting SAP: deleting router interface with ID: ", portID)
			intOpts := routers.RemoveInterfaceOpts{
				SubnetID: subnetID,
			}
			_, err := routers.RemoveInterface(client.networkClient, routerID, intOpts).Extract()
			if err != nil {
				log.Error("Error deleting SAP: error during interface removal of router ", routerID, " with subnetID "+subnetID)
				return ErrPortRemoval
			} else {
				log.Debug("Deleting SAP: delete router interface with ID: ", portID)
			}
		} else {
			log.Info("Empty portID parameter, impossible to delete the router interface")
		}
		// 2. delete the router
		log.Trace("Deleting SAP: deleting router with ID ", routerID)
		err := routers.Delete(client.networkClient, routerID).ExtractErr()
		if err != nil {
			log.Error("Error deleting SAP: error during router removal with ID " + routerID)
			return ErrRouterRemoval
		} else {
			log.Debug("Deleting SAP: deleted router with ID ", routerID)
		}
	} else {
		log.Info("Empty routerID or SubnetID parameters, impossible to delete the router")
	}
	// 3. delete network and subnet
	if networkID != "" {
		err := client.DeleteNetwork(networkID, subnetID)
		if err != nil {
			log.Error("Error deleting SAP: error during network removal with ID " + networkID)
			return err
		}
	}
	return nil
}
