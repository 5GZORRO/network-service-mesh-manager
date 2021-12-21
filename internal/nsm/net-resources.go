package nsm

import (
	"errors"
	"net/http"
	nsmmapi "nextworks/nsm/api"
	"nextworks/nsm/internal/vim"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// GetNetResources retrieves all the network resources created on the VIMs or the set
// of network resources created for a slice, using the query param
func (obj *ServerInterfaceImpl) GetNetResources(c *gin.Context, params nsmmapi.GetNetResourcesParams) {

	if params.SliceId != nil {
		log.Trace("Requested retrieval of network resources for slice-id: ", *params.SliceId)
		res, result := RetrieveResourcesFromDBbySliceID(obj.DB, *params.SliceId)
		if result != nil {
			log.Error("Impossible to retrieve network resources. Error reading from DB: ", result)
			if errors.Is(result, gorm.ErrRecordNotFound) {
				log.Error(result)
				SetErrorResponse(c, http.StatusNotFound, ErrSliceNotExists)
				return
			}
			SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
			return
		}
		SetNetResourcesResponse(c, http.StatusOK, *res)
	} else {
		log.Trace("Requested retrieval of all network resources")
		var resources []ResourceSet
		result := obj.DB.Find(&resources)
		if result.Error != nil {
			log.Error("Impossible to retrieve network resources. Error reading from DB: ", result.Error)
			SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
			return
		}
		for i := range resources {
			err := LoadAssociationFromDB(obj.DB, &resources[i])
			if err != nil {
				log.Error("Error retrieving associations of resource set with ID: ", resources[i].ID)
				SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
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
		log.Error("Impossible to create network resources. Error in the request, wrong json body")
		SetErrorResponse(c, http.StatusBadRequest, ErrBodyMissingInfo)
		return
	}
	log.Trace("Requested creation of network resources for SliceId: ", jsonBody.SliceId, " on VIM: ", jsonBody.VimName)

	// TODO select vim
	// Check if a Vim with this name exists
	if !obj.Vims.Exists(jsonBody.VimName) {
		log.Error("Impossible to create network resources. Vim with name: ", jsonBody.VimName, " does not exist")
		SetErrorResponse(c, http.StatusForbidden, vim.ErrVimNotFound)
		return
	}

	// Check if resources for slice-id already exists
	var netres ResourceSet
	result = obj.DB.First(&netres, "slice_id = ?", jsonBody.SliceId)

	if result.Error == nil {
		log.Error("Impossible to create network resources. Network resources for SliceID: ", jsonBody.SliceId, " already exist")
		SetErrorResponse(c, http.StatusForbidden, ErrSliceExists)
		return
	} else {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Error("Impossible to create network resources. Error reading from DB")
			SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
			return
		}
	}
	// scan resources and create them on the VIM
	resset := &ResourceSet{
		Status:  CREATING,
		VimName: jsonBody.VimName,
		SliceId: jsonBody.SliceId,
	}

	result = obj.DB.Save(resset)
	if result.Error != nil {
		log.Error("Impossible to create network resources. Error reading from DB")
		SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
		return
	}
	log.Trace("PostNetResources saved resources set: ", *resset)

	// retrieve VIM
	vim := obj.Vims.GetVim(jsonBody.VimName)

	// create networks:
	for _, net := range jsonBody.Networks {
		// TODO create on vim
		(*vim).CreateNetwork()
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
		(*vim).CreateSAP()
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
		log.Error("Impossible to create network resources. Error saving in DB")
		SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
		return
	}
	log.Trace("PostNetResources updated ", *resset)
	SetNetResourcesResponse(c, http.StatusOK, *resset)
}

// DeleteNetResources start the removal of all the network resources associated to a slice-id created on the VIM
// only if they are in WAIT_FOR_GATEWAY or CREATION_ERROR
// and set the resource set status to DELETING
// the actual removal is done in an async. way by a dedicated go routine
func (obj *ServerInterfaceImpl) DeleteNetResources(c *gin.Context, params nsmmapi.DeleteNetResourcesParams) {
	var netres ResourceSet

	if params.SliceId == "" {
		log.Error("Impossible to delete network resources. Error in the request, missing slice-id query param")
		SetErrorResponse(c, http.StatusBadRequest, ErrMissingQueryParameter)
		return
	}

	log.Trace("Received request to delete network resources for slice-id: ", params.SliceId)
	// Read from DB
	result := obj.DB.First(&netres, "slice_id = ?", params.SliceId)

	if result.Error != nil {
		log.Error("Impossible to delete network resources. Error reading from DB: ", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, http.StatusNotFound, ErrSliceNotExists)
			return
		}
		SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
		return
	}
	obj.deleteNetResources(c, &netres)
}

// GetNetResourcesId retrieves the info of a set of network resources by its ID
func (obj *ServerInterfaceImpl) GetNetResourcesId(c *gin.Context, id int) {
	resources, error := RetrieveResourcesFromDB(obj.DB, id)
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
		log.Error("Impossible to delete network resources. Error reading from DB: ", error)
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, http.StatusNotFound, ErrSliceNotExists)
			return
		} else {
			SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
			return
		}
	}
	obj.deleteNetResources(c, netres)
}

func (obj *ServerInterfaceImpl) deleteNetResources(c *gin.Context, netres *ResourceSet) {
	// check status
	if netres.Status != WAIT_FOR_GATEWAY_CONFIG && netres.Status != CREATION_ERROR {
		log.Error("Impossible to delete network resources. The current state is ", netres.Status)
		SetErrorResponse(c, http.StatusForbidden, ErrResourcesCantBeDeleted)
		return
	}

	// Check/Retrieve VIM
	vim := obj.Vims.GetVim(netres.VimName)
	if *vim == nil {
		log.Error("Network resources cannot be canceled: vim ", netres.VimName, " does not exist")
		SetErrorResponse(c, http.StatusInternalServerError, ErrVimNotExists)
		return
	}

	log.Trace("Removing network resources with ID: ", netres.ID)
	log.Trace("Removing network resources - ", *netres)

	// Set current state to DELETING and the async. delete them from VIM and DB
	// and all the resources associated to the set
	netres.Status = DELETING_RESOURCES
	result := obj.DB.Save(netres)
	if result.Error != nil {
		log.Error("Network resources cannot be canceled: error saving on DB")
		SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
		return
	}
	log.Trace("Removing network resources - staring asynch job to delete network resource set with ID: ", netres.ID)
	go deleteResources(obj.DB, vim, netres)
	c.Status(http.StatusOK)
}
