package nbi

import (
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// TODO 2 resources associated to a slice should be saved somewhere

// gateway/connectivity
func (env *Env) CreateGatewayConnectivity(c *gin.Context) {
	sliceId := c.Query("sliceId")

	log.Info("CreateGatewayConnectivity: sliceID: " + sliceId)
	var json Network
	if err := c.ShouldBindJSON(&json); err != nil {
		log.Error("JSON body not well formatted")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check params
	log.Println("CreateGatewayConnectivity: creating a private network with name: " + json.Name + " and CIDR: " + json.CIDR)

	// TODO create a priv net and a routed associated to this sliceId
	// and store the infos locally

	// 1. create privnet
	// network, subnet, err := env.Client.CreateNetwork(networkName, json.CIDR)
	// if err != nil {
	// 	log.Error(err)
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// } else {
	// 	c.JSON(http.StatusOK, gin.H{"network": network, "subnet": subnet})
	// }
}

// gateway/connectivity
func (env *Env) RetrieveGatewayConnectivity(c *gin.Context) {
	sliceId := c.Query("sliceId")

	log.Info("RetrieveGatewayConnectivity: sliceID: " + sliceId)

	// TODO retrieve objects (privnet, router) associated to the sliceID
}

// gateway/connectivity
func (env *Env) DeleteGatewayConnectivity(c *gin.Context) {
	sliceId := c.Query("sliceId")

	log.Info("DeleteGatewayConnectivity: sliceID: " + sliceId)
	// TODO retrieve the slice associated objects (privnet, router)
	// and delete them
}
