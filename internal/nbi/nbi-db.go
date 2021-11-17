package nbi

import (
	// "fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// POST gateway/connectivity
func (env *Env) GetDB(c *gin.Context) {
	log.Info("Retrieving all slice istances")
	c.JSON(http.StatusOK, gin.H{"gateways": env.DB})
}

// POST gateway/connectivity
func (env *Env) AddDBEntry(c *gin.Context) {
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
	log.Info("Adding new slice istance with ", json)
	log.Info(sliceId, networkId, subnetId, routerId)

	index, _, err := env.RetrieveGatewayConnectivityFromDB(sliceId)
	if err != nil && index == -1 {
		env.DB = append(env.DB, json)
		log.Info("Inserted a new slice with name ", sliceId, " [DB len: ", len(env.DB), " capacity: ", cap(env.DB), "] \nDB:", env.DB)
		c.JSON(http.StatusOK, gin.H{"gateway": json})
	}
	c.Status(http.StatusBadRequest)
}

// POST gateway/connectivity
func (env *Env) DeleteDBEntry(c *gin.Context) {
	sliceId := c.Query("sliceId")
	log.Info("Deleting slice instance with sliceID: ", sliceId)
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
