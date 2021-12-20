package nsm

import (
	"errors"
	"net"
	"net/http"
	nsmmapi "nextworks/nsm/api"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func checkGatewayConfigurationParams(input nsmmapi.Gateway) error {
	_, err := parsePort(input.MgmtPort)
	if err != nil {
		return ErrGatewayConfigMgmtPort
	}
	if _, _, err = net.ParseCIDR(input.SubnetToExpose); err != nil {
		return ErrGatewayConfigSubnet
	}
	mngmIp := net.ParseIP(input.MgmtIp)
	externalIP := net.ParseIP(input.ExternalIp)
	if mngmIp == nil {
		return ErrGatewayConfigMgmtIp
	}
	if externalIP == nil {
		return ErrGatewayConfigExternalIp
	}
	return nil
}

func (obj *ServerInterfaceImpl) GetNetResourcesIdGateway(c *gin.Context, id int) {
	// Retrive the Resource Set state and check it, if the gateway can be configured

	log.Trace("GetNetResourcesIdGateway - requested retrieve of gateway configuration for resource set with ID: ", id)
	resource, error := RetrieveResourcesFromDB(obj.DB, id)
	if error != nil {
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, "GetNetResourcesIdGateway", http.StatusNotFound, ErrSliceNotExists)
			return
		} else {
			SetErrorResponse(c, "GetNetResourcesIdGateway", http.StatusInternalServerError, ErrGeneral)
			return
		}
	}

	// check if it is configured
	if resource.Gateway.ExternalIp == "" {
		SetErrorResponse(c, "GetNetResourcesIdGateway", http.StatusNotFound, ErrGatewayNotConfigured)
		return
	}
	SetGatewayResponse(c, http.StatusOK, *resource)
}

func (obj *ServerInterfaceImpl) PutNetResourcesIdGateway(c *gin.Context, id int) {
	var jsonBody nsmmapi.Gateway

	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		SetErrorResponse(c, "PutNetResourcesIdGateway", http.StatusBadRequest, ErrRequestConfigurationGateway)
		return
	}

	log.Trace("PutNetResourcesIdGateway - requested configuration of gateway for resource set with ID: ", id)
	resource, error := RetrieveResourcesFromDB(obj.DB, id)
	if error != nil {
		log.Error(error)
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, "PutNetResourcesIdGateway", http.StatusNotFound, ErrResourcesNotExists)
			return
		} else {
			SetErrorResponse(c, "PutNetResourcesIdGateway", http.StatusInternalServerError, ErrGeneral)
			return
		}
	}
	// check status
	if resource.Status != WAIT_FOR_GATEWAY_CONFIG {
		log.Error("PutNetResourcesIdGateway - impossibile to configure gateway. The current state is ", resource.Status)
		SetErrorResponse(c, "PutNetResourcesIdGateway", http.StatusForbidden, ErrConfiguringGateway)
		return
	}

	if err := checkGatewayConfigurationParams(jsonBody); err != nil {
		SetErrorResponse(c, "PutNetResourcesIdGateway", http.StatusBadRequest, err)
		return
	}

	resource.Gateway.ExternalIp = jsonBody.ExternalIp
	resource.Gateway.MgmtIp = jsonBody.MgmtIp
	resource.Gateway.MgmtPort, _ = parsePort(jsonBody.MgmtPort)
	resource.Gateway.ExposedNets = jsonBody.SubnetToExpose
	resource.Gateway.PubKey = jsonBody.PubKey
	// If it is configurable update the state to -> CONFIGURING
	// and store params
	resource.Status = CONFIGURING
	log.Trace("PutNetResourcesIdGateway - updating gateway configuration for resource set with ID: ", id)
	err := obj.DB.Save(&resource)
	if err.Error != nil {
		log.Error("PutNetResourcesIdGateway - error updating resource set with gateway configuration ", err.Error)
		SetErrorResponse(c, "PutNetResourcesIdGateway", http.StatusInternalServerError, ErrGeneral)
		return
	}
	// go routine with httpclient to configure the VPN server
	// and the update the state to -> READY
	go configureGateway(obj.DB, resource)

	SetGatewayResponse(c, http.StatusCreated, *resource)
}

func (obj *ServerInterfaceImpl) DeleteNetResourcesIdGateway(c *gin.Context, id int) {
	log.Trace("DeleteNetResourcesIdGateway - requested removal of gateway configuration for resource set with ID: ", id)
	resource, error := RetrieveResourcesFromDB(obj.DB, id)
	if error != nil {
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, "DeleteNetResourcesIdGateway", http.StatusNotFound, ErrResourcesNotExists)
			return
		} else {
			SetErrorResponse(c, "DeleteNetResourcesIdGateway", http.StatusInternalServerError, ErrGeneral)
			return
		}
	}
	// check status
	if resource.Status != READY {
		log.Error("DeleteNetResourcesIdGateway - impossibile to configure gateway. The current state is ", resource.Status)
		SetErrorResponse(c, "DeleteNetResourcesIdGateway", http.StatusForbidden, ErrDeleteConfigurationGateway)
		return
	}
	// If it is configurable update the state to -> CONFIGURING
	resource.Status = DELETING_CONFIGURATION
	err := obj.DB.Save(&resource)
	if err.Error != nil {
		SetErrorResponse(c, "DeleteNetResourcesIdGateway", http.StatusInternalServerError, ErrGeneral)
		return
	}
	// TODO creates go routine with httpclient to reset VPN server?
	go resetGateway(obj.DB, resource)

	// and the update the DB with nil param
	c.Status(http.StatusNoContent)
}

// TODO configureGateway is a goroutine to configure the VM gateway, using
// an HTTP client
func configureGateway(database *gorm.DB, res *ResourceSet) {
	time.Sleep(time.Second * 10)
	// TODO configure VM gateway

	// update the state
	res.Status = READY
	result := database.Save(&res)
	if result.Error != nil {
		log.Error("Error updating resource set status with ID: ", res.ID, " and slice-id: ", res.SliceId)
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
