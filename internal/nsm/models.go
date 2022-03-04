package nsm

import "time"

type ResourceSet struct {
	ID          int
	SliceId     string `gorm:"unique;<-:create;"`
	Status      string
	VimName     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Gateway     Gateway      `gorm:"embedded;embeddedPrefix:gw_"`
	Networks    []Network    `gorm:"foreignKey:ResourceSetId;constraint:OnDelete:CASCADE"`
	Saps        []Sap        `gorm:"foreignKey:ResourceSetId;constraint:OnDelete:CASCADE;"`
	Connections []Connection `gorm:"foreignKey:ResourceSetId;constraint:OnDelete:CASCADE;"`
}

type Gateway struct {
	External ExternalIP `gorm:"embedded"`
	Config   Config     `gorm:"embedded"`
}

type ExternalIP struct {
	ExternalIp string
	PortID     string
	PortName   string
	FloatingID string
}

type Config struct {
	MgmtIp          string
	MgmtPort        uint16 // NOTE: It is the Server and wireguard port!!
	PrivateVpnRange string
	ExposedNets     string
}

type Network struct {
	ID            int `gorm:"autoIncrement"`
	ResourceSetId int
	NetworkId     string
	NetworkName   string
	SubnetId      string
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
	FloatingNetID   string
	FloatingNetName string
}

type Connection struct {
	ID            int `gorm:"autoIncrement"`
	ResourceSetId int
	PublicKey     string
	PrivateKey    string
	PeerIp        string
	PeerPort      string
	PeerNets      string
}
