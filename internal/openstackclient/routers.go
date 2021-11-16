package openstackclient

import (
	// "errors"
	// "github.com/gophercloud/gophercloud"
	// "github.com/gophercloud/gophercloud/openstack"
	// "github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	// "github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	log "github.com/sirupsen/logrus"
)

var floatingNetID = "85e329ed-1bed-4bb6-8b1c-a11a7eb133fa"
var routerAvailabilityZoneHints = []string{"nova"}
var routerAdminStateUp = true
var routerGatewayInfo = routers.GatewayInfo{
	NetworkID: floatingNetID,
}

// CreateRouters creates a router on the OpenStack TenantID and
// connects it to a private_net passed as params
func (client *OpenStackClient) CreateRouter(routerName string, subnetID string) (*routers.Router, *routers.InterfaceInfo, error) {
	log.Info("Creating router...")

	createOpts := routers.CreateOpts{
		Name:                  routerName,
		AdminStateUp:          &routerAdminStateUp,
		AvailabilityZoneHints: routerAvailabilityZoneHints,
		GatewayInfo:           &routerGatewayInfo,
	}

	router, err := routers.Create(client.networkClient, createOpts).Extract()
	if err != nil {
		return nil, nil, err
	}

	// Add interface to router connecting it to the privatesubnet
	intOpts := routers.AddInterfaceOpts{
		SubnetID: subnetID,
	}

	port, err := routers.AddInterface(client.networkClient, router.ID, intOpts).Extract()
	if err != nil {
		return nil, nil, err
	}
	log.Info("Added interface to router :", port.PortID)
	return router, port, nil
}

// RetrieveRouter retrieve info associated to a router
// passing as params its ID
func (client *OpenStackClient) RetrieveRouterById(routerID string) (*routers.Router, error) {
	router, err := routers.Get(client.networkClient, routerID).Extract()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return router, nil
}

// DeleteRouter deletes the router passed as params removing before its interface
func (client *OpenStackClient) DeleteRouter(routerID string, subnetID string) error {

	if subnetID != "" {
		log.Info("Deleting router with ID ", routerID, " with port on subnetID ", subnetID)
		_ = client.deleteInterface(routerID, subnetID)
	} else {
		log.Info("Deleting router with ID ", routerID, " - No port to be deleted")
	}
	// delete router
	err := routers.Delete(client.networkClient, routerID).ExtractErr()
	if err != nil {
		log.Error("Error deleting router with ID " + routerID)
		return err
	}
	return nil
}

// DeleteRouter deletes the router passed as params removing before its interface
func (client *OpenStackClient) deleteInterface(routerID string, subnetID string) error {
	log.Info("Deleting router interface of router ", routerID)

	// Delete interface
	intOpts := routers.RemoveInterfaceOpts{
		SubnetID: subnetID,
	}

	_, err := routers.RemoveInterface(client.networkClient, routerID, intOpts).Extract()
	if err != nil {
		log.Error("Error deleting interface router with subnetID " + subnetID)
		return err
	}
	return nil
}
