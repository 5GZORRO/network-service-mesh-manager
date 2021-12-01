package nsm

import "time"

type Gateway struct {
	ID                 int
	SliceID            string `gorm:"unique;<-:create;"`
	Status             string
	VimName            string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	VimResourceId      int
	Resources          OpenstackResource `gorm:"foreignKey:VimResourceId;references:ID"`
	ExternalIp         string
	ManagementIP       string
	ManagementPort     uint16
	VPNServerPort      uint16
	VPNServerInterface string
}

// Object referring to table where to store all OS gateway resources
type OpenstackResource struct {
	ID              int `gorm:"autoIncrement"`
	NetworkVimID    string
	NetworkVimName  string
	SubnetVimID     string
	SubnetVimName   string
	SubnetCidr      string
	RouterVimID     string
	RouterVimName   string
	RouterVimPortId string
}
