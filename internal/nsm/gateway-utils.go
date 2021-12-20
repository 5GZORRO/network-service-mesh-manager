package nsm

import (
	"net"
	nsmmapi "nextworks/nsm/api"

	"github.com/gin-gonic/gin"
)

func checkGatewayConfigurationParams(input nsmmapi.Gateway) error {
	_, err := parsePort(input.MgmtPort)
	if err != nil {
		return ErrGatewayConfigMgmtPort
	}
	if _, _, err = net.ParseCIDR(input.SubnetToExpose); err != nil {
		return ErrGatewayConfigSubnet
	}
	mngmIp := net.ParseIP(input.MgmtIp)
	externalIP := net.ParseIP(input.ExternalIp)
	if mngmIp == nil {
		return ErrGatewayConfigMgmtIp
	}
	if externalIP == nil {
		return ErrGatewayConfigExternalIp
	}
	return nil
}

// SetNetResourcesResponse creates the return type for api
//
func SetGatewayResponse(ctx *gin.Context, status int, res ResourceSet) {
	var gateway nsmmapi.Gateway
	gateway.ExternalIp = res.Gateway.ExternalIp
	gateway.MgmtIp = res.Gateway.MgmtIp
	gateway.MgmtPort = parsePortToString(res.Gateway.MgmtPort)
	// Todo publickKey of the server
	// gateway.PubKey = res.Gateway.
	gateway.SubnetToExpose = res.Gateway.ExposedNets
	gateway.PubKey = res.Gateway.PubKey
	ctx.JSON(status, gateway)
}
