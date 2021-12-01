// Package to implement all of the handlers in the ServerInterface
package nbi

import (
	"errors"
	"net/http"
	NsmmApi "nextworks/nsm/api"
	"nextworks/nsm/internal/nsm"
	"nextworks/nsm/internal/vim"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setErrorResponse(ctx *gin.Context, method string, errorStatus int, err error) {
	log.Error(method, " - ", err.Error())
	outputJson := NsmmApi.ErrorResponse{Error: err.Error()}
	ctx.JSON(errorStatus, outputJson)
}

func buildGatewayResponse(gc nsm.Gateway) NsmmApi.ResponseGateway {
	return NsmmApi.ResponseGateway{
		Id:      gc.ID,
		SliceId: gc.SliceID,
		Status:  gc.Status,
		VimName: gc.VimName,
		Subnet:  gc.Resources.SubnetCidr,
	}
}

func setGatewayResponse(ctx *gin.Context, status int, gc nsm.Gateway) {
	output := buildGatewayResponse(gc)
	ctx.JSON(status, output)
}

// (GET /gateways)

func (obj *ServerInterfaceImpl) GetGateways(ctx *gin.Context, params NsmmApi.GetGatewaysParams) {
	log.Info("GetGateways")
	var gco nsm.Gateway
	var as nsm.OpenstackResource
	var gcos []nsm.Gateway
	var result *gorm.DB

	if params.SliceId != nil {
		log.Info("GetGateways - filtered by sliceId ", *params.SliceId)
		result = obj.DB.First(&gco, "slice_id = ?", params.SliceId)

		if result.RowsAffected == 0 {
			log.Error("GetGateways - 0 rows read")
			setErrorResponse(ctx, "GetGateways", http.StatusNotFound, nsm.ErrGatewayNotFound)
			return
		}
		if result.Error != nil {
			log.Error("GetGateways - Error reading from DB")
			ctx.Status(http.StatusInternalServerError)
			return
		}
		// Read association
		err := obj.DB.Model(&gco).Association("Resources").Find(&as)
		if err != nil {
			log.Error(result.Error)
			setErrorResponse(ctx, "PostGateways", http.StatusInternalServerError, nsm.ErrGeneral)
			return
		}
		gco.Resources = as
		log.Trace(gco)
		setGatewayResponse(ctx, http.StatusOK, gco)
	} else {
		log.Info("GetGateways - without query parameters")
		result = obj.DB.Find(&gcos)

		if result.RowsAffected == 0 {
			log.Error("GetGateways - 0 rows read")
			setErrorResponse(ctx, "GetGateways", http.StatusNotFound, nsm.ErrGatewayNotFound)
			return
		}
		if result.Error != nil {
			log.Error("GetGateways - Error reading from DB")
			ctx.Status(http.StatusInternalServerError)
			return
		}
		var response NsmmApi.ResponseGatewaysList
		for _, s := range gcos {
			// read resources
			err := obj.DB.Model(&s).Association("Resources").Find(&as)
			if err != nil {
				log.Error(result.Error)
				setErrorResponse(ctx, "PostGateways", http.StatusInternalServerError, nsm.ErrGeneral)
				return
			}
			s.Resources = as
			log.Trace(s)
			response = append(response, buildGatewayResponse(s))
		}
		ctx.JSON(http.StatusOK, response)
	}
}

func (obj *ServerInterfaceImpl) PostGateways(ctx *gin.Context) {
	// Retrieve body
	var jsonBody NsmmApi.PostGateway
	if err := ctx.ShouldBindJSON(&jsonBody); err != nil {
		log.Error(nsm.ErrBodyMissingInfo.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": nsm.ErrBodyMissingInfo.Error()})
		return
	}
	// TODO check body params and format ex. subnet and in case return 400
	// TODO check if subnet is already in use

	log.Debug("PostGateways - requesting creation of Gateway connection for SliceId: ", jsonBody.SliceId, " on VIM '", jsonBody.VimName, "' with subnet ", jsonBody.Subnet)
	// Check if a Vim with this name exists
	if !obj.Vims.Exists(jsonBody.VimName) {
		setErrorResponse(ctx, "PostGateways", http.StatusForbidden, vim.ErrVimNotFound)
		return
	}

	// Check if a gateway for the slice already exists
	var gco nsm.Gateway
	result := obj.DB.First(&gco, "slice_id = ?", jsonBody.SliceId)

	if result.Error == nil {
		setErrorResponse(ctx, "PostGateways", http.StatusForbidden, nsm.ErrGatewayExists)
		return
	} else {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			setErrorResponse(ctx, "PostGateways", http.StatusInternalServerError, nsm.ErrGeneral)
			return
		}
	}
	// gateway does not exists, create a new one
	log.Trace("PostGateways - gateway for requested slice not found")
	gc := nsm.Gateway{
		SliceID: jsonBody.SliceId,
		Status:  nsm.CREATING,
		VimName: jsonBody.VimName}
	gc.Resources = nsm.OpenstackResource{
		SubnetCidr: jsonBody.Subnet}
	log.Trace(gc)
	result = obj.DB.Create(&gc)

	if result.RowsAffected != 1 {
		log.Error("Error writing in DB")
		setErrorResponse(ctx, "PostGateways", http.StatusInternalServerError, nsm.ErrGeneral)
		return
	}
	log.Info("PostGateways - gateway for requested slice initialized")
	// TODO create il gatewayconnectivity su Openstack and update all info
	// drivers.VimDriver.CreateGatewayConnectivity()
	createResources(obj.DB, &gc)

	setGatewayResponse(ctx, http.StatusOK, gc)
}

func (obj *ServerInterfaceImpl) GetGatewaysId(ctx *gin.Context, id int) {
	var gc nsm.Gateway
	var as nsm.OpenstackResource
	log.Debug("GetGatewaysId - requested Gateway info with ID: ", id)
	// Read from DB, check if exists
	result := obj.DB.First(&gc, id)
	log.Trace(gc)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			setErrorResponse(ctx, "GetGatewaysId", http.StatusNotFound, nsm.ErrGatewayNotFound)
			return
		}
		log.Error(result.Error)
		setErrorResponse(ctx, "PostGateways", http.StatusInternalServerError, nsm.ErrGeneral)
		return
	} else {
		// Read association
		err := obj.DB.Model(&gc).Association("Resources").Find(&as)

		if err != nil {
			log.Error(err)
			setErrorResponse(ctx, "PostGateways", http.StatusInternalServerError, nsm.ErrGeneral)
			return
		}
		log.Trace(as)
		gc.Resources = as
		setGatewayResponse(ctx, http.StatusOK, gc)
		return
	}
}

