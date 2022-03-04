package nsm

import (
	"errors"
	"net/http"
	vimdriver "nextworks/nsm/internal/vim"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (obj *ServerInterfaceImpl) PutNetResourcesIdGatewayExternalIp(c *gin.Context, id int) {
	log.Trace("PutNetResourcesIdGatewayExternalIp - requested allocation of floating IP for ResourceSet with ID: ", id)

	// retrieve resource
	resource, error := RetrieveResourcesFromDB(obj.DB, id)
	if error != nil {
		log.Error("Impossible to allocate external-IP to gateway. Error reading from DB: ", error)
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, http.StatusNotFound, ErrResourcesNotExists)
			return
		} else {
			SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
			return
		}
	}

	// check state - it should be CREATED
	if resource.Status != CREATED {
		log.Error("Impossibile to sallocate external gateway IP. The current state is ", resource.Status)
		SetErrorResponse(c, http.StatusForbidden, ErrConfiguringGateway)
		return
	}

	// retrieve VIM
	vim, err := obj.Vims.GetVim(resource.VimName)
	if err == vimdriver.ErrVimNotFound {
		log.Error("VIM of resource set with id ", resource.ID, " does not exist")
		SetErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	log.Trace("Allocating floating IP... ")
	// ASSOCIATE FLOATING-IP USING VIM
	// For each SAP network, try to allocate a floating IP for each compute, found on that network
	// retrieve all SAP networks of this resource set
	err = LoadSAPsFromDB(obj.DB, resource)
	if err != nil {
		SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
		return
	}
	log.Trace("Loaded SAP: ", resource.Saps)
	//for the SAP matching the NetworkPrefix search a VM and allocate a floating-ip
	var selectedSap *Sap = nil
	for _, sap := range resource.Saps {
		log.Trace("Loaded SAP: ", sap.NetworkName)
		if strings.HasPrefix(sap.NetworkName, obj.Netconfig.GatewayNetworkNamePrefix) {
			log.Info("Selected SAP network with name ", sap.NetworkName)
			selectedSap = &sap
			break
		}
	}
	if selectedSap == nil {
		log.Error("No SAP network found with the specified prefix-name ", obj.Netconfig.GatewayNetworkNamePrefix, " for resource set with ID: ", resource.ID)
		SetErrorResponse(c, http.StatusInternalServerError, ErrGatewayNoNetworkFound)
		return
	} else {
		portid, portname, fipid, fip, err := (*vim).AllocateFloatingIP(selectedSap.NetworkId)
		if err != nil {
			log.Error("Error allocating FloatingIP for gateway of resource set with ID: ", resource.ID)
			SetErrorResponse(c, http.StatusInternalServerError, err)
			return
		} else {
			log.Trace(" PortID: ", portid, " PortName: ", portname, " FloatingID: ", fipid, " FloatingIP: ", fip)
			ext := ExternalIP{}
			ext.ExternalIp = fip
			ext.PortID = portid
			ext.PortName = portname
			ext.FloatingID = fipid
			resource.Gateway.External = ext
			resource.Status = WAIT_FOR_GATEWAY_CONFIG
			result := obj.DB.Save(&resource)
			if result.Error != nil {
				log.Error("Error updating resource set status with ID: ", resource.ID, " and slice-id: ", resource.SliceId)
			}
			c.Status(http.StatusCreated)
		}
	}

}

func (obj *ServerInterfaceImpl) DeleteNetResourcesIdGatewayExternalIp(c *gin.Context, id int) {
	log.Trace("DeleteNetResourcesIdGatewayExternalIp - requested de-allocation of floating IP for ResourceSet with ID: ", id)

	resource, error := RetrieveResourcesFromDB(obj.DB, id)
	if error != nil {
		log.Error("Impossible to de-allocate external gateway IP. Error reading from DB: ", error)
		if errors.Is(error, gorm.ErrRecordNotFound) {
			SetErrorResponse(c, http.StatusNotFound, ErrResourcesNotExists)
			return
		} else {
			SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
			return
		}
	}

	// check state - it should be WAIT_FOR_CONFIG
	if resource.Status != WAIT_FOR_GATEWAY_CONFIG {
		log.Error("Impossibile to de-allocate external gateway IP. The current state is ", resource.Status)
		SetErrorResponse(c, http.StatusForbidden, ErrDeleteExternalIP)
		return
	}

	// retrieve VIM
	vim, err := obj.Vims.GetVim(resource.VimName)
	if err == vimdriver.ErrVimNotFound {
		log.Error("VIM of resource set with id ", resource.ID, " does not exist")
		SetErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	// deallocate FIP, if it exists
	if resource.Gateway.External.FloatingID != "" {
		log.Trace("Deallocating and deleting floatingIP...")
		err := (*vim).DeallocateFloatingIP(resource.Gateway.External.PortID, resource.Gateway.External.FloatingID)
		if err != nil {
			log.Error("Error deallocating/deleting floating IP with ID: ", resource.Gateway.External.FloatingID)
		}
	} else {
		log.Info("No floatingIP to be deallocated/deleted")
	}

	// update the state of the gateway to CREATED, and reset all fields of Gateway
	resource.Status = CREATED
	resource.Gateway = Gateway{}
	result := obj.DB.Save(&resource)
	if result.Error != nil {
		log.Error("Error updating resource set status with ID: ", resource.ID, " and slice-id: ", resource.SliceId)
	}

}
