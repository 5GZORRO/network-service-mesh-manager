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
	}
	// Retrieve network and subnet object
	network, _ := env.Client.RetrieveNetworkById(slice.PrivNetID)
	subnet, _ := env.Client.RetrieveSubnetById(slice.SubnetID)

	// Retrieve router and port object
	router, _ := env.Client.RetrieveRouterById(slice.RouterID)
	c.JSON(http.StatusOK, gin.H{"sliceId": sliceId, "network": network, "subnet": subnet, "router": router})
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
	_, err = env.RemoveSliceConnectivity(sliceId)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.Status(http.StatusOK)
	}
}

// gateway/connectivity/ip?sliceId=test&vmId=af9e6bcd-356e-46f6-b625-1dbba06e21dc
// function retrieve FIP and store it in the DB with the VM ID
func (env *Env) RetriveGatewayFloatingIP(c *gin.Context) {
	sliceId := c.Query("sliceId")
	vmId := c.Query("vmId")

	_, gatewayInfo, err := env.RetrieveSliceConnectivity(sliceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	network, err := env.Client.RetrieveNetworkById(gatewayInfo.PrivNetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fip, err := env.Client.RetriveGatewayFloatingIP(vmId, network.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Info("Found floatingIP ", fip, " for sliceID ", sliceId)
	gatewayInfo.FloatingIP = fip
	gatewayInfo.VmGatewayID = vmId
	c.JSON(http.StatusOK, gin.H{"gatewayInfo": gatewayInfo})
}

// TODO delete gatewayIP to be implemented?
