package nsm

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	nsmmapi "nextworks/nsm/api"
	gatewayconfig "nextworks/nsm/internal/gateway-config"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// (GET /net-resources/{{id}}/gateway/connections)
func (obj *ServerInterfaceImpl) GetNetResourcesIdGatewayConnections(c *gin.Context, id int) {

	// Retrieve the ResNet element with all the active connections
	res, error := RetrieveResourcesFromDB(obj.DB, id)
	if error != nil {
		log.Error("Impossible to retrieve network resources. Error reading from DB: ", error)
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, http.StatusNotFound, ErrResourcesNotExists)
			return
		} else {
			SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
			return
		}
	}
	// check status, it should be in running to have active connections
	if res.Status != RUNNING {
		log.Error("Impossibile to retrieve VPN connections. The current state is ", res.Status)
		SetErrorResponse(c, http.StatusForbidden, ErrGatewayNotRunning)
		return
	}

	err := LoadConnectionAssociationFromDB(obj.DB, res)
	if err != nil {
		SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
		return
	}
	log.Trace("Requested connection: ", res.Connections)

	// List all the connections found
	if len(res.Connections) == 0 {
		SetErrorResponse(c, http.StatusNotFound, ErrNoConnection)
	} else {
		var output []nsmmapi.Connection
		for _, conn := range res.Connections {
			apicon := nsmmapi.Connection{
				Id:                 conn.ID,
				PrivateKey:         &conn.PrivateKey,
				PublicKey:          &conn.PublicKey,
				RemotePeerIp:       conn.PeerIp,
				RemotePeerPort:     conn.PeerPort,
				PeerExposedSubnets: SubnetsToArray(conn.PeerNets),
			}
			output = append(output, apicon)
		}
		c.JSON(http.StatusOK, output)
	}

}

// (POST /net-resources/{{id}}/gateway/connections)
func (obj *ServerInterfaceImpl) PostNetResourcesIdGatewayConnections(c *gin.Context, id int) {
	var jsonBody nsmmapi.PostConnection

	// read from DB the requested resNet and check current status
	res, error := RetrieveResourcesFromDB(obj.DB, id)
	if error != nil {
		log.Error("Impossible to retrieve network-resources. Error reading from DB: ", error)
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, http.StatusNotFound, ErrResourcesNotExists)
			return
		} else {
			SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
			return
		}
	}
	// check status
	// !(!= READY && != RUNNING)
	if res.Status != READY && res.Status != RUNNING {
		log.Error("Impossibile to create a VPN connection. The current state is ", res.Status)
		SetErrorResponse(c, http.StatusForbidden, ErrGatewayNotConfigured)
		return
	}
	// read json body
	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		log.Error("Impossible to create a VPN connection. Error in the request, wrong json body")
		SetErrorResponse(c, http.StatusBadRequest, ErrBodyWrongInfo)
		return
	}

	out := net.ParseIP(jsonBody.RemotePeerIp)
	remoteSubnets, err := ParseExposedSubnets(jsonBody.PeerExposedSubnets)
	if out == nil || err != nil {
		SetErrorResponse(c, http.StatusBadRequest, ErrConnectionParameters)
		return
	}

	// TODO decide how to handle and manage pub/priv keys of the connection
	publicKey := ""
	privateKey := ""
	if jsonBody.PrivateKey != nil && jsonBody.PublicKey != nil {
		log.Debug("Private and Public keys passed as parameters")
		privateKey = *jsonBody.PrivateKey
		publicKey = *jsonBody.PublicKey
	}

	// build the VPNaaS service
	// var client gatewayconfig.VPNHttpClient
	client := gatewayconfig.New(net.ParseIP(res.Gateway.Config.MgmtIp), fmt.Sprint(res.Gateway.Config.MgmtPort), obj.VpnaasConfig.Environment)

	// call the Connect_to_VPN
	output := client.Connect(jsonBody.RemotePeerIp, fmt.Sprint(obj.VpnaasConfig.VpnaasPort), remoteSubnets, res.Gateway.Config.ExposedNets)
	log.Debug("VPNaaS connect output: ", output)
	if output {
		log.Trace("Creating a VPN connection object in DB")
		// create the state for the new VPN connection and save in BD
		conn := Connection{
			ResourceSetId: res.ID,
			PrivateKey:    privateKey,
			PublicKey:     publicKey,
			PeerIp:        jsonBody.RemotePeerIp,
			PeerPort:      fmt.Sprint(obj.VpnaasConfig.VpnaasPort),
			PeerNets:      SubnetsToString(jsonBody.PeerExposedSubnets),
		}

		result := obj.DB.Create(&conn)
		if result.Error != nil {
			log.Error("Error creation VPN connection in DB for resource set with ID ", res.ID, " and slice-id: ", res.SliceId)
			SetErrorResponse(c, http.StatusInternalServerError, ErrSavingConnectionDB)
			return
		}

		// update the state of netres if necessary
		if res.Status == READY {
			log.Trace("Updating resource-set state to RUNNING in DB...")
			res.Status = RUNNING
			result = obj.DB.Save(&res)
			if result.Error != nil {
				log.Error("Error updating resource-set for VPN connection creation, resource set with ID: ", res.ID, " and slice-id: ", res.SliceId)
				SetErrorResponse(c, http.StatusInternalServerError, ErrUpdatingGatewayInDB)
				return
			}
		}
		output := nsmmapi.Connection{
			Id:                 conn.ID,
			PrivateKey:         &privateKey,
			PublicKey:          &publicKey,
			RemotePeerIp:       conn.PeerIp,
			RemotePeerPort:     conn.PeerPort,
			PeerExposedSubnets: SubnetsToArray(conn.PeerNets),
		}
		c.JSON(http.StatusCreated, output)
	} else {
		log.Error("Error creating VPN connection with peer: ", jsonBody.RemotePeerIp)
		SetErrorResponse(c, http.StatusInternalServerError, ErrCreatingConnection)
		return
	}
}

