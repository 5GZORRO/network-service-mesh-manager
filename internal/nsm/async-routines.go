package nsm

import (
	"fmt"
	"net"
	gatewayconfig "nextworks/nsm/internal/gateway-config"
	vimdriver "nextworks/nsm/internal/vim"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// deleteResources is a goroutine to delete in an async way all the network resources
func deleteResources(database *gorm.DB, vim *vimdriver.VimDriver, res *ResourceSet) {
	// time.Sleep(time.Second * 5)

	log.Trace("Async routine to delete resources stared")
	if len(res.Networks) == 0 {
		log.Error("No networks to be deleted found")
	}
	// try to delete all resources
	for _, net := range res.Networks {
		err := (*vim).DeleteNetwork(net.NetworkId, net.SubnetId)
		if err != nil {
			log.Error("Error deleting network name ", net.NetworkName)
		}
	}
	if len(res.Saps) == 0 {
		log.Error("No SAPs to be deleted found")
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
	log.Trace("Async routine to delete resources ended")
}

// configureGateway is a goroutine to configure and start the VPN server in the gateway,
// using an HTTP client
func configureGateway(database *gorm.DB, res *ResourceSet) {

	log.Trace("Async routine to configure gateway stared")
	// configure VM gateway, starting the VPN server
	client := gatewayconfig.New(net.ParseIP(res.Gateway.MgmtIp), fmt.Sprint(res.Gateway.MgmtPort))

	vpnIp := res.Gateway.PrivateVpnRange
	output := client.Launch(vpnIp, res.Gateway.PortName, fmt.Sprint(res.Gateway.MgmtPort))
	log.Debug(output)
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
	log.Trace("Async routine to configure gateway ended")
}

// resetGateway is a goroutine to reset the VM gateway, using an HTTP client
// TODO it should reset also the VPNaaS, which does not have this functionality now
func resetGateway(database *gorm.DB, vim *vimdriver.VimDriver, res *ResourceSet) {
	log.Trace("Async routine to reset gateway stared")

	// deallocate FIP, if it exists
	if res.Gateway.FloatingID != "" {
		log.Trace("Deallocating and deleting floatingIP ")
		err := (*vim).DeallocateFloatingIP(res.Gateway.PortID, res.Gateway.FloatingID)
		if err != nil {
			log.Error("Error deallocating/deleting floating IP with ID: ", res.Gateway.FloatingID)
		}
	} else {
		log.Trace("No floatingIP to be deallocated/deleted")
	}

	// update the state of the gateway to WAIT_FOR
	res.Status = WAIT_FOR_GATEWAY_CONFIG
	res.Gateway = Gateway{}
	result := database.Save(&res)
	if result.Error != nil {
		log.Error("Error updating resource set status with ID: ", res.ID, " and slice-id: ", res.SliceId)
	}
	log.Trace("Async routine to reset gateway ended")
}
