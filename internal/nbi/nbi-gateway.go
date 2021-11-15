package nbi

import (
	// "fmt"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type SliceInfo struct {
	CIDR string `json:"cidr" binding:"required"` // CIDR of the subnet
}

// POST gateway/connectivity
func (env *Env) CreateGatewayConnectivity(c *gin.Context) {
	sliceId := c.Query("sliceId")

	log.Info("CreateGatewayConnectivity: sliceID: " + sliceId)
	var json SliceInfo
	if err := c.ShouldBindJSON(&json); err != nil {
		log.Error("JSON body not well formatted")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add slice info in the DB
	_, err := env.AddSliceConnectivity(sliceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		// otherwise the slise does not exist, so go ahead checking param
		log.Info("CreateGatewayConnectivity: creating a private network and subnet with CIDR: ", json.CIDR)

		// 1. Create a private net with a subnet
		networkName := sliceId + "_network"
		network, subnet, err := env.Client.CreateNetwork(networkName, json.CIDR)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 2. Crate a router connected with the floating net
		// with an interface on the private subnet
		routerName := sliceId + "_router"
		router, port, err := env.Client.CreateRouter(routerName, subnet.ID)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"network": network, "subnet": subnet, "router": router, "port": port})
		}

		// 3. Store these info locally in DB
		log.Info("CreateGatewayConnectivity: updating info in DB")
		env.UpdateSliceConnectivity(sliceId, network.ID, subnet.ID, router.ID, port.PortID)
	}
	fmt.Printf("%v", env.DB)
}

// GET gateway/connectivity
func (env *Env) RetrieveGatewayConnectivity(c *gin.Context) {
	sliceId := c.Query("sliceId")

	log.Info("RetrieveGatewayConnectivity: sliceID: " + sliceId)
	fmt.Printf("%v", env.DB)

	_, slice, err := env.RetrieveSliceConnectivity(sliceId)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"gatewayConnectivity": slice})
	}
}

// DELETE gateway/connectivity
func (env *Env) DeleteGatewayConnectivity(c *gin.Context) {
	sliceId := c.Query("sliceId")

	log.Info("DeleteGatewayConnectivity: sliceID: " + sliceId)

	// retrieve the slice associated objects (privnet, router) and delete them
	_, gc, err := env.RetrieveSliceConnectivity(sliceId)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Delete router and its interface
	if gc.RouterID == "" || gc.InterfaceID == "" {
		output := "gatewayConnectivity does not contain a RouterID or InterfaceID"
		log.Info("gatewayConnectivity does not contain a RouterID or InterfaceID : ", gc.RouterID, " ", gc.InterfaceID)
		c.JSON(http.StatusBadRequest, gin.H{"error": output})
		return
	}

	err = env.Client.DeleteRouter(gc.RouterID, gc.SubnetID)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 2. Delete privnet and subnet
	err = env.Client.DeleteNetworkByID(gc.PrivNetID)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// if everything is ok
	slice, err := env.RemoveSliceConnectivity(sliceId)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"slice": slice})
	}
}
