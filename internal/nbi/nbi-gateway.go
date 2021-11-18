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

	log.Info("CreateGatewayConnectivity: called with param sliceID: " + sliceId)
	var json SliceInfo
	if err := c.ShouldBindJSON(&json); err != nil {
		log.Error("JSON body not correctly formatted")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add gatewayConnectivity obj info in the DB
	_, err := env.AddGatewayConnectivityInDB(sliceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		// otherwise the slise does not exist, so go ahead checking param
		log.Info("CreateGatewayConnectivity: creating a gateway network and subnet with CIDR: ", json.CIDR)

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
		log.Info("CreateGatewayConnectivity: updating GatewayConnectivity info in DB")
		env.UpdateGatewayConnectivityInDB(sliceId, network.ID, subnet.ID, router.ID, port.PortID)
	}
	fmt.Printf("%v", env.DB)
}

// GET gateway/connectivity
func (env *Env) RetrieveGatewayConnectivity(c *gin.Context) {
	sliceId := c.Query("sliceId")

	log.Info("RetrieveGatewayConnectivity: called with param sliceID: " + sliceId)
	fmt.Printf("%v", env.DB)

	_, gc, err := env.RetrieveGatewayConnectivityFromDB(sliceId)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Retrieve network and subnet object
	network, _ := env.Client.RetrieveNetworkById(gc.NetworkID)
	subnet, _ := env.Client.RetrieveSubnetById(gc.SubnetID)

	// Retrieve router and port object
	router, _ := env.Client.RetrieveRouterById(gc.RouterID)
	// TODO add also floating IP and VMID
	c.JSON(http.StatusOK, gin.H{"sliceId": sliceId, "network": network, "subnet": subnet, "router": router})
}

// DELETE gateway/connectivity
func (env *Env) DeleteGatewayConnectivity(c *gin.Context) {
	sliceId := c.Query("sliceId")

	log.Info("DeleteGatewayConnectivity: called with param sliceID: " + sliceId)

	// retrieve the slice associated objects (network, router) and delete them
	_, gc, err := env.RetrieveGatewayConnectivityFromDB(sliceId)
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
	// 2. Delete network and subnet
	err = env.Client.DeleteNetworkByID(gc.NetworkID)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// if everything is ok
	_, err = env.RemoveGatewayConnectivityFromDB(sliceId)
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

	log.Info("RetriveGatewayFloatingIP: called with params  sliceID: ", sliceId, " gatewayVmID", vmId)
	_, gatewayInfo, err := env.RetrieveGatewayConnectivityFromDB(sliceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	network, err := env.Client.RetrieveNetworkById(gatewayInfo.NetworkID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fip, err := env.Client.RetrieveGatewayFloatingIP(vmId, network.Name)
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
