// Package to implement all of the handlers in the ServerInterface
package NsmmApi

import (
	"errors"
	"net/http"
	"nextworks/nsm/internal/nsm"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Shared object between different HTTP REST handlers
// it should contain the DBconnection
type ServerInterfaceImpl struct {
	DB *gorm.DB
	// Gateways string
	// Lock     sync.Mutex
}

func NewServerInterfaceImpl(DBconnection *gorm.DB) *ServerInterfaceImpl {
	return &ServerInterfaceImpl{
		DB: DBconnection,
	}
}

func (obj *ServerInterfaceImpl) GetGateway(ctx *gin.Context, params GetGatewayParams) {
	var gc nsm.Gateway
	sliceId := params.SliceId
	log.Info("GetGateway - requested SliceId: " + sliceId)
	// Read from DB, check if exists
	result := obj.DB.First(&gc, "slice_id = ?", sliceId)

	// TODO decide what should be returned in the response body
	log.Info(gc)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.Status(http.StatusNotFound)
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"gateway": gc})
	}
}

func (obj *ServerInterfaceImpl) PostGateway(ctx *gin.Context, params PostGatewayParams) {
	sliceId := params.SliceId

	// Retrieve body
	var jsonBody PostGatewayJSONRequestBody
	if err := ctx.ShouldBindJSON(&jsonBody); err != nil {
		log.Error(nsm.ErrMissingInfo.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": nsm.ErrMissingInfo.Error()})
		return
	}
	// TODO check body params and format ex. subnet

	log.Info("PostGateway - requesting creation of Gateway connection for SliceId: ", sliceId, " with subnet ", jsonBody.Subnet)
	// Check if a gateway for the slice already exists
	var gco nsm.Gateway
	result := obj.DB.First(&gco, "slice_id = ?", sliceId)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// gateway does not exists, create a new one
		log.Info("PostGateway - gateway for requested slice not found")
		gc := nsm.Gateway{SliceID: sliceId, Status: nsm.CREATING}
		result := obj.DB.Create(&gc)

		if result.RowsAffected == 1 {
			log.Info("PostGateway - gateway for requested slice initialized")
			ctx.Status(http.StatusOK)
		} else {
			log.Error("PostGateway - error saving gateway info")
			ctx.Status(http.StatusInternalServerError)
		}
		// TODO TBF
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": nsm.ErrGatewayExists.Error()})
	}
}

// TODO to be implemented and define new errors
func (obj *ServerInterfaceImpl) DeleteGateway(ctx *gin.Context, params DeleteGatewayParams) {

	sliceId := params.SliceId
	var gc nsm.Gateway
	log.Info("DeleteGateway - requesting removal of Gateway connection for SliceId: " + sliceId)
	// Read from DB
	result := obj.DB.First(&gc, "slice_id = ?", sliceId)

	log.Info(gc)
	if result.Error != nil {
		log.Error("DeleteGateway - error retrieving gateway for requested slice")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Error("DeleteGateway - gateway for requested slice not found")
			ctx.Status(http.StatusNotFound)
		}
	}
	// check status
	if gc.Status == nsm.CREATING {
		log.Info("DeleteGateway - gateway for requested slice is in state ", nsm.CREATING)
		ctx.Status(http.StatusInternalServerError)
	} else {
		// TODO Delete from DB
		log.Info("DeleteGateway - deleting gateway")
		ctx.Status(http.StatusNoContent)
	}
}
