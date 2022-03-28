package nsm

import "time"

type ResourceSet struct {
	ID          int
	SliceId     string `gorm:"unique;<-:create;"`
	Status      string
	VimName     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	StaticSap   bool         `gorm:"default:false"` // indicates if the SAP GW is a PNF
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
	ExternalIp string // Floating IP associated to the GW VM
	PortID     string // Interface ID in OpenStack of the interface associated with the floating IP
	PortName   string // Interface name in the VM of the interface associated with the floating IP, ex. ens3, ens4
	FloatingID string // Floating IP ID in OpenStack
}

type Config struct {
	MgmtIp          string // Management IP of the GW
	MgmtPort        uint16 // NOTE: It is the Server and wireguard port!!
	PrivateVpnRange string
	ExposedNets     string  // Subnets to be exposed through the VPN connection
	Keys            KeyPair `gorm:"embedded;embeddedPrefix:key_"`
}

type KeyPair struct {
	Did       string
	PubK      string
	PrivK     string
	Timestamp string
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
	PeerIp        string
	PeerPort      string
	PeerNets      string
}
