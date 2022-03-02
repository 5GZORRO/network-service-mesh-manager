package vim

// Vim Types, only 2 vims supported
type VimType string

const (
	Openstack  VimType = "openstack"
	Kubernetes VimType = "kubernetes"
	None       VimType = "none"
)

type VimDriver interface {
	Authenticate()
	// CreateNetwork() creates a network with a subnet
	// It returns:
	// networkID, subnetID, subnetName
	CreateNetwork(name string, cidr string, gateway bool) (string, string, string, error)
	DeleteNetwork(networkID string, subnetID string) error
	// CreateSAP() creates the infrastructure to have a floating-ip, it could be for a gateway or for other sap of the ns
	// It returns:
	// networkID, subnetID, subnetName, routerID, routerName, portID
	CreateSAP(floatingNetName string, networkName string, cidr string) (string, string, string, string, string, string, error)
	DeleteSAP(networkID string, subnetID string, routerID string, portID string) error
	//
	RetrieveNetwork(id string)
	// Function to return the FloatingNetworkName, set at init in the VimDriver object
	RetrieveFloatingNetworkName() string
	Revoke()
}
