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

// GET /management/db
func (env *Env) GetDB(c *gin.Context) {
	log.Info("GetDB called: retrieving all GatewayConnectivity instances")
	c.JSON(http.StatusOK, gin.H{"gateways": env.DB})
}

// POST /management/db
func (env *Env) AddDBEntry(c *gin.Context) {
	log.Info("AddDBEntry: called")
	var json GatewayConnectivity
	if err := c.ShouldBindJSON(&json); err != nil {
		log.Error("JSON body not well formatted")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sliceId := json.SliceID
	networkId := json.NetworkID
	subnetId := json.SubnetID
	routerId := json.RouterID
	log.Info("Adding new GatewayConnectivity istance with infos ", json)
	log.Info(sliceId, networkId, subnetId, routerId)

	index, _, err := env.RetrieveGatewayConnectivityFromDB(sliceId)
	if err != nil && index == -1 {
		env.DB = append(env.DB, json)
		log.Info("Added a new GatewayConnectivty for slice with sliceID: ", sliceId, " \n[DB len: ", len(env.DB), " capacity: ", cap(env.DB), "] \nDB:", env.DB)
		c.JSON(http.StatusOK, gin.H{"gateway": json})
		return
	}
	c.Status(http.StatusBadRequest)
}

// DELETE /management/db?sliceId=test
func (env *Env) DeleteDBEntry(c *gin.Context) {
	sliceId := c.Query("sliceId")
	log.Info("DeleteDBEntry: called with param sliceID: " + sliceId)
	if sliceId != "" {
		_, _, err := env.RetrieveGatewayConnectivityFromDB(sliceId)
		if err != nil {
			c.Status(http.StatusNotFound)
		}
		env.RemoveGatewayConnectivityFromDB(sliceId)
		c.Status(http.StatusOK)
	}
	c.Status(http.StatusBadRequest)
}

//
// Managementc clean functions
//

// GET /management/clean clean up procedure
func (env *Env) CleanUp(c *gin.Context) {
	for i := range env.DB {
		obj := env.DB[i]
		log.Info("Clean ", obj)
		// try to remove port and router router
		env.Client.DeleteRouter(obj.RouterID, obj.SubnetID)
		// try to remove network
		env.Client.DeleteNetworkByID(obj.NetworkID)
	}
	// TODO reset DB
	c.Status(http.StatusOK)
}

// POST /management/clean clean up procedure with
// gatewayInfo as body. It deletes the resources passed as params
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

	if routerId != "" && subnetId != "" {
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
