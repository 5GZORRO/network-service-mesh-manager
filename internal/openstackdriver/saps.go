package openstackdriver

import (
	log "github.com/sirupsen/logrus"
)

func (obj *OpenStackDriver) CreateSAP(floatingNetName string, networkName string, cidr string) (string, string, string, string, string, string, error) {
	log.Info("Creating SAP...")
	return "", "", "", "", "", "", nil
}

func (obj *OpenStackDriver) DeleteSAP(networkID string, subnetID string, routerID string, portID string) error {
	log.Info("Deleting SAP...")
	return nil
}
