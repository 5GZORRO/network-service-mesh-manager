// Package to implement all of the handlers in the ServerInterface
package nbi

import (
	"errors"
	"net/http"
	nsmmapi "nextworks/nsm/api"
	"nextworks/nsm/internal/nsm"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func setGatewayResponse(ctx *gin.Context, status int, res nsm.ResourceSet) {
	output := nsmmapi.SliceResources{Status: "gyyf"}

	ctx.JSON(status, output)
}

func (obj *ServerInterfaceImpl) GetNetResources(c *gin.Context, params nsmmapi.GetNetResourcesParams) {
	var res nsm.ResourceSet
	// var result *gorm.DB

	if params.SliceId != nil {
		log.Info("GetGateways - filtered by sliceId ", *params.SliceId)
		// result = obj.DB.First(&gco, "slice_id = ?", params.SliceId)

		setGatewayResponse(c, http.StatusOK, res)
	} else {
		log.Info("GetGateways - without query parameters")
		c.JSON(http.StatusOK, "{}")
	}
}

func (obj *ServerInterfaceImpl) PostNetResources(c *gin.Context) {
	var jsonBody nsmmapi.PostSliceResources
	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		SetErrorResponse(c, "PostNetResources", http.StatusBadRequest, nsm.ErrBodyMissingInfo)
		return
	}
	// TODO check body params and format ex. subnet and in case return 400
	log.Debug("PostNetResources - requesting creation of network resources for SliceId: ", jsonBody.SliceId, " on VIM '", jsonBody.VimName)

	// Check if a Vim with this name exists
	// if !obj.Vims.Exists(jsonBody.VimName) {
	// 	SetErrorResponse(c, "PostNetResources", http.StatusForbidden, vim.ErrVimNotFound)
	// 	return
	// }

	// Check if resources for slice-id already exists
	var netres nsm.ResourceSet
	result := obj.DB.First(&netres, "slice_id = ?", jsonBody.SliceId)

	if result.Error == nil {
		SetErrorResponse(c, "PostNetResources", http.StatusForbidden, nsm.ErrSliceExists)
		return
	} else {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, "PostNetResources", http.StatusInternalServerError, nsm.ErrGeneral)
			return
		}
	}
	// resources for requested slice do not exist
}

func (obj *ServerInterfaceImpl) DeleteNetResources(c *gin.Context, params nsmmapi.DeleteNetResourcesParams) {
}

func (obj *ServerInterfaceImpl) GetNetResourcesId(c *gin.Context, id int) {}

func (obj *ServerInterfaceImpl) DeleteNetResourcesId(c *gin.Context, id int) {}
