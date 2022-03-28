package gatewayconfig

import identityclient "nextworks/nsm/internal/identity"

type VPNHttpClient interface {
	GetCurrentConfiguration() *VpnInfo
	Launch(ipRange string, netInterface string, port string, keyPair *identityclient.KeyPair, idependpoint string) bool
	Connect(peerIp string, peerPort string, remoteIPs string, localIPs string) bool
	Disconnect(peerIP string, peerPort string) bool
}
