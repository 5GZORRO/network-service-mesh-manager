package nsm

import (
	"fmt"
	"net"
	gatewayconfig "nextworks/nsm/internal/gateway-config"
	"nextworks/nsm/internal/vim"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// deleteResources is a goroutine to delete in an async way all the network resources
func deleteResources(database *gorm.DB, vim *vim.VimDriver, res *ResourceSet) {
	time.Sleep(time.Second * 5)

	// try to delete all resources
	for _, net := range res.Networks {
		err := (*vim).DeleteNetwork(net.NetworkId, net.SubnetId)
		if err != nil {
			log.Error("Error deleting network name ", net.NetworkName)
		}
	}
	for _, sap := range res.Saps {
		err := (*vim).DeleteSAP(sap.NetworkId, sap.SubnetId, sap.RouterId, sap.RouterPortId)
		if err != nil {
			log.Error("Error deleting SAP with network name ", sap.NetworkName)
		}
	}

	// if removal from VIM is OK then delete it from DB
	result := database.Delete(&res)
	if result.Error != nil {
		log.Error("Error deleting network resource set with ID: ", res.ID, " and slice-id: ", res.SliceId, " from DB")
	}
}

// configureGateway is a goroutine to configure and start the VPN server in the gateway,
// using an HTTP client
func configureGateway(database *gorm.DB, res *ResourceSet) {
	time.Sleep(time.Second * 5)

	// TODO change mgmt port to string?
	// configure VM gateway, starting the VPN server
	client := gatewayconfig.New(net.ParseIP(res.Gateway.MgmtIp), fmt.Sprint(res.Gateway.MgmtPort))
	// TODO check parameters
	output := client.Start("", "", "")
	if output {
		res.Status = READY
	} else {
		res.Status = CONFIGURATION_ERROR
	}
	log.Trace("Update gateway and resource set state in DB")
	// update the state
	result := database.Save(&res)
	if result.Error != nil {
		log.Error("Error updating resource after gateway configuration in DB, resource set with ID: ", res.ID, " and slice-id: ", res.SliceId)
	}
}

// TODO resetGateway is a goroutine to configure the VM gateway, using
// an HTTP client
func resetGateway(database *gorm.DB, res *ResourceSet) {
	time.Sleep(time.Second * 10)
	// TODO reset VM gateway

	// update the state of the gateway to WAIT_FOR
	res.Status = WAIT_FOR_GATEWAY_CONFIG
	res.Gateway = Gateway{}
	result := database.Save(&res)
	if result.Error != nil {
		log.Error("Error updating resource set status with ID: ", res.ID, " and slice-id: ", res.SliceId)
	}
}
