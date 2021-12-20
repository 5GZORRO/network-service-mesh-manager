package nsm

import (
	"errors"
	"net/http"
	nsmmapi "nextworks/nsm/api"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// GetNetResources retrieves all the network resources created on the VIMs or the set
// of network resources created for a slice, using the query param
func (obj *ServerInterfaceImpl) GetNetResources(c *gin.Context, params nsmmapi.GetNetResourcesParams) {

	if params.SliceId != nil {
		log.Trace("GetNetResources - requested retrieval of network resources for slice-id: ", *params.SliceId)
		res, result := RetrieveResourcesFromDBbySliceID(obj.DB, *params.SliceId)
		if result != nil {
			if errors.Is(result, gorm.ErrRecordNotFound) {
				log.Error(result)
				SetErrorResponse(c, "GetNetResources", http.StatusNotFound, ErrSliceNotExists)
				return
			}
			SetErrorResponse(c, "GetNetResources", http.StatusInternalServerError, ErrGeneral)
			return
		}
		SetNetResourcesResponse(c, http.StatusOK, *res)
	} else {
		log.Trace("GetNetResources - requested retrieval of all network resources")
		var resources []ResourceSet
		result := obj.DB.Find(&resources)
		if result.Error != nil {
			SetErrorResponse(c, "GetNetResources", http.StatusInternalServerError, ErrGeneral)
			return
		}
		for i := range resources {
			err := LoadAssociationFromDB(obj.DB, &resources[i])
			if err != nil {
				log.Error("GetNetResources error retrieving associations of resource set with ID: ", resources[i].ID)
				SetErrorResponse(c, "GetNetResources", http.StatusInternalServerError, ErrGeneral)
				return
			}
			log.Trace(resources[i])
		}
		log.Trace(resources)
		SetNetResourcesListResponse(c, http.StatusOK, resources)
	}

}

// PostNetResources creates a new set of resources and creates them on the VIM
func (obj *ServerInterfaceImpl) PostNetResources(c *gin.Context) {
	var result *gorm.DB
	var jsonBody nsmmapi.PostSliceResources

	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		SetErrorResponse(c, "PostNetResources", http.StatusBadRequest, ErrBodyMissingInfo)
		return
	}
	log.Trace("PostNetResources - requested creation of network resources for SliceId: ", jsonBody.SliceId, " on VIM: ", jsonBody.VimName)

	// Check if a Vim with this name exists
	// if !obj.Vims.Exists(jsonBody.VimName) {
	// 	SetErrorResponse(c, "PostNetResources", http.StatusForbidden, vim.ErrVimNotFound)
	// 	return
	// }

	// Check if resources for slice-id already exists
	var netres ResourceSet
	result = obj.DB.First(&netres, "slice_id = ?", jsonBody.SliceId)

	if result.Error == nil {
		SetErrorResponse(c, "PostNetResources", http.StatusForbidden, ErrSliceExists)
		return
	} else {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, "PostNetResources", http.StatusInternalServerError, ErrGeneral)
			return
		}
	}
	// resources for requested slice do not exist

	// scan resources and create them on the VIM
	resset := &ResourceSet{
		Status:  CREATING,
		VimName: jsonBody.VimName,
		SliceId: jsonBody.SliceId,
	}

	result = obj.DB.Save(resset)
	if result.Error != nil {
		log.Error("PostNetResources - error writing in DB")
		SetErrorResponse(c, "", http.StatusInternalServerError, ErrGeneral)
		return
	}
	log.Trace("PostNetResources saved resources set: ", *resset)

	// create networks:
	for _, net := range jsonBody.Networks {
		// TODO create on vim
		log.Trace("PostNetResources - creating Network ", net)
		ne := Network{
			ResourceSetId: resset.ID,
			NetworkName:   net.NetworkName,
		}
		resset.Networks = append(resset.Networks, ne)
	}

	// create saps:
	for _, sap := range jsonBody.ServiceAccessPoints {
		// TODO crate on vim
		log.Trace("PostNetResources - creating SAP ", sap)
		ap := Sap{
			ResourceSetId:   resset.ID,
			NetworkName:     sap.NetworkName,
			FloatingNetName: sap.FloatingNetworkName,
		}
		resset.Saps = append(resset.Saps, ap)
	}

	log.Trace("ResourceSet with additional infos: ", *resset)
	resset.Status = WAIT_FOR_GATEWAY_CONFIG
	result = obj.DB.Save(resset)
	if result.Error != nil {
		log.Error("PostNetResources - error writing in DB for slice: ", resset.SliceId)
		SetErrorResponse(c, "", http.StatusInternalServerError, ErrGeneral)
		return
	}
	log.Debug("PostNetResources resource set with ID: ", resset.ID, " and Slice-ID: ", resset.SliceId, "updated")
	SetNetResourcesResponse(c, http.StatusOK, *resset)
}

