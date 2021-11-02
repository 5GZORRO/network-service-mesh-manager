package main

import (
	// "fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

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
func createNetwork(c *gin.Context) {
	networkName := c.Query("name")

	var json Network
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check params
	log.Println("createNetwork: " + networkName + " " + json.Name)
	if networkName != json.Name {
		log.Error("Wrong parameters: network name in the URL is different from network name in JSON body ")
		c.Status(http.StatusBadRequest)
		return
	}
	// logic
	network, subnet, err := CreateNetwork(networkName, json.CIDR)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"network": network, "subnet": subnet})
	}
}

// /network?name=name
func deleteNetwork(c *gin.Context) {
	networkName := c.Query("name")

	log.Info("deleteNetwork:" + networkName)
	err := DeleteNetwork(networkName)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.Status(http.StatusOK)
	}
}

// /network?name=name
func retrieveNetwork(c *gin.Context) {
	networkName := c.Query("name")
	log.Info("retrieveNetwork:" + networkName)
	// logic
	network, err := RetrieveNetwork(networkName)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		log.Info(*network)
		c.JSON(http.StatusOK, network)
	}
}
