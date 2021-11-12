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

// TODO
// CreateRouters creates a router on the OpenStack TenantID and
// connects it to a private_net passed as params
func (client *OpenStackClient) CreateRouter(routerName string) (*routers.Router, error) {
	log.Info("Creating router...")

	createOpts := routers.CreateOpts{
		Name:                  routerName,
		AdminStateUp:          &routerAdminStateUp,
		AvailabilityZoneHints: routerAvailabilityZoneHints,
		GatewayInfo:           &routerGatewayInfo,
	}

	router, err := routers.Create(client.networkClient, createOpts).Extract()
	if err != nil {
		return nil, err
	}
	return router, nil
}

// TODO
func (client *OpenStackClient) DeleteRouter(routerID string) error {
	log.Info("Deleting router")
	err := routers.Delete(client.networkClient, routerID).ExtractErr()
	if err != nil {
		return err
	}
	return nil
}
