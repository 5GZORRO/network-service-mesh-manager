package nsm

import (
	"fmt"
	"net"
	config "nextworks/nsm/internal/config"
	gatewayconfig "nextworks/nsm/internal/gateway-config"
	identityclient "nextworks/nsm/internal/identity"
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
// First if it is in Production mode - it requests the key pair then it invokes the VPNaaS Client to Launch the Server
//
func configureGateway(database *gorm.DB, res *ResourceSet, vpnaasenv string, idep *config.IdepConfigurations) {

	log.Trace("Async routine to configure gateway stared.")
	log.Info("Configuration is in ", vpnaasenv, " mode.")

	// TODO TO BE MODIFIED
	if vpnaasenv == gatewayconfig.Prod {
		log.Trace("Configuration is in PROD mode.")
		idepClient := identityclient.New(net.ParseIP(idep.Host), fmt.Sprint(idep.Port), idep.Secret)
		keyPair, err := idepClient.CreateKeyPair()
		if err != nil {
			log.Error("Error obtaining key pair from ID&P")
			res.Status = CONFIGURATION_ERROR
			result := database.Save(&res)
			if result.Error != nil {
				log.Error("Error updating resource after gateway configuration in DB, resource set with ID: ", res.ID, " and slice-id: ", res.SliceId)
			}
			log.Trace("Async routine to configure gateway ended")
		}
		log.Trace("KeyPair received: ", keyPair)
		// handle keyPair
		res.Gateway.Config.Keys.Did = keyPair.Did
		res.Gateway.Config.Keys.PrivK = keyPair.PrivKey
		res.Gateway.Config.Keys.PubK = keyPair.PubKey
		res.Gateway.Config.Keys.Timestamp = keyPair.Timestamp

		// configure VM gateway, starting the VPN server
		// var client gatewayconfig.VPNHttpClient
		// TODO VPNaaS client should be updated with keypair passed as param if PROD, otherwise no additional config
		client := gatewayconfig.New(net.ParseIP(res.Gateway.Config.MgmtIp), fmt.Sprint(res.Gateway.Config.MgmtPort), vpnaasenv)

		vpnIp := res.Gateway.Config.PrivateVpnRange
		log.Trace(keyPair)
		output := client.Launch(vpnIp, res.Gateway.External.PortName, fmt.Sprint(res.Gateway.Config.MgmtPort), keyPair)
		log.Debug(output)
		if output {
			res.Status = READY
		} else {
			res.Status = CONFIGURATION_ERROR
		}
	} else {
		log.Trace("Configuration is in TEST mode. Key pair is not configured")
		// No additional configuration needed

		// configure VM gateway, starting the VPN server
		// var client gatewayconfig.VPNHttpClient
		// TODO VPNaaS client should be updated with keypair passed as param if PROD, otherwise no additional config
		client := gatewayconfig.New(net.ParseIP(res.Gateway.Config.MgmtIp), fmt.Sprint(res.Gateway.Config.MgmtPort), vpnaasenv)

		vpnIp := res.Gateway.Config.PrivateVpnRange
		output := client.Launch(vpnIp, res.Gateway.External.PortName, fmt.Sprint(res.Gateway.Config.MgmtPort), nil)
		log.Debug(output)
		if output {
			res.Status = READY
		} else {
			res.Status = CONFIGURATION_ERROR
		}
	}
	// Parte del codice sopra e' replicato sopra, clean up

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
func resetGateway(database *gorm.DB, res *ResourceSet) {
	log.Trace("Async routine to reset gateway stared")

	// update the state of the gateway to WAIT_FOR
	res.Status = WAIT_FOR_GATEWAY_CONFIG
	res.Gateway.Config = Config{}
	result := database.Save(&res)
	if result.Error != nil {
		log.Error("Error updating resource set status with ID: ", res.ID, " and slice-id: ", res.SliceId)
	}
	log.Trace("Async routine to reset gateway ended")
}
