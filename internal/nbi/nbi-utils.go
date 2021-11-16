package nbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type GatewayInfo struct {
	NetworkID string `json:"networkID" binding:"required"`
	SubnetID  string `json:"subnetID" binding:"required"`
	RouterID  string `json:"routerID" binding:"required"`
}

// clean up procedure
func (env *Env) CleanUp(c *gin.Context) {
	for i := range env.DB {
		obj := env.DB[i]
		log.Info("Clean ", obj)
		// try to remove port and router router
		env.Client.DeleteRouter(obj.RouterID, obj.SubnetID)
		// try to remove network
		env.Client.DeleteNetworkByID(obj.PrivNetID)
	}
	// TODO reset DB
	c.Status(http.StatusOK)
}

// clean up procedure
func (env *Env) CleanUpGateway(c *gin.Context) {
	var json GatewayInfo
	if err := c.ShouldBindJSON(&json); err != nil {
		log.Error("JSON body not well formatted")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	networkId := json.NetworkID
	subnetId := json.SubnetID
	routerId := json.RouterID
	log.Info("Clean procedure...")

	if routerId != "" {
		log.Info("Removing router with ID ", routerId)
		env.Client.DeleteRouter(routerId, subnetId)
	}
	if networkId != "" {
		log.Info("Removing network with ID ", networkId, " and subnet with ID ", subnetId)
		env.Client.DeleteNetworkByID(networkId)
	}
	c.Status(http.StatusOK)
}

// /test endpoint
func (env *Env) Test(c *gin.Context) {
	log.Info("Test:" + env.Client.TenantID)
}
