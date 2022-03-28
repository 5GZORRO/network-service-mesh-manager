package nsm

import (
	"errors"
	"net/http"
	nsmmapi "nextworks/nsm/api"
	vimdriver "nextworks/nsm/internal/vim"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// PostNetResourcesStaticSap creates a new set of resources and creates them on the VIM
// excluding some CIDR passed as optional parameters
func (obj *ServerInterfaceImpl) PostNetResourcesStaticSap(c *gin.Context) {
	var result *gorm.DB
	var jsonBody nsmmapi.PostSliceResourcesFixedSap

	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		log.Error("Impossible to create network resources. Error in the request, wrong json body")
		SetErrorResponse(c, http.StatusBadRequest, ErrBodyWrongInfo)
		return
	}
	log.Trace("Requested creation of network resources with STATIC SAP for SliceId: ", jsonBody.SliceId, " on VIM: ", jsonBody.VimName)
	if jsonBody.ExcludeSubnet != nil {
		log.Trace("with excluded subnet: ", *jsonBody.ExcludeSubnet)
		err := checkExcludedSubnetsParams(jsonBody.ExcludeSubnet)
		if err != nil {
			log.Error("Impossible to create network resources. Error in JSON Body: exclude-subnets")
			SetErrorResponse(c, http.StatusBadRequest, err)
			return
		}
	}

	// Check if a Vim with this name exists, and retrieve it
	// retrieve VIM
	vim, err := obj.Vims.GetVim(jsonBody.VimName)
	if err != nil {
		log.Error("Impossible to create network resources. Vim with name: ", jsonBody.VimName, " does not exist")
		SetErrorResponse(c, http.StatusNotFound, vimdriver.ErrVimNotFound)
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
		Status:    CREATING,
		VimName:   jsonBody.VimName,
		SliceId:   jsonBody.SliceId,
		StaticSap: true,
	}

	result = obj.DB.Save(resset)
	if result.Error != nil {
		log.Error("Impossible to create network resources. Error writing in DB")
		SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
		return
	}
	log.Trace("PostNetResourcesFixedSap saved resources set: ", *resset)

	// management of allocated networks
	var netmng *Network_manager
	if jsonBody.ExcludeSubnet == nil {
		netmng = NewNetworkManager(obj.Netconfig.Start, jsonBody.ExcludeSubnet != nil)
	} else {
		netmng = NewNetworkManager(*jsonBody.ExcludeSubnet, jsonBody.ExcludeSubnet != nil)
	}
	// log.Info(netmng.next_subnet)

	// create networks:
	for _, net := range jsonBody.Networks {
		ipnet := netmng.NextSubnet()
		if ipnet != nil {
			cidr := ipnet.String()
			netID, subnetID, subnetName, err := (*vim).CreateNetwork(net.NetworkName, cidr, false)
			ne := Network{
				ResourceSetId: resset.ID,
				NetworkId:     netID,
				NetworkName:   net.NetworkName,
				SubnetId:      subnetID,
				SubnetName:    subnetName,
				SubnetCidr:    cidr,
			}
			resset.Networks = append(resset.Networks, ne)
			if err != nil {
				log.Error("Impossible to create network resources. Error creating network ", net.NetworkName)
				SetErrorResponse(c, http.StatusInternalServerError, ErrVimCreatingNetwork)
				resset.Status = CREATION_ERROR
				_ = obj.DB.Save(resset)
				return
			} else {
				log.Trace("PostNetResourcesFixedSap - created Network ", net.NetworkName, " with cidr: ", cidr)
			}
		} else {
			log.Error("Impossible to create network resources. Error allocating IP addresses")
			SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
			resset.Status = CREATION_ERROR
			_ = obj.DB.Save(resset)
			return
		}
	}

	// the sap is static, so no operation on the vim is invoked, only store information coming from API
	for _, sap := range jsonBody.ServiceAccessPoints {
		ap := Sap{
			ResourceSetId:   resset.ID,
			NetworkId:       sap.NetworkId,
			NetworkName:     sap.NetworkName,
			SubnetId:        sap.SubnetId,
			SubnetName:      sap.SubnetName,
			SubnetCidr:      sap.SubnetCidr,
			RouterId:        "",
			RouterName:      "",
			RouterPortId:    "",
			FloatingNetID:   (*vim).RetrieveFloatingNetworkID(),
			FloatingNetName: (*vim).RetrieveFloatingNetworkName(),
		}
		resset.Saps = append(resset.Saps, ap)
		log.Trace("PostNetResourcesFixedSap - stored SAP ", sap.NetworkName, " with cidr: "+sap.SubnetCidr)
	}

	ext := ExternalIP{}
	ext.ExternalIp = jsonBody.GwExternalIp
	ext.PortID = jsonBody.GwPortName
	ext.PortName = ""
	ext.FloatingID = ""
	resset.Gateway.External = ext
	resset.Status = WAIT_FOR_GATEWAY_CONFIG

	log.Trace("ResourceSet with relations infos: ", *resset)
	resset.Status = WAIT_FOR_GATEWAY_CONFIG
	result = obj.DB.Save(resset)
	if result.Error != nil {
		log.Error("Impossible to create network resources. Error saving in DB")
		SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
		return
	}
	log.Trace("PostNetResourcesFixedSap updated ", *resset)

	SetNetResourcesResponse(c, http.StatusCreated, *resset)

}