func (obj *ServerInterfaceImpl) DeleteGatewaysId(ctx *gin.Context, id int) {
	var gc nsm.Gateway
	var as nsm.OpenstackResource

	log.Info("DeleteGatewaysId - requesting removal of Gateway connection with ID: ", id)
	// Read from DB
	result := obj.DB.First(&gc, id)

	if result.Error != nil {
		log.Error("DeleteGatewaysId - error retrieving gateway with ID ", id)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			setErrorResponse(ctx, "DeleteGatewaysId", http.StatusNotFound, nsm.ErrGatewayNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}
	// check status
	if gc.Status != nsm.WAIT_FOR_GATEWAY_CONFIG && gc.Status != nsm.CREATION_ERROR {
		log.Info("DeleteGatewaysId - impossible to delete gateway. The current state is ", gc.Status)
		setErrorResponse(ctx, "DeleteGatewaysId", http.StatusForbidden, nsm.ErrGatewayCantBeDeleted)
		return
	}
	// Read association
	err := obj.DB.Model(&gc).Association("Resources").Find(&as)

	if err != nil {
		log.Error(err)
		setErrorResponse(ctx, "PostGateways", http.StatusInternalServerError, nsm.ErrGeneral)
		return
	}
	gc.Resources = as
	log.Trace(gc)

	// TODO delete from VIM
	// drivers.VimDriver.DeleteGatewayConnectivity()

	// Delete from DB - the resources and the gateway
	log.Info("DeleteGatewaysId - deleting gateway ", gc)
	obj.DB.Delete(gc)
	obj.DB.Delete(gc.Resources)
	ctx.Status(http.StatusNoContent)
}

// TODO implement SYNC go routine for creation and removal of resources
func createResources(database *gorm.DB, gc *nsm.Gateway) {
	time.Sleep(time.Second * 5)
	gc.Status = nsm.WAIT_FOR_GATEWAY_CONFIG
	log.Info("createResources ", gc.Status)
	database.Save(&gc)
}
