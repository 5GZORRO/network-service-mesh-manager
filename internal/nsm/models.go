package nsm

import "time"

type Gateway struct {
	ID                 int
	SliceID            string `gorm:"unique;<-:create"`
	Status             string
	VimName            string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	VimResourceId      int
	ExternalIp         string
	ManagementIP       string
	ManagementPort     uint16
	VPNServerPort      uint16
	VPNServerInterface string
}

// Object referring to table where to store all OS gateway resources
type OpenstackResource struct {
	ID              int
	NetworkVimID    string
	NetworkVimName  string
	SubnetVimID     string
	SubnetVimName   string
	SubentCidr      string
	RouterVimID     string
	RouterVimName   string
	RouterVimPortId string
}
