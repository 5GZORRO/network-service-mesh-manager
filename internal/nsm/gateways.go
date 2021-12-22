package nsm

import (
	"errors"
	"net"
	"net/http"
	nsmmapi "nextworks/nsm/api"

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

// TODO added params in the NBI
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
	if resource.Status != WAIT_FOR_GATEWAY_CONFIG {
		log.Error("Impossibile to create gateway configuration. The current state is ", resource.Status)
		SetErrorResponse(c, http.StatusForbidden, ErrConfiguringGateway)
		return
	}

	if err := checkGatewayConfigurationParams(jsonBody); err != nil {
		log.Error("Impossible to create gateway configuration - error in json body ", err)
		SetErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	resource.Gateway.ExternalIp = jsonBody.ExternalIp
	resource.Gateway.MgmtIp = jsonBody.MgmtIp
	resource.Gateway.MgmtPort, _ = parsePort(jsonBody.MgmtPort)
	resource.Gateway.ExposedNets = SubnetsToString(jsonBody.SubnetToExpose)
	resource.Gateway.PubKey = jsonBody.PubKey

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

	// If it is configurable update the state to -> CONFIGURING
	// and store params
	resource.Status = CONFIGURING
	log.Trace("Creating gateway configuration - updating network resource set with ID: ", id)
	err := obj.DB.Save(&resource)
	if err.Error != nil {
		log.Error("Impossible to create gateway configuration - error saving in DB when updating network resource set  ", err.Error)
		SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
		return
	}
	// TODO go routine with httpclient to configure the VPN server
	// and the update the state to -> READY
	go configureGateway(obj.DB, resource)

	SetGatewayResponse(c, http.StatusCreated, *resource)
}

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
	if resource.Status != READY {
		log.Error("Impossibile to configure gateway. The current state is ", resource.Status)
		SetErrorResponse(c, http.StatusForbidden, ErrDeleteConfigurationGateway)
		return
	}
	// If it is configurable update the state to -> CONFIGURING
	resource.Status = DELETING_CONFIGURATION
	err := obj.DB.Save(&resource)
	if err.Error != nil {
		log.Error("Impossible to delete gateway configuration. Error saving in DB: ", error)
		SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
		return
	}
	// TODO creates go routine with httpclient to reset VPN server?
	go resetGateway(obj.DB, resource)

	// and the update the DB with nil param
	c.Status(http.StatusNoContent)
}
