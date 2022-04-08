package nsm

import (
	"errors"
	"net/http"
	nsmmapi "nextworks/nsm/api"
	vimdriver "nextworks/nsm/internal/vim"
	"strings"

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
	// Being a Static GW resource creation, check if the VIM has a static GW associated
	gwNetworkName, gwSubnetCidr, gwExternalInterfaceName, gwFloatingIP, gwInstanceID, err := (*vim).GetStaticGatewayInfo()
	if err != nil {
		log.Error("Impossible to create network resources with static GW. Static GW not specified for VIM: ", jsonBody.VimName)
		SetErrorResponse(c, http.StatusNotFound, vimdriver.ErrStaticGatewayNotFound)
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
		StaticGW: StaticGatewayAdditionalInfo{
			Enabled:    true,
			InstanceID: gwInstanceID,
		},
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
	ap := Sap{
		ResourceSetId:   resset.ID,
		NetworkId:       "",
		NetworkName:     gwNetworkName,
		SubnetId:        "",
		SubnetName:      "",
		SubnetCidr:      gwSubnetCidr,
		RouterId:        "",
		RouterName:      "",
		RouterPortId:    "",
		FloatingNetID:   (*vim).RetrieveFloatingNetworkID(),
		FloatingNetName: (*vim).RetrieveFloatingNetworkName(),
	}
	resset.Saps = append(resset.Saps, ap)
	log.Trace("PostNetResourcesFixedSap - stored SAP ", ap.NetworkName, " with cidr: "+ap.SubnetCidr)

	// Add the additional interface on the GW-VM on the exposed network, in order to connect the GW-VM
	// with the network to be exposed on the VPN
	// NOTE Select first exposed net
	var exposedNetID string = ""
	for _, network := range resset.Networks {
		log.Trace("Network: with name: ", network.NetworkName)
		if strings.HasPrefix(network.NetworkName, obj.Netconfig.ExposedNetworksNamePrefix) {
			log.Info("Selected Network with name ", network.NetworkName)
			exposedNetID = network.NetworkId
			break
		}

	}
	if exposedNetID == "" {
		log.Error("No network has name matching the ExposedNetworkNamePrefix, so GW-VM will not be connected to the exposed network, continue without port")
		resset.StaticGW.PortID = ""
	} else {
		portId, err := (*vim).CreateInterfacePort(gwInstanceID, exposedNetID)
		if err != nil {
			log.Error("Error adding the interface to GW-VM on the exposed network, continue without port...")
			resset.StaticGW.PortID = ""
		}
		resset.StaticGW.PortID = portId
	}

	ext := ExternalIP{}
	ext.ExternalIp = gwFloatingIP
	ext.PortID = ""
	ext.PortName = gwExternalInterfaceName
	ext.FloatingID = ""
	resset.Gateway.External = ext
	resset.Status = WAIT_FOR_GATEWAY_CONFIG

	result = obj.DB.Save(resset)
	if result.Error != nil {
		log.Error("Impossible to create network resources. Error saving in DB")
		SetErrorResponse(c, http.StatusInternalServerError, ErrGeneral)
		return
	}
	log.Trace("PostNetResourcesFixedSap updated and saved ", *resset)

	SetNetResourcesResponse(c, http.StatusCreated, *resset)

}