// (DELETE /gateways/{id}/vpn/connections/{{cid}}/)
func (obj *ServerInterfaceImpl) DeleteNetResourcesIdGatewayConnectionsCid(c *gin.Context, id int, connectionid int) {
	// Delete the connection with ID connectionid and if no other connection are active and state is running switch to READY

	// Retrieve the ResNet element with all the active connections
	res, error := RetrieveResourcesFromDB(obj.DB, id)
	if error != nil {
		log.Error("Impossible to retrieve network-resources. Error reading from DB: ", error)
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, http.StatusNotFound, ErrResourcesNotExists)
			return
		} else {
			SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
			return
		}
	}
	// check status
	if res.Status != RUNNING {
		log.Error("Impossibile to delete a VPN connection. The current state is ", res.Status)
		SetErrorResponse(c, http.StatusForbidden, ErrGatewayNotRunning)
		return
	}

	// Search connection ID connectionid
	conn, err := SearchConnectionAssociationFromDB(obj.DB, res, connectionid)
	if err != nil {
		SetErrorResponse(c, http.StatusNotFound, ErrNoConnection)
		return
	}
	log.Trace("Requested connection: ", conn)

	// Create umu client to shutdown the connection
	// var client gatewayconfig.VPNHttpClient
	client := gatewayconfig.New(net.ParseIP(res.Gateway.Config.MgmtIp), fmt.Sprint(res.Gateway.Config.MgmtPort), obj.VpnaasConfig.Environment)

	// call the Connect_to_VPN
	output := client.Disconnect(conn.PeerIp, conn.PeerPort)
	log.Debug("VPNaaS disconnect output: ", output)
	if !output {
		log.Error("Error removing VPN connection with peer: ", conn.PeerIp)
		SetErrorResponse(c, http.StatusInternalServerError, ErrCreatingConnection)
		return
	}
	log.Trace("Closed VPN connection with peer: ", conn.PeerIp, " removing it from DB...")

	// remove from DB
	result := obj.DB.Delete(&conn)
	if result.Error != nil {
		log.Error("Error deleting connection, for resource set with ID: ", res.ID, " and slice-id: ", res.SliceId)
		SetErrorResponse(c, http.StatusInternalServerError, ErrSavingConnectionDB)
		return
	}
	//check if other connections exists
	// Search connection ID connectionid
	err = LoadConnectionAssociationFromDB(obj.DB, res)
	if err != nil {
		SetErrorResponse(c, http.StatusNotFound, ErrNoConnection)
		return
	}
	if len(res.Connections) == 0 {
		log.Debug("Remaining connections for this resource-set are 0, switch back to READY status")
		// Move back to READY STATE
		res.Status = READY
		result = obj.DB.Save(&res)
		if result.Error != nil {
			log.Error("Error updating resource-set VPN for connection deletion, resource set with ID: ", res.ID, " and slice-id: ", res.SliceId)
			SetErrorResponse(c, http.StatusInternalServerError, ErrUpdatingGatewayInDB)
			return
		}
	}
	c.Status(http.StatusOK)
}

// (GET /net-resources/{{id}}/gateway/connections/{{cid}})
func (obj *ServerInterfaceImpl) GetNetResourcesIdGatewayConnectionsCid(c *gin.Context, id int, connectionid int) {

	// Retrieve the ResNet element with all the active connections
	res, error := RetrieveResourcesFromDB(obj.DB, id)
	if error != nil {
		log.Error("Impossible to retrieve network-resources. Error reading from DB: ", error)
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, http.StatusNotFound, ErrResourcesNotExists)
			return
		} else {
			SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
			return
		}
	}

	// check status, it should be in running to have active connections
	if res.Status != RUNNING {
		log.Error("Impossibile to retrieve VPN connections. The current state is ", res.Status)
		SetErrorResponse(c, http.StatusForbidden, ErrGatewayNotRunning)
		return
	}

	conn, err := SearchConnectionAssociationFromDB(obj.DB, res, connectionid)
	if err != nil {
		SetErrorResponse(c, http.StatusNotFound, ErrNoConnection)
		return
	}
	log.Trace("Requested connection: ", conn)

	output := nsmmapi.Connection{
		Id:                 conn.ID,
		PrivateKey:         &conn.PrivateKey,
		PublicKey:          &conn.PublicKey,
		RemotePeerIp:       conn.PeerIp,
		RemotePeerPort:     conn.PeerPort,
		PeerExposedSubnets: SubnetsToArray(conn.PeerNets),
	}

	c.JSON(http.StatusOK, output)
}
