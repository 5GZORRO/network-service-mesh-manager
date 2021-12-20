package nsm

import "time"

type ResourceSet struct {
	ID        int
	SliceId   string `gorm:"unique;<-:create;"`
	Status    string
	VimName   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Gateway   Gateway   `gorm:"embedded;embeddedPrefix:gw_"`
	Networks  []Network `gorm:"foreignKey:ResourceSetId;OnDelete:CASCADE"`
	Saps      []Sap     `gorm:"foreignKey:ResourceSetId;OnDelete:CASCADE;"`
}

type Gateway struct {
	MgmtIp      string
	MgmtPort    uint16
	ExternalIp  string
	PubKey      string
	ExposedNets string
}

type Network struct {
	ID            int `gorm:"autoIncrement"`
	ResourceSetId int
	NetworkID     string
	NetworkName   string
	SubnetID      string
	SubnetName    string
	SubnetCidr    string
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
