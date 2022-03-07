package nsm

import (
	"net"
	nsmmapi "nextworks/nsm/api"

	"github.com/gin-gonic/gin"
)

func checkGatewayConfigurationParams(input nsmmapi.PostGateway) error {
	mngmIp := net.ParseIP(input.MgmtIp)
	if mngmIp == nil {
		return ErrGatewayConfigMgmtIp
	}
	return nil
}

// SetNetResourcesResponse creates the return type for api
//
func SetGatewayResponse(ctx *gin.Context, status int, res ResourceSet) {
	var gateway nsmmapi.Gateway
	gateway.MgmtIp = res.Gateway.Config.MgmtIp
	gateway.MgmtPort = parsePortToString(res.Gateway.Config.MgmtPort)
	gateway.ExposedSubnets = SubnetsToArray(res.Gateway.Config.ExposedNets)
	gateway.PrivateVpnRange = res.Gateway.Config.PrivateVpnRange
	ctx.JSON(status, gateway)
}
