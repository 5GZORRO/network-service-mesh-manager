package vim

// Vim Types, only 2 vims supported
type VimType string

const (
	Openstack  VimType = "openstack"
	Kubernetes VimType = "kubernetes"
	None       VimType = "none"
)

type VimDriver interface {
	// GetStaticGatewayInfo to retrieve from the actual implementation object the information of the static gw for this VIM
	GetStaticGatewayInfo() (string, string, string, string, string, error)
	Authenticate()
	// AllocateFloatingIP allocates a floatingIp identified to a (compute) port found on the SAP network, selected using the prefix network
	// it returns the portID, the interface name associated, the floatingID, the floating ip-address
	AllocateFloatingIP(networkid string) (string, string, string, string, error)
	DeallocateFloatingIP(portID string, fipID string) error
	// CreateNetwork() creates a network with a subnet
	// It returns:
	// networkID, subnetID, subnetName
	CreateNetwork(name string, cidr string, gateway bool) (string, string, string, error)
	//
	RetrieveNetwork(id string)
	DeleteNetwork(networkID string, subnetID string) error
	// CreateSAP() creates the infrastructure to have a floating-ip, it could be for a gateway or for other sap of the ns
	// It returns:
	// networkID, subnetID, subnetName, routerID, routerName, portID
	CreateSAP(floatingNetName string, networkName string, cidr string) (string, string, string, string, string, string, error)
	DeleteSAP(networkID string, subnetID string, routerID string, portID string) error
	// Function to return the FloatingNetworkName, set at init in the VimDriver object
	RetrieveFloatingNetworkName() string
	RetrieveFloatingNetworkID() string
	// Attach a VM to a Network (used for now in case of static GW to attach to the VM to the exposed network)
	CreateInterfacePort(string, string) (string, error)
	DeleteInterfacePort(string, string) error
	Revoke()
}
