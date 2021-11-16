// Package nbi provides struct and methods used in the server routes
// It defines an Env struct which stores the global status for the calls
// (such as DB connection pool or OpenStack sessions)
// Each Gin server routes corresponds to a Env method
package nbi

import (
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// TODO 1 add a sliceId param in all the network request
// TODO 2 resources associated to a slice should be saved somewhere

type Network struct {
	Name string `json:"name" binding:"required"`
	CIDR string `json:"cidr" binding:"required"`
}

// /network?name=name
// Input Body
// {
//     "name": "name",
//     "cidr": "cidr"
// }
// Query name parameters should be equal to name in the JSON body
// This method first creates the Network, then it creates the Subnet associated to it
func (env *Env) CreateNetwork(c *gin.Context) {
	networkName := c.Query("name")

	var json Network
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check params
	log.Println("createNetwork: creating a network with name: " + json.Name + " and CIDR: " + json.CIDR)
	if networkName != json.Name {
		log.Error("Wrong parameters: network name in the URL is different from network name in JSON body ")
		c.Status(http.StatusBadRequest)
		return
	}
	// logic
	network, subnet, err := env.Client.CreateNetwork(networkName, json.CIDR)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"network": network, "subnet": subnet})
	}
}

// /network?name=name
func (env *Env) DeleteNetwork(c *gin.Context) {
	networkName := c.Query("name")

	log.Info("deleteNetwork:" + networkName)
	err := env.Client.DeleteNetworkByName(networkName)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.Status(http.StatusOK)
	}
}

// /network?name=name
func (env *Env) RetrieveNetwork(c *gin.Context) {
	networkName := c.Query("name")
	log.Info("retrieveNetwork:" + networkName)
	// logic
	network, err := env.Client.RetrieveNetworkByName(networkName)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		log.Info(*network)
		c.JSON(http.StatusOK, network)
	}
}
