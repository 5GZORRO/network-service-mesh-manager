package main

import (
	"net"
	gatewayconfig "nextworks/nsm/internal/gateway-config"

	log "github.com/sirupsen/logrus"
)

type VpnInfo struct {
	Did       string `json:"did,omitempty"`
	PublicKey string `json:"public_key,omitempty"`
	IPRange   string `json:"IP_range,omitempty"`
	VpnPort   int32  `json:"vpn_port,omitempty"`
}

type PostLaunch struct {
	IpRange      string `json:"ip_range"`
	NetInterface string `json:"net_interface"`
	Port         string `json:"port"`
}

type PostDisconnect struct {
	IpAddressServer string `json:"ip_address_server"`
	PortServer      string `json:"port_server"`
}

type PostConnect struct {
	IpAddressServer   string `json:"ip_address_server"`
	PortServer        string `json:"port_server"`
	IPRangeToRedirect string `json:"IP_range_to_redirect"`
}

type ClientAccepted struct {
	AssignedIp      string `json:"assigned_ip,omitempty"`
	VpnPort         int32  `json:"vpn_port,omitempty"`
	ServerPublicKey string `json:"server_public_key,omitempty"`
}

func main() {
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	log.SetLevel(log.TraceLevel)

	// mng := nsm.NewNetworkManager("192.168.161.16/28", true)
	// log.Info(*mng)
	// nn := mng.NextSubnet()
	// log.Info(*mng)
	// log.Info(nn.String())

	// nn = mng.NextSubnet()
	// log.Info(*mng)
	// log.Info(nn)

	// Test HTTP client

	umuclient := gatewayconfig.New(net.ParseIP("172.22.0.1"), "8082")
	res1 := umuclient.Start("", "", "")
	log.Debug(res1)

	res2 := umuclient.GetCurrentConfiguration()
	log.Debug(res2)

	res3 := umuclient.Connect("", "", "")
	log.Debug(res3)

	res4 := umuclient.Disconnect("", "")
	log.Debug(res4)
}
