package nbi

// t represent all connectivity information for an interdomain
// communication:
// for each slice, I assume a
// priv_net with a subnet
// router connected to the floating net
// a gateway vm
// a floating ip
type GatewayConnectivity struct {
	SliceID     string
	PrivNetID   string
	SubnetID    string
	RouterID    string
	InterfaceID string
}
