package nsm

import (
	"errors"
	"net/http"
	nsmmapi "nextworks/nsm/api"
	"strconv"
	"strings"

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

	// check if it has an external ip allocated
	if resource.Gateway.External.ExternalIp == "" {
		log.Error("Impossible to retrieve gateway configuration. It is not yet configured")
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

	//  SET CONFIGURING STATE AND SAVE
	resource.Status = CONFIGURING
	log.Trace("Start association of a floating IPs - updating state to CONFIGURINF of resource set with ID: ", id)
	output := obj.DB.Save(&resource)
	if output.Error != nil {
		log.Error("Impossible to update resource set - error saving in DB ", output.Error)
		SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
		return
	}
	// Floating IP should be already created

	config := Config{}
	// TODO this should be different, API not exposed on external IP
	config.MgmtIp = resource.Gateway.External.ExternalIp
	config.MgmtPort, _ = parsePort(strconv.Itoa(int(obj.VpnaasConfig.VpnaasPort)))
	resource.Gateway.Config = config

	// resource.Gateway.ExposedNets = SubnetsToString(jsonBody.SubnetToExpose)
	err := LoadNetworkAssociationFromDB(obj.DB, resource)
	if err != nil {
		resource.Status = CONFIGURATION_ERROR
		result := obj.DB.Save(&resource)
		if result.Error != nil {
			log.Error("Error updating resource set status with ID: ", resource.ID, " and slice-id: ", resource.SliceId)
		}
		log.Error("Error retrieving network associations of resource set with ID: ", resource.ID)
		SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
		return
	}
	var exposedNetworks []string
	for _, network := range resource.Networks {
		log.Trace("Network: with name: ", network.NetworkName)
		if strings.HasPrefix(network.NetworkName, obj.Netconfig.ExposedNetworksNamePrefix) {
			log.Info("Selected Network with name ", network.NetworkName)
			exposedNetworks = append(exposedNetworks, network.SubnetCidr)
		}

	}
	log.Trace("ExposedNetworks selected: ", exposedNetworks)
	resource.Gateway.Config.ExposedNets = SubnetsToString(exposedNetworks)
	log.Trace("ExposedNetworks stored: ", resource.Gateway.Config.ExposedNets)

	// NO mode check in the PrivateVPNRange to pass to VPNaaS
	resource.Gateway.Config.PrivateVpnRange = obj.Netconfig.PrivateVpnRange
	log.Info("Setting private VPN Range as ", resource.Gateway.Config.PrivateVpnRange)

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
	if resource.Status != READY && resource.Status != CONFIGURATION_ERROR {
		log.Error("Impossibile to delete gateway configuration. The current state is ", resource.Status)
		SetErrorResponse(c, http.StatusForbidden, ErrDeleteConfigurationGateway)
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
	go resetGateway(obj.DB, resource)

	c.Status(http.StatusNoContent)
}
