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
	SliceID     string `gorm:"primaryKey;<-:create"`
	Status      string
	NetworkID   string
	SubnetID    string
	RouterID    string
	InterfaceID string
	FloatingIP  string
	VmGatewayID string
}

type Network struct {
	NetworkID   string `gorm:"primaryKey;<-:create"`
	NetworkName string
}

type Subnet struct {
	SubnetID   string `gorm:"primaryKey;<-:create"`
	SubnetName string
	SubnetCidr string
}

type Router struct {
	RouterID   string `gorm:"primaryKey;<-:create"`
	RouterName string
}
