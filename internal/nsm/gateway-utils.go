package nsm

import (
	"net"
	nsmmapi "nextworks/nsm/api"
	"strings"

	"github.com/gin-gonic/gin"
)

func checkGatewayConfigurationParams(input nsmmapi.Gateway) error {
	_, err := parsePort(input.MgmtPort)
	if err != nil {
		return ErrGatewayConfigMgmtPort
	}
	for _, sub := range input.SubnetToExpose {
		if _, _, err = net.ParseCIDR(sub); err != nil {
			return ErrGatewayConfigSubnet
		}
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
	gateway.SubnetToExpose = SubnetsToArray(res.Gateway.ExposedNets)
	gateway.PubKey = res.Gateway.PubKey
	ctx.JSON(status, gateway)
}

func SubnetsToString(subnets []string) string {
	stringSubs := ""
	for i, sub := range subnets {
		if i == 0 {
			stringSubs = stringSubs + sub
		} else {
			stringSubs = stringSubs + "," + sub
		}
	}
	return stringSubs
}

func SubnetsToArray(subnets string) []string {
	return strings.Split(subnets, ",")
}
