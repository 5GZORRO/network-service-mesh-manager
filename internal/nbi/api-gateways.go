// Package to implement all of the handlers in the ServerInterface
package nbi

import (
	"errors"
	"net/http"
	NsmmApi "nextworks/nsm/api"
	"nextworks/nsm/internal/nsm"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setErrorResponse(ctx *gin.Context, method string, errorStatus int, err error) {
	log.Error(method, " - ", err.Error())
	outputJson := NsmmApi.ErrorResponse{Error: err.Error()}
	ctx.JSON(errorStatus, outputJson)
}

// TODO add subnet, retrieving from DB
func buildGatewayObjectResponse(gc nsm.Gateway, subnet string) NsmmApi.ResponseGatewayObject {
	return NsmmApi.ResponseGatewayObject{
		Id:      gc.ID,
		SliceId: gc.SliceID,
		Status:  gc.Status,
		Subnet:  subnet,
	}
}

// TODO add subnet, retrieving from DB
func setGatewayObjectResponse(ctx *gin.Context, status int, gc nsm.Gateway) {
	output := buildGatewayObjectResponse(gc, "")
	ctx.JSON(status, output)
}

// (GET /gateways)
// TODO returns a list of ResponsegatewayObject
func (obj *ServerInterfaceImpl) GetGateways(ctx *gin.Context, params NsmmApi.GetGatewaysParams) {
	log.Info("GetGateways")
	var gco nsm.Gateway
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
		setGatewayObjectResponse(ctx, http.StatusOK, gco)
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
		// TODO define response object

		ctx.JSON(http.StatusOK, gin.H{"gateways": gcos})
	}
}

func (obj *ServerInterfaceImpl) PostGateways(ctx *gin.Context) {
	// Retrieve body
	var jsonBody NsmmApi.PostGatewayBody
	if err := ctx.ShouldBindJSON(&jsonBody); err != nil {
		log.Error(nsm.ErrBodyMissingInfo.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": nsm.ErrBodyMissingInfo.Error()})
		return
	}
	// TODO check body params and format ex. subnet and in case return 400
	// TODO check if subnet is already in use

	log.Info("PostGateways - requesting creation of Gateway connection for SliceId: ", jsonBody.SliceId, " with subnet ", jsonBody.Subnet)
	// Check if a gateway for the slice already exists
	var gco nsm.Gateway
	result := obj.DB.First(&gco, "slice_id = ?", jsonBody.SliceId)

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		setErrorResponse(ctx, "PostGateways", http.StatusForbidden, nsm.ErrGatewayExists)
		return
	}
	// gateway does not exists, create a new one
	log.Info("PostGateways - gateway for requested slice not found")
	gc := nsm.Gateway{SliceID: jsonBody.SliceId, Status: nsm.CREATING}
	result = obj.DB.Create(&gc)

	if result.RowsAffected != 1 {
		log.Error("PostGateways - error saving gateway info")
		ctx.Status(http.StatusInternalServerError)
		return
	}
	log.Info("PostGateways - gateway for requested slice initialized")
	// TODO create il gatewayconnectivity su Openstack and update all info
	// drivers.VimDriver.CreateGatewayConnectivity()

	setGatewayObjectResponse(ctx, http.StatusOK, gc)
}

func (obj *ServerInterfaceImpl) GetGatewaysId(ctx *gin.Context, id int) {
	var gc nsm.Gateway
	log.Info("GetGatewaysId - requested Gateway info with ID: ", id)
	// Read from DB, check if exists
	result := obj.DB.First(&gc, id)

	log.Info(gc)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		setErrorResponse(ctx, "GetGatewaysId", http.StatusNotFound, nsm.ErrGatewayNotFound)
		return
	} else {
		setGatewayObjectResponse(ctx, http.StatusOK, gc)
		return
	}
}

func (obj *ServerInterfaceImpl) DeleteGatewaysId(ctx *gin.Context, id int) {
	var gc nsm.Gateway

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
	if gc.Status != nsm.WAIT_FOR_GATEWAY && gc.Status != nsm.CREATION_ERROR {
		log.Info("DeleteGatewaysId - impossible to delete gateway. The current state is ", gc.Status)
		setErrorResponse(ctx, "DeleteGatewaysId", http.StatusForbidden, nsm.ErrGatewayCantBeDeleted)
		return
	}
	// TODO delete from VIM
	// drivers.VimDriver.DeleteGatewayConnectivity()

	// Delete from DB
	log.Info("DeleteGatewaysId - deleting gateway ", gc)
	obj.DB.Delete(gc)
	ctx.Status(http.StatusNoContent)
}
