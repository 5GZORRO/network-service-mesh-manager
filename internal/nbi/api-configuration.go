package nbi

import (
	"errors"
	"net"
	"net/http"
	NsmmApi "nextworks/nsm/api"
	"nextworks/nsm/internal/nsm"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

// TODO different object?
// (GET /gateways/{id}/configuration)
func (obj *ServerInterfaceImpl) GetGatewaysIdConfiguration(c *gin.Context, id int) {}

func checkPutGatewaysIdConfigurationBody(c *gin.Context) (*NsmmApi.PutGatewayConfigurationBody, error) {
	var jsonBody NsmmApi.PutGatewayConfigurationBody
	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		setErrorResponse(c, "PutGatewaysIdConfiguration", http.StatusBadRequest, nsm.ErrBodyMissingInfo)
		return nil, nsm.ErrBodyMissingInfo
	}
	// check
	if net.ParseIP(jsonBody.ExternalIp) == nil || net.ParseIP(jsonBody.ManagementIp) == nil {
		setErrorResponse(c, "PutGatewaysIdConfiguration", http.StatusBadRequest, nsm.ErrGatewayConfiguration)
		return nil, nsm.ErrGatewayConfiguration
	}

	// TODO check port
	// port, err = parsePort(jsonBody.ManagementPort)
	// if err != nil {
	// 	setErrorResponse(c, "PutGatewaysIdConfiguration", http.StatusBadRequest, nsm.ErrBodyWrongInfo)
	// }
	return &jsonBody, nil
}

// (PUT /gateways/{id}/configuration)
func (obj *ServerInterfaceImpl) PutGatewaysIdConfiguration(c *gin.Context, id int) {
	// Retrieve and check JSON body
	jsonBody, err := checkPutGatewaysIdConfigurationBody(c)
	if err != nil {
		return
	}

	var gc nsm.Gateway
	log.Info("PutGatewaysIdConfiguration - requested configuration of gateway with ID: ", id)
	// Read gateway from DB
	result := obj.DB.First(&gc, id)

	if result.Error != nil {
		log.Error("PutGatewaysIdConfiguration - error retrieving gateway")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			setErrorResponse(c, "PutGatewaysIdConfiguration", http.StatusNotFound, nsm.ErrGatewayNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}
	// check status
	if gc.Status != nsm.WAIT_FOR_GATEWAY {
		log.Info("PutGatewaysIdConfiguration - impossible to configure gateway. The current state is ", gc.Status)
		setErrorResponse(c, "PutGatewaysIdConfiguration", http.StatusForbidden, nsm.ErrConfiguringGateway)
		return
	}

	// TODO configure
	gc.ExternalIp = jsonBody.ExternalIp
	gc.ManagementIP = jsonBody.ManagementIp
	// TODO fix port
	// gc.ManagementPort = jsonBody.ManagementPort

	// Update database
	result = obj.DB.Save(&gc)
	if result.Error != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

// (DELETE /gateways/{id}/configuration)
func (obj *ServerInterfaceImpl) DeleteGatewaysIdConfiguration(c *gin.Context, id int) {}
