package nsm

import (
	"errors"
	"net"
	"net/http"
	nsmmapi "nextworks/nsm/api"
	vimdriver "nextworks/nsm/internal/vim"
	"strconv"
	"strings"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (obj *ServerInterfaceImpl) GetNetResourcesIdGateway(c *gin.Context, id int) {
	// Retrive the Resource Set state and check it, if the gateway can be configured

	log.Trace("GetNetResourcesIdGateway - requested retrieve of gateway configuration for resource set with ID: ", id)
	resource, error := RetrieveResourcesFromDB(obj.DB, id)
	if error != nil {
		log.Error("Impossible to retrieve gateway configuration. Error reading from DB: ", error)
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, http.StatusNotFound, ErrResourcesNotExists)
			return
		} else {
			SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
			return
		}
	}

	// check if it is configured
	if resource.Gateway.ExternalIp == "" {
		log.Error("Impossible to retrieve gateway configuration. It does not exist")
		SetErrorResponse(c, http.StatusNotFound, ErrGatewayNotConfigured)
		return
	}
	SetGatewayResponse(c, http.StatusOK, *resource)
}

func (obj *ServerInterfaceImpl) PutNetResourcesIdGateway(c *gin.Context, id int) {
	var jsonBody nsmmapi.PostGateway

	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		log.Error("Impossible to create a gateway configuration. Error in the request, wrong json body")
		SetErrorResponse(c, http.StatusBadRequest, ErrRequestConfigurationGateway)
		return
	}

	log.Trace("PutNetResourcesIdGateway - requested configuration of gateway for resource set with ID: ", id)
	resource, error := RetrieveResourcesFromDB(obj.DB, id)
	if error != nil {
		log.Error("Impossible to create gateway configuration. Error reading from DB: ", error)
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, http.StatusNotFound, ErrResourcesNotExists)
			return
		} else {
			SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
			return
		}
	}
	// check status
	// TODO this could be also configuring, if a floating ip has been added previously
	if resource.Status != WAIT_FOR_GATEWAY_CONFIG && resource.Status != CONFIGURATION_ERROR {
		log.Error("Impossibile to create gateway configuration. The current state is ", resource.Status)
		SetErrorResponse(c, http.StatusForbidden, ErrConfiguringGateway)
		return
	}

	if err := checkGatewayConfigurationParams(jsonBody); err != nil {
		log.Error("Impossible to create gateway configuration - error in json body ", err)
		SetErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	// retrieve VIM
	vim, err := obj.Vims.GetVim(resource.VimName)
	if err == vimdriver.ErrVimNotFound {
		log.Error("Error: VIM of resource set with id ", resource.ID, " does not exist")
		SetErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	//  SET CONFIGURING STATE AND SAVE
	resource.Status = CONFIGURING
	log.Trace("Start association of a floating IPs - updating state to CONFIGURINF of resource set with ID: ", id)
	output := obj.DB.Save(&resource)
	if output.Error != nil {
		log.Error("Impossible to update resource set - error saving in DB ", output.Error)
		SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
		return
	}

	// IF no FloatingIP is associated to the Gateway, continue allocating one
	if resource.Gateway.ExternalIp == "" {
		log.Debug("Gateway without an external-ip, allocating it... ")

		// ASSOCIATE FLOATING-IP USING VIM
		// For each SAP network, try to allocate a floating IP for each compute, found on that network
		// retrieve all SAP networks of this resource set
		err = LoadSAPsFromDB(obj.DB, resource)
		if err != nil {
			resource.Status = CONFIGURATION_ERROR
			_ = obj.DB.Save(&resource)
			log.Error("Error retrieving associations of resource set with ID: ", resource.ID)
			SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
			return
		}
		log.Trace("Loaded SAP: ", resource.Saps)
		//for the SAP matching the NetworkPrefix search a VM and allocate a floating-ip
		var selectedSap *Sap = nil
		for _, sap := range resource.Saps {
			log.Trace("Loaded SAP: ", sap.NetworkName)
			if strings.HasPrefix(sap.NetworkName, obj.Netconfig.GatewayNetworkNamePrefix) {
				log.Info("Selected SAP network with name ", sap.NetworkName)
				selectedSap = &sap
				break
			}
		}
		if selectedSap == nil {
			// TODO error
			log.Error("Error no SAP network found with the specified prefix-name ", obj.Netconfig.GatewayNetworkNamePrefix, " for resource set with ID: ", resource.ID)
			resource.Status = CONFIGURATION_ERROR
			_ = obj.DB.Save(&resource)
			SetErrorResponse(c, http.StatusInternalServerError, ErrGatewayNoNetworkFound)
			return
		} else {
			portid, portname, fipid, fip, err := (*vim).AllocateFloatingIP(selectedSap.NetworkId)
			if err != nil {
				log.Error("Error allocating FloatingIP for gatewa of resource set with ID: ", resource.ID)
				resource.Status = CONFIGURATION_ERROR
				result := obj.DB.Save(&resource)
				if result.Error != nil {
					log.Error("Error updating resource with ID: ", resource.ID)
				}
				SetErrorResponse(c, http.StatusInternalServerError, err)
				return
			} else {
				log.Trace(" PortID: ", portid, " PortName: ", portname, " FloatingID: ", fipid, " FloatingIP: ", fip)
				resource.Gateway.ExternalIp = fip
				resource.Gateway.PortID = portid
				resource.Gateway.PortName = portname
				resource.Gateway.FloatingID = fipid
			}
		}
	} else {
		// FloatingIP is already configured, we could be in CONFIGURATION_ERROR
		log.Debug("Gateway has already an external-ip, skip the allocation phase ")
	}

	resource.Gateway.MgmtIp = jsonBody.MgmtIp
	resource.Gateway.MgmtPort, _ = parsePort(strconv.Itoa(int(obj.VpnaasConfig.VpnaasPort)))
	resource.Gateway.ExposedNets = SubnetsToString(jsonBody.SubnetToExpose)

	// ranges and private ips
	_, privNet, _ := net.ParseCIDR(jsonBody.PrivateVpnRange)
	resource.Gateway.PrivateVpnRange = privNet.String()
	log.Info("Setting private VPN range as ", resource.Gateway.PrivateVpnRange)
	// TODO to check
	if jsonBody.PrivateVpnPeerIp != nil {
		peerIp := net.ParseIP(*jsonBody.PrivateVpnPeerIp)
		resource.Gateway.PrivateVpnIp = cidr.Inc(peerIp).String()
		log.Info("Setting private VPN IP as ", resource.Gateway.PrivateVpnIp, " knowing the IP of the peer")
	} else {
		resource.Gateway.PrivateVpnIp = cidr.Inc(privNet.IP).String()
		log.Info("Setting private VPN IP as ", resource.Gateway.PrivateVpnIp)
	}

	// Updating other fields
	log.Trace("Creating gateway configuration - updating network resource set with ID: ", id)
	output = obj.DB.Save(&resource)
	if output.Error != nil {
		log.Error("Impossible to create gateway configuration - error saving in DB when updating network resource set  ", output.Error)
		SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
		return
	}
	// go routine with httpclient to configure the VPN server
	// and the update the state to -> READY
	go configureGateway(obj.DB, resource)

	SetGatewayResponse(c, http.StatusCreated, *resource)
}

// Delete gateway configuration
func (obj *ServerInterfaceImpl) DeleteNetResourcesIdGateway(c *gin.Context, id int) {
	log.Trace("DeleteNetResourcesIdGateway - requested removal of gateway configuration for resource set with ID: ", id)
	resource, error := RetrieveResourcesFromDB(obj.DB, id)
	if error != nil {
		log.Error("Impossible to delete gateway configuration. Error reading from DB: ", error)
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, http.StatusNotFound, ErrResourcesNotExists)
			return
		} else {
			SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
			return
		}
	}
	// check status
	// READY state could be excluded in this check
	if resource.Status != READY && resource.Status != CONFIGURATION_ERROR {
		log.Error("Impossibile to delete gateway configuration. The current state is ", resource.Status)
		SetErrorResponse(c, http.StatusForbidden, ErrDeleteConfigurationGateway)
		return
	}

	// retrieve VIM
	vim, err := obj.Vims.GetVim(resource.VimName)
	if err == vimdriver.ErrVimNotFound {
		log.Error("Error: VIM of resource set with id ", resource.ID, " does not exist")
		SetErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	// If it is configurable update the state to -> CONFIGURING
	resource.Status = DELETING_CONFIGURATION
	output := obj.DB.Save(&resource)
	if output.Error != nil {
		log.Error("Impossible to delete gateway configuration. Error saving in DB: ", output.Error)
		SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
		return
	}
	// TODO creates go routine with httpclient to reset VPN server?
	go resetGateway(obj.DB, vim, resource)

	// and the update the DB with nil param
	c.Status(http.StatusNoContent)
}