// DeleteNetResources start the removal of all the network resources associated to a slice-id created on the VIM
// only if they are in WAIT_FOR_GATEWAY or CREATION_ERROR
// and set the resource set status to DELETING
// the actual removal is done in an async. way by a dedicated go routine
func (obj *ServerInterfaceImpl) DeleteNetResources(c *gin.Context, params nsmmapi.DeleteNetResourcesParams) {
	var netres ResourceSet

	if params.SliceId == "" {
		log.Error("DeleteNetResources - no query param slice-id specified")
		SetErrorResponse(c, "DeleteNetResources", http.StatusBadRequest, ErrMissingQueryParameter)
		return
	}

	log.Trace("DeleteNetResources - Received request to delete network resources for slice-id: ", params.SliceId)
	// Read from DB
	result := obj.DB.First(&netres, "slice_id = ?", params.SliceId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, "DeleteNetResources", http.StatusNotFound, ErrSliceNotExists)
			return
		}
		SetErrorResponse(c, "DeleteNetResources", http.StatusInternalServerError, ErrGeneral)
		return
	}
	// check status
	if netres.Status != WAIT_FOR_GATEWAY_CONFIG && netres.Status != CREATION_ERROR {
		log.Info("DeleteNetResources - impossible to delete network resources. The current state is ", netres.Status)
		SetErrorResponse(c, "DeleteGatewaysId", http.StatusForbidden, ErrResourcesCantBeDeleted)
		return
	}
	log.Trace("DeleteNetResources - deleting network resources for slice-id: ", params.SliceId)
	// Set current stats to DELETING and the async. delete them from VIM and DB
	// and all the resources associated to the set
	netres.Status = DELETING_RESOURCES
	result = obj.DB.Save(netres)
	if result.Error != nil {
		SetErrorResponse(c, "DeleteNetResources", http.StatusInternalServerError, ErrGeneral)
		return
	}
	log.Trace("DeleteNetResources - starting asynch routine for deleting resources of Slice-ID: ", params.SliceId)
	go deleteResources(obj.DB, &netres)
	c.Status(http.StatusOK)
}

// GetNetResourcesId retrieves the info of a set of network resources by its ID
func (obj *ServerInterfaceImpl) GetNetResourcesId(c *gin.Context, id int) {
	resources, error := RetrieveResourcesFromDB(obj.DB, id)
	if error != nil {
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, "PostNetResources", http.StatusNotFound, ErrSliceNotExists)
			return
		} else {
			SetErrorResponse(c, "PostNetResources", http.StatusInternalServerError, ErrGeneral)
			return
		}
	}
	SetNetResourcesResponse(c, http.StatusOK, *resources)
}

// DeleteNetResourcesId start the removal of all the network resources by its ID created on the VIM
// only if they are in WAIT_FOR_GATEWAY or CREATION_ERROR
// and set the resource set status to DELETING
// the actual removal is done in an async. way by a dedicated go routine
func (obj *ServerInterfaceImpl) DeleteNetResourcesId(c *gin.Context, id int) {
	log.Trace("DeleteNetResourcesId - Received request to delete network resources with ID: ", id)
	netres, error := RetrieveResourcesFromDB(obj.DB, id)
	if error != nil {
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, "DeleteNetResourcesId", http.StatusNotFound, ErrSliceNotExists)
			return
		} else {
			SetErrorResponse(c, "DeleteNetResourcesId", http.StatusInternalServerError, ErrGeneral)
			return
		}
	}
	// check status
	if netres.Status != WAIT_FOR_GATEWAY_CONFIG && netres.Status != CREATION_ERROR {
		log.Info("DeleteNetResourcesId - impossible to delete network resources. The current state is ", netres.Status)
		SetErrorResponse(c, "DeleteNetResourcesId", http.StatusForbidden, ErrResourcesCantBeDeleted)
		return
	}
	log.Trace("DeleteNetResourcesId - removing network resources with ID: ", id)
	// Set currrent stats to DELETING and the async. delete them from VIM and DB
	// and all the resources associated to the set
	netres.Status = DELETING_RESOURCES
	result := obj.DB.Save(netres)
	if result.Error != nil {
		SetErrorResponse(c, "DeleteNetResourcesId", http.StatusInternalServerError, ErrGeneral)
		return
	}
	log.Trace("DeleteNetResourcesId - staring asynch job to delete network resource set with ID: ", id)
	go deleteResources(obj.DB, netres)
	c.Status(http.StatusOK)
}

// deleteResources is a goroutine to delete in an async way all the network resources
func deleteResources(database *gorm.DB, res *ResourceSet) {
	time.Sleep(time.Second * 5)
	// TODO Delete resources from VIM

	// if removal from VIM is OK then delete it from DB
	result := database.Delete(&res)
	if result.Error != nil {
		log.Error("Error deleting network resource set with ID: ", res.ID, " and slice-id: ", res.SliceId)
	}
}
