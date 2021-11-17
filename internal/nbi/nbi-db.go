package nbi

import (
	// "fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// POST gateway/connectivity
func (env *Env) GetDB(c *gin.Context) {
	log.Info("GetDB called: retrieving all GatewayConnectivity instances")
	c.JSON(http.StatusOK, gin.H{"gateways": env.DB})
}

// POST gateway/connectivity
func (env *Env) AddDBEntry(c *gin.Context) {
	log.Info("AddDBEntry: called")
	var json GatewayConnectivity
	if err := c.ShouldBindJSON(&json); err != nil {
		log.Error("JSON body not well formatted")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sliceId := json.SliceID
	networkId := json.PrivNetID
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

// POST gateway/connectivity
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
