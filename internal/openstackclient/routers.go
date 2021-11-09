package openstackclient

import (
	// "errors"
	// "github.com/gophercloud/gophercloud"
	// "github.com/gophercloud/gophercloud/openstack"
	// "github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	// "github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	log "github.com/sirupsen/logrus"
)

// TODO
// CreateRouters creates a router on the OpenStack TenantID and
// connects it to a private_net passed as params
func (client *OpenStackClient) CreateRouter(privnet string) error {
	log.Info("Creating router")
	return nil
}

// TODO
func (client *OpenStackClient) DeleteRouter(name string) error {
	log.Info("Creating router")
	return nil
}
