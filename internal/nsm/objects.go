package nsm

// States
var CREATING = "CREATING"
var WAIT_FOR_GATEWAY = "WAIT_FOR_GATEWAY_VM"
var CONFIGURING_GATEWAY = "CONFIGURING_GATEWAY"
var RUNNING = "RUNNING"

// Error state
var CREATION_ERROR = "CREATION_ERROR"
var ASSOCIATION_ERROR = "ASSOCIATION_ERROR"
var CONFIGURATION_ERROR = "CONFIGURATION_ERROR"

type Gateway struct {
	SliceID        string `gorm:"primaryKey;<-:create"`
	Status         string
	NetworkID      string
	NetworkName    string
	SubnetID       string
	SubnetName     string
	SubnetCidr     string
	RouterID       string
	RouterName     string
	PortID         string
	FloatingIP     string
	VmGatewayID    string
	GatewayRole    string
	RemoteEndpoint string
}
