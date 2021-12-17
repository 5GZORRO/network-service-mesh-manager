package nsm

import "time"

type ResourceSet struct {
	ID        int
	SliceId   string `gorm:"unique;<-:create;"`
	Status    string
	VimName   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Gateway   Gateway `gorm:"embedded;embeddedPrefix:gw_"`
}

type Gateway struct {
	MgmtIp       string
	MgmtPort     uint16
	ExternalIp   string
	ExposedNets  string
	VpnPort      uint16
	VpnInterface string
}

type Network struct {
	ID              int `gorm:"autoIncrement"`
	ResourceSetId   int
	NetworkVimID    string
	NetworkVimName  string
	SubnetVimID     string
	SubnetVimName   string
	SubnetCidr      string
	RouterVimID     string
	RouterVimName   string
	RouterVimPortId string
}

type Sap struct {
	ID              int `gorm:"autoIncrement"`
	ResourceSetId   int
	NetworkId       string
	NetworkName     string
	SubnetId        string
	SubnetName      string
	SubnetCidr      string
	RouterId        string
	RouterName      string
	RouterPortId    string
	FloatingNetId   string
	FloatingNetName string
}
