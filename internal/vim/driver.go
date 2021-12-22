package vim

// Vim Types, only 2 vims supported
type VimType string

const (
	Openstack  VimType = "openstack"
	Kubernetes VimType = "kubernetes"
)

type VimDriver interface {
	Authenticate()
	// // Set of methods to prepare the environment before NS instantiation
	// // CreateNetwork() creates a network with a subnet
	CreateNetwork()
	RetrieveNetwork()
	DeleteNetwork()
	// // CreateSAP() creates the infrastructure to have a floating-ip, it could be for a gateway or for other sap of the ns
	CreateSAP()
	DeleteSAP()
	Revoke()
}
