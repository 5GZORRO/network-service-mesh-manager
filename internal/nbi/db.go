package nbi

// t represent all connectivity information for an interdomain
// communication:
// for each slice, I assume a
// priv_net with a subnet
// router connected to the floating net
// a gateway vm
// a floating ip
type GatewayConnectivity struct {
	SliceID     string `json:"sliceID" binding:"required"`
	PrivNetID   string `json:"networkID" binding:"required"`
	SubnetID    string `json:"subnetID" binding:"required"`
	RouterID    string `json:"routerID" binding:"required"`
	InterfaceID string `json:"interfaceID"`
}
