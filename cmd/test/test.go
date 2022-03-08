package main

// Executable to test the VPNaaS Client and the interaction with the VPNaaS module (UMU)

import (
	"net"
	gatewayconfig "nextworks/nsm/internal/gateway-config"
	"time"

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

	gw1ip := "10.30.8.84"
	gw2ip := "10.30.8.129"
	gw1port := "8080"
	gw2port := "8080"

	umuclient1 := gatewayconfig.New(net.ParseIP(gw1ip), gw1port, "local")
	umuclient2 := gatewayconfig.New(net.ParseIP(gw2ip), gw2port, "local")
	res1 := umuclient1.Launch("192.168.1.1/24", "ens4", gw1port)
	log.Debug(res1)
	res2 := umuclient2.Launch("192.168.2.1/24", "ens4", gw2port)
	log.Debug(res2)

	// res2 := umuclient.GetCurrentConfiguration()
	// log.Debug(res2)

	res3 := umuclient1.Connect(gw2ip, gw2port, "192.168.2.1/32, 192.168.162.0/24", "192.168.161.0/24")
	log.Debug(res3)

	time.Sleep(10 * time.Second)

	res4 := umuclient1.Disconnect(gw2ip, gw2port)
	log.Debug(res4)
}
