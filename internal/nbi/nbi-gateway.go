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

// gateway/connectivity
func (env *Env) CreateGatewayConnectivity(c *gin.Context) {
	sliceId := c.Query("sliceId")

	log.Info("CreateGatewayConnectivity: sliceID: " + sliceId)
	var json SliceInfo
	if err := c.ShouldBindJSON(&json); err != nil {
		log.Error("JSON body not well formatted")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//
	_, err := env.AddSliceConnectivity(sliceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		// otherwise, the slide does not exist, so go ahead checking param
		log.Println("CreateGatewayConnectivity: creating a private network and subnet with CIDR: " + json.CIDR)

		// TODO
		// 1. Create a private net with a subnet
		// networkName := sliceId + "_network"
		// network, subnet, err := env.Client.CreateNetwork(networkName, json.CIDR)
		// if err != nil {
		// 	log.Error(err)
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// }
		// network and sunet stores OS object
		// c.JSON(http.StatusOK, gin.H{"network": network, "subnet": subnet})

		// 2. Crate a router connected with the floating net
		routerName := sliceId + "_router"
		router, err := env.Client.CreateRouter(routerName)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"router": router})
		}

		// 3. Add interface to the router on the private net
		// 4. Store these info locally in DB

		// If all steps are successful, update infos in the slice
		// or delete it in case of error
		env.UpdateSliceConnectivity(sliceId, "", "", router.ID)
	}
	fmt.Printf("%v", env.DB)
}

// TODO
// gateway/connectivity
func (env *Env) RetrieveGatewayConnectivity(c *gin.Context) {
	sliceId := c.Query("sliceId")

	log.Info("RetrieveGatewayConnectivity: sliceID: " + sliceId)
	fmt.Printf("%v", env.DB)

	// TODO retrieve objects (privnet, router) associated to the sliceID
	_, slice, err := env.RetrieveSliceConnectivity(sliceId)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"gatewayConnectivity": slice})
	}
}

// gateway/connectivity
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
	// TODO delete all the associated resources

	// 1. Delete router
	if gc.RouterID == "" {
		output := "gatewayConnectivity does not contain a RouterID"
		log.Error(output)
		c.JSON(http.StatusBadRequest, gin.H{"error": output})
	}
	err = env.Client.DeleteRouter(gc.RouterID)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 2. Delete privnet and subnet
	// TODO

	// if everything is ok
	slice, err := env.RemoveSliceConnectivity(sliceId)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"slice": slice})
	}
}
