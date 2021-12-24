package gatewayconfig

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