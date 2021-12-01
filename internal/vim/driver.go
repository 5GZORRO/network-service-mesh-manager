package vim

// Vim Types, only 2 vims supported
type VimType string

const (
	Openstack  VimType = "openstack"
	Kubernetes VimType = "kubernetes"
)

type VimDriver interface {
	Authenticate()
	CreateGatewayConnectivity(sliceId string, subnet string) (string, string, string, error)
	DeleteGatewayConnectivity(networkId string, subnetId string, routerId string) error
	GetGatewayConnectivity(networkId string, subnetId string, routerId string)
	Revoke()
}
