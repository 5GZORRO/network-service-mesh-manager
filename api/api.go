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

//
func (obj *ServerInterfaceImpl) GetGateway(ctx *gin.Context, params GetGatewayParams) {
	var gc nsm.Gateway
	sliceId := params.SliceId
	log.Info("GetGateway - requested SliceId: " + sliceId)
	result := obj.DB.First(&gc, "slice_id = ?", sliceId)

	log.Info(gc)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.Status(http.StatusNotFound)
	} else {
		ctx.Status(http.StatusAccepted)
	}
}

func (obj *ServerInterfaceImpl) PostGateway(ctx *gin.Context, params PostGatewayParams) {

	ctx.Status(http.StatusNoContent)
}

func (obj *ServerInterfaceImpl) DeleteGateway(ctx *gin.Context, params DeleteGatewayParams) {
	ctx.Status(http.StatusNoContent)
}
