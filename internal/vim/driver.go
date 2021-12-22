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
	// // Set of methods to prepare the environment before NS instantiation
	// // CreateNetwork() creates a network with a subnet
	CreateNetwork(name string, cidr string) (string, string, string, error)
	RetrieveNetwork(id string)
	DeleteNetwork(id string)
	// // CreateSAP() creates the infrastructure to have a floating-ip, it could be for a gateway or for other sap of the ns
	CreateSAP()
	DeleteSAP()
	Revoke()
}
