package nbi

import (
	"errors"
	"net"
	"net/http"
	NsmmApi "nextworks/nsm/api"
	"nextworks/nsm/internal/nsm"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

func checkPutGatewaysIdConfigurationBody(c *gin.Context) (*NsmmApi.PutGatewayConfigurationBody, uint16, uint16, error) {
	var jsonBody NsmmApi.PutGatewayConfigurationBody
	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		setErrorResponse(c, "PutGatewaysIdConfiguration", http.StatusBadRequest, nsm.ErrBodyMissingInfo)
		return nil, 0, 0, nsm.ErrBodyMissingInfo
	}
	// check ips
	if net.ParseIP(jsonBody.ExternalIp) == nil || net.ParseIP(jsonBody.ManagementIp) == nil {
		setErrorResponse(c, "PutGatewaysIdConfiguration", http.StatusBadRequest, nsm.ErrRequestConfigurationGateway)
		return nil, 0, 0, nsm.ErrRequestConfigurationGateway
	}

	// check ports
	mgmtPort, err := parsePort(jsonBody.ManagementPort)
	if err != nil {
		setErrorResponse(c, "PutGatewaysIdConfiguration", http.StatusBadRequest, nsm.ErrRequestConfigurationGateway)
		return nil, 0, 0, nsm.ErrRequestConfigurationGateway
	}

	vpnPort, err := parsePort(jsonBody.VpnServerPort)
	if err != nil {
		setErrorResponse(c, "PutGatewaysIdConfiguration", http.StatusBadRequest, nsm.ErrRequestConfigurationGateway)
		return nil, 0, 0, nsm.ErrRequestConfigurationGateway
	}
	return &jsonBody, mgmtPort, vpnPort, nil
}

// TODO to be implemented
// if no configuration exists, return error
// (GET /gateways/{id}/configuration)
func (obj *ServerInterfaceImpl) GetGatewaysIdConfiguration(c *gin.Context, id int) {
	var gc nsm.Gateway
	log.Info("GetGatewaysIdConfiguration - requested GET of gateway with ID: ", id)
	// Read gateway from DB
	result := obj.DB.First(&gc, id)

	// check error
	if result.Error != nil {
		log.Error("GetGatewaysIdConfiguration - error retrieving gateway from DB")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			setErrorResponse(c, "GetGatewaysIdConfiguration", http.StatusNotFound, nsm.ErrGatewayNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}
	// check if exists
	if gc.ExternalIp == "" || gc.ManagementIP == "" {
		setErrorResponse(c, "GetGatewaysIdConfiguration", http.StatusNotFound, nsm.ErrConfigurationGatewayNotExists)
		return
	}
	var response NsmmApi.ResponseGatewayConfigurationObject
	response.ExternalIp = gc.ExternalIp
	response.ManagementIp = gc.ManagementIP
	response.ManagementPort = strconv.Itoa(int(gc.ManagementPort))
	c.JSON(http.StatusOK, response)
}

// (PUT /gateways/{id}/configuration)
func (obj *ServerInterfaceImpl) PutGatewaysIdConfiguration(c *gin.Context, id int) {
	// Retrieve and check JSON body
	jsonBody, mnmtPort, vpnPort, err := checkPutGatewaysIdConfigurationBody(c)
	if err != nil {
		return
	}

	var gc nsm.Gateway
	log.Info("PutGatewaysIdConfiguration - requested configuration of gateway with ID: ", id)
	// Read gateway from DB
	result := obj.DB.First(&gc, id)

	if result.Error != nil {
		log.Error("PutGatewaysIdConfiguration - error retrieving gateway from DB")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			setErrorResponse(c, "PutGatewaysIdConfiguration", http.StatusNotFound, nsm.ErrGatewayNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}
	// check status
	if gc.Status != nsm.WAIT_FOR_GATEWAY && gc.Status != nsm.CONFIGURATION_ERROR {
		log.Info("PutGatewaysIdConfiguration - impossible to configure gateway configuration. The current state is ", gc.Status)
		setErrorResponse(c, "PutGatewaysIdConfiguration", http.StatusForbidden, nsm.ErrConfiguringGateway)
		return
	}

	gc.ExternalIp = jsonBody.ExternalIp
	gc.ManagementIP = jsonBody.ManagementIp
	gc.ManagementPort = mnmtPort
	gc.VPNServerInterface = jsonBody.VpnServerInterface
	gc.VPNServerPort = vpnPort
	gc.Status = nsm.CONFIGURING

	// Update database
	result = obj.DB.Save(&gc)
	if result.Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	go configureGateway(obj.DB, gc.ID)

	c.Status(http.StatusOK)
}

// (DELETE /gateways/{id}/configuration)
func (obj *ServerInterfaceImpl) DeleteGatewaysIdConfiguration(c *gin.Context, id int) {
	var gc nsm.Gateway
	log.Info("DeleteGatewaysIdConfiguration - requested removal of gateway with ID: ", id)
	// Read gateway from DB
	result := obj.DB.First(&gc, id)

	// check result
	if result.Error != nil {
		log.Error("DeleteGatewaysIdConfiguration - error retrieving gateway from DB")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			setErrorResponse(c, "DeleteGatewaysIdConfiguration", http.StatusNotFound, nsm.ErrGatewayNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}
	// check status
	if gc.Status != nsm.READY && gc.Status != nsm.CONFIGURATION_ERROR {
		log.Info("DeleteGatewaysIdConfiguration - impossible to delete gateway configuration. The current state is ", gc.Status)
		setErrorResponse(c, "DeleteGatewaysIdConfiguration", http.StatusForbidden, nsm.ErrDeleteConfigurationGateway)
		return
	}

	go resetGateway(obj.DB, gc.ID)

	// Remove configuration params from DB
	gc.ExternalIp = ""
	gc.ManagementIP = ""
	gc.ManagementPort = 0
	gc.VPNServerInterface = ""
	gc.VPNServerPort = 0
	gc.Status = nsm.DELETING_CONFIGURATION

	// Update database
	result = obj.DB.Save(&gc)
	if result.Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// TODO implement async go routine for configuration
func configureGateway(database *gorm.DB, id int) {
	var gc nsm.Gateway
	time.Sleep(time.Second * 5)
	database.First(&gc, id)
	gc.Status = nsm.READY
	log.Info("configureGateway ", gc.Status)
	database.Save(&gc)
}

func resetGateway(database *gorm.DB, id int) {
	var gc nsm.Gateway
	time.Sleep(time.Second * 5)
	database.First(&gc, id)
	gc.Status = nsm.WAIT_FOR_GATEWAY
	log.Info("resetGateway ", gc.Status)
	database.Save(&gc)
}
