// Package to implement all of the handlers in the ServerInterface
package nbi

import (
	"errors"
	"net/http"
	nsmmapi "nextworks/nsm/api"
	"nextworks/nsm/internal/nsm"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// TODO add list of networks and saps
func setNetResourcesResponse(ctx *gin.Context, status int, res nsm.ResourceSet) {
	var netlist []nsmmapi.Network
	var saplist []nsmmapi.Sap

	for _, sap := range res.Networks {
		netlist = append(netlist, nsmmapi.Network{NetworkName: sap.NetworkName})
	}
	for _, sap := range res.Saps {
		saplist = append(saplist, nsmmapi.Sap{
			FloatingNetworkName: sap.FloatingNetName,
			NetworkName:         sap.NetworkName})
	}

	output := nsmmapi.SliceResources{
		Id:                  res.ID,
		Status:              res.Status,
		SliceId:             res.SliceId,
		VimName:             res.VimName,
		Networks:            netlist,
		ServiceAccessPoints: saplist,
	}
	ctx.JSON(status, output)
}

func (obj *ServerInterfaceImpl) GetNetResources(c *gin.Context, params nsmmapi.GetNetResourcesParams) {
	// var res nsm.ResourceSet
	// var result *gorm.DB

	if params.SliceId != nil {
		log.Info("GetGateways - filtered by sliceId ", *params.SliceId)
		// result = obj.DB.First(&gco, "slice_id = ?", params.SliceId)
		// setGatewayResponse(c, http.StatusOK, res)
		// TODO
	} else {
		log.Info("GetGateways - without query parameters")
		c.JSON(http.StatusOK, "{}")
		// TODO
	}

}

func (obj *ServerInterfaceImpl) PostNetResources(c *gin.Context) {
	var result *gorm.DB
	var jsonBody nsmmapi.PostSliceResources

	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		SetErrorResponse(c, "PostNetResources", http.StatusBadRequest, nsm.ErrBodyMissingInfo)
		return
	}
	// TODO check body params and format

	log.Debug("PostNetResources - requested creation of network resources for SliceId: ", jsonBody.SliceId, " on VIM: ", jsonBody.VimName)

	// Check if a Vim with this name exists
	// if !obj.Vims.Exists(jsonBody.VimName) {
	// 	SetErrorResponse(c, "PostNetResources", http.StatusForbidden, vim.ErrVimNotFound)
	// 	return
	// }

	// Check if resources for slice-id already exists
	var netres nsm.ResourceSet
	result = obj.DB.First(&netres, "slice_id = ?", jsonBody.SliceId)

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

	// scan resources and create them on the VIM
	log.Debug("PostNetResources creating requested resources for slice-id: ", jsonBody.SliceId)
	resset := &nsm.ResourceSet{
		Status:  nsm.CREATING,
		VimName: jsonBody.VimName,
		SliceId: jsonBody.SliceId,
	}

	result = obj.DB.Save(resset)
	if result.Error != nil {
		log.Error("PostNetResources - error writing in DB")
		SetErrorResponse(c, "", http.StatusInternalServerError, nsm.ErrGeneral)
		return
	}
	log.Debug("PostNetResources created resources set: ", resset.ID)

	// create networks:
	for _, net := range jsonBody.Networks {
		// TODO create on vim
		log.Info("PostNetResources - creating Network ", net)
		ne := nsm.Network{
			ResourceSetId: resset.ID,
			NetworkName:   net.NetworkName,
		}
		resset.Networks = append(resset.Networks, ne)
	}

	// create saps:
	for _, sap := range jsonBody.ServiceAccessPoints {
		// TODO crate on vim
		log.Info("PostNetResources - creating SAP ", sap)
		ap := nsm.Sap{
			ResourceSetId:   resset.ID,
			NetworkName:     sap.NetworkName,
			FloatingNetName: sap.FloatingNetworkName,
		}
		resset.Saps = append(resset.Saps, ap)
	}

	log.Info("---", resset)
	resset.Status = nsm.WAIT_FOR_GATEWAY_CONFIG
	log.Debug("PostNetResources updating resource set status for slice-id: ", resset.SliceId)
	result = obj.DB.Save(resset)
	if result.Error != nil {
		log.Error("PostNetResources - error writing in DB for slice: ", resset.SliceId)
		SetErrorResponse(c, "", http.StatusInternalServerError, nsm.ErrGeneral)
		return
	}

	setNetResourcesResponse(c, http.StatusOK, *resset)
}

func (obj *ServerInterfaceImpl) DeleteNetResources(c *gin.Context, params nsmmapi.DeleteNetResourcesParams) {
	var netres nsm.ResourceSet
	// var as nsm.OpenstackResource

	if params.SliceId == "" {
		log.Error("DeleteNetResources - no query param slice-id specified")
		SetErrorResponse(c, "DeleteNetResources", http.StatusBadRequest, nsm.ErrMissingQueryParameter)
	}

	log.Debug("DeleteNetResources - requesting deletion of all network resources for slice-id: ", params.SliceId)
	// Read from DB
	result := obj.DB.First(&netres, "slice_id = ?", params.SliceId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, "DeleteNetResources", http.StatusNotFound, nsm.ErrSliceNotExists)
			return
		}
		SetErrorResponse(c, "DeleteNetResources", http.StatusInternalServerError, nsm.ErrGeneral)
		return
	}
	// check status
	if netres.Status != nsm.WAIT_FOR_GATEWAY_CONFIG && netres.Status != nsm.CREATION_ERROR {
		log.Info("DeleteNetResources - impossible to delete network resources. The current state is ", netres.Status)
		SetErrorResponse(c, "DeleteGatewaysId", http.StatusForbidden, nsm.ErrResourcesCantBeDeleted)
		return
	}
	// delete
	// TODO delete from VIM

}

func (obj *ServerInterfaceImpl) GetNetResourcesId(c *gin.Context, id int) {
	resources, error := RetrieveResourcesFromDB(obj.DB, id)
	if error != nil {
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, "PostNetResources", http.StatusNotFound, nsm.ErrSliceNotExists)
			return
		} else {
			SetErrorResponse(c, "PostNetResources", http.StatusInternalServerError, nsm.ErrGeneral)
			return
		}
	}
	setNetResourcesResponse(c, http.StatusOK, *resources)
}

func (obj *ServerInterfaceImpl) DeleteNetResourcesId(c *gin.Context, id int) {
	// with onDelete cascade OK
	resources, error := RetrieveResourcesFromDB(obj.DB, id)
	if error != nil {
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, "DeleteNetResourcesId", http.StatusNotFound, nsm.ErrSliceNotExists)
			return
		} else {
			SetErrorResponse(c, "DeleteNetResourcesId", http.StatusInternalServerError, nsm.ErrGeneral)
			return
		}
	}
	// TODO check status, delete only on some state

	go deleteResources(obj.DB, resources)
	c.Status(http.StatusNoContent)
}

// TODO implement SYNC go routine for creation and removal of resources
func deleteResources(database *gorm.DB, res *nsm.ResourceSet) {
	time.Sleep(time.Second * 10)
	// with onDelete cascade OK
	_ = database.Delete(&res)
}
