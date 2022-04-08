package nsm

import (
	"net"
	nsmmapi "nextworks/nsm/api"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func checkExcludedSubnetsParams(subs *string) error {
	if _, _, err := net.ParseCIDR(*subs); err != nil {
		return ErrNetResourcesExcludeSubnetsWrongInfo
	}
	return nil
}

// SetNetResourcesListResponse
func SetNetResourcesListResponse(ctx *gin.Context, status int, resources []ResourceSet) {
	output := []nsmmapi.SliceResources{} // initialization to return with marshal []

	for _, resource := range resources {
		var netlist []nsmmapi.Network
		var saplist []nsmmapi.Sap
		for _, net := range resource.Networks {
			netlist = append(netlist, nsmmapi.Network{
				NetworkName: net.NetworkName,
				SubnetCidr:  net.SubnetCidr,
			})
		}
		for _, sap := range resource.Saps {
			saplist = append(saplist, nsmmapi.Sap{
				FloatingNetworkName: sap.FloatingNetName,
				NetworkName:         sap.NetworkName,
				SubnetCidr:          sap.SubnetCidr})
		}
		outResource := nsmmapi.SliceResources{
			Id:                  resource.ID,
			Status:              resource.Status,
			SliceId:             resource.SliceId,
			VimName:             resource.VimName,
			Networks:            netlist,
			ServiceAccessPoints: saplist,
		}
		log.Trace(outResource)
		output = append(output, outResource)
	}

	ctx.JSON(status, output)
}

// SetNetResourcesResponse creates the return type for api
// ResourceSet -> SliceResources
func SetNetResourcesResponse(ctx *gin.Context, status int, res ResourceSet) {

	var netlist []nsmmapi.Network
	var saplist []nsmmapi.Sap

	for _, net := range res.Networks {
		netlist = append(netlist, nsmmapi.Network{
			NetworkName: net.NetworkName,
			SubnetCidr:  net.SubnetCidr,
		})
	}
	for _, sap := range res.Saps {
		saplist = append(saplist, nsmmapi.Sap{
			FloatingNetworkName: sap.FloatingNetName,
			NetworkName:         sap.NetworkName,
			SubnetCidr:          sap.SubnetCidr,
		})
	}

	output := nsmmapi.SliceResources{
		Id:                  res.ID,
		Status:              res.Status,
		SliceId:             res.SliceId,
		VimName:             res.VimName,
		Networks:            netlist,
		ServiceAccessPoints: saplist,
		StaticSap:           res.StaticGW.Enabled,
	}
	ctx.JSON(status, output)
}

// LoadAssociationFromDB load from DB associations of a resource set
// networks and sap
func LoadNetworkAssociationFromDB(database *gorm.DB, netres *ResourceSet) error {
	var nets []Network
	var saps []Sap

	// Retrieve associations - networks and saps
	a := database.Model(&netres).Association("Networks")
	result := a.Find(&nets)
	if result != nil {
		return result
	}
	netres.Networks = nets

	b := database.Model(&netres).Association("Saps")
	result = b.Find(&saps)
	if result != nil {
		return result
	}

	netres.Saps = saps
	return nil
}

// LoadSAPsFromDB load from DB associations of a resource set
// only for SAP
func LoadSAPsFromDB(database *gorm.DB, netres *ResourceSet) error {
	var saps []Sap
	b := database.Model(&netres).Association("Saps")
	result := b.Find(&saps)
	if result != nil {
		return result
	}

	netres.Saps = saps
	return nil
}

// RetrieveResourcesFromDB load from DB a ResourceSet by its ID and all the associations
// networks and sap
func RetrieveResourcesFromDB(database *gorm.DB, id int) (*ResourceSet, error) {
	var result *gorm.DB
	netres := ResourceSet{
		ID: id,
	}

	result = database.First(&netres)
	if result.Error != nil {
		return nil, result.Error
	}
	// Retrieve associations - networks and saps
	err := LoadNetworkAssociationFromDB(database, &netres)
	if err != nil {
		return &netres, err
	}

	return &netres, nil
}

// RetrieveResourcesFromDB load from DB a ResourceSet by its slice-id and all the associations
// networks and sap
func RetrieveResourcesFromDBbySliceID(database *gorm.DB, sliceId string) (*ResourceSet, error) {
	var result *gorm.DB
	var netres ResourceSet

	result = database.First(&netres, "slice_id = ?", sliceId)
	if result.Error != nil {
		return nil, result.Error
	}
	// Retrieve associations - networks and saps
	err := LoadNetworkAssociationFromDB(database, &netres)
	if err != nil {
		return &netres, err
	}

	return &netres, nil
}
