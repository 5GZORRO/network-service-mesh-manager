package nsm

import "time"

type Gateway struct {
	ID             int
	SliceID        string `gorm:"unique;<-:create"`
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	NetworkID      string
	SubnetID       string
	RouterID       string
	ExternalIp     string
	ManagementIP   string
	ManagementPort uint16
	VPNPort        string
	VPNInterface   string
}

type Network struct {
	ID             int
	VimNetworkID   string
	VimNetworkName string
}

type Subnet struct {
	ID            int
	VimSubnetID   string
	VimSubnetName string
	SubentCidr    string
}

type Router struct {
	ID              int
	VimRouterID     string
	VimRouterName   string
	VimRouterPortId string
}
