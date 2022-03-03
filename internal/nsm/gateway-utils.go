package nsm

import (
	"net"
	nsmmapi "nextworks/nsm/api"

	"github.com/gin-gonic/gin"
)

func checkGatewayConfigurationParams(input nsmmapi.PostGateway) error {
	for _, sub := range input.SubnetToExpose {
		if _, _, err := net.ParseCIDR(sub); err != nil {
			return ErrGatewayConfigSubnet
		}
	}
	mngmIp := net.ParseIP(input.MgmtIp)
	if mngmIp == nil {
		return ErrGatewayConfigMgmtIp
	}
	_, privateVpnNet, err := net.ParseCIDR(input.PrivateVpnRange)
	if err != nil {
		return ErrGatewayVpnPrivateRange
	}
	if input.PrivateVpnPeerIp != nil {
		privateVpnPeerIp := net.ParseIP(*input.PrivateVpnPeerIp)
		if privateVpnPeerIp == nil {
			return ErrGatewayVpnPeerPrivateIp
		}
		// check if privateVpnPeerIp belongs to privateVpnNet
		if !privateVpnNet.Contains(privateVpnPeerIp) {
			return ErrGatewayVpnPeerPrivateIpRange
		}
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
	gateway.PrivateVpnIp = res.Gateway.PrivateVpnIp
	gateway.PrivateVpnRange = res.Gateway.PrivateVpnRange
	ctx.JSON(status, gateway)
}
