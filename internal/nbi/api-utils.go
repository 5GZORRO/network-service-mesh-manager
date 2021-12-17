package nbi

import (
	"errors"
	"strconv"

	nsmmapi "nextworks/nsm/api"
	"nextworks/nsm/internal/nsm"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SetErrorResponse(ctx *gin.Context, method string, errorStatus int, err error) {
	log.Error(method, " - ", err.Error())
	outputJson := nsmmapi.ErrorResponse{Error: err.Error()}
	ctx.JSON(errorStatus, outputJson)
}

func RetrieveResourcesFromDB(database *gorm.DB, id int) (*nsm.ResourceSet, error) {
	var result *gorm.DB
	var nets []nsm.Network
	var saps []nsm.Sap
	netres := nsm.ResourceSet{
		ID: id,
	}

	result = database.First(&netres)
	if result.Error != nil {
		return nil, result.Error
	}
	// Retrieve associations - networks and saps
	a := database.Model(&netres).Association("Networks")
	_ = a.Find(&nets)
	netres.Networks = nets

	b := database.Model(&netres).Association("Saps")
	_ = b.Find(&saps)

	netres.Saps = saps
	log.Info(netres)
	return &netres, nil
}

func parsePort(port string) (uint16, error) {
	portInt, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return 0, err
	}
	if portInt == 0 {
		return 0, errors.New("0 is not a valid port number")
	}
	result := uint16(portInt)

	return result, nil
}
