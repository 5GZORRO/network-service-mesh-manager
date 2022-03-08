package gatewayconfig

type VPNHttpClient interface {
	GetCurrentConfiguration() *VpnInfo
	Launch(ipRange string, netInterface string, port string) bool
	Connect(peerIp string, peerPort string, remoteIPs string, localIPs string) bool
	Disconnect(peerIP string, peerPort string) bool
}
