package nbi

import (
	"errors"
	NsmmApi "nextworks/nsm/api"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func SetErrorResponse(ctx *gin.Context, method string, errorStatus int, err error) {
	log.Error(method, " - ", err.Error())
	outputJson := NsmmApi.ErrorResponse{Error: err.Error()}
	ctx.JSON(errorStatus, outputJson)
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
