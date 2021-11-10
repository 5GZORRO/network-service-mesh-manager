package nbi

import "nextworks/nsm/internal/openstackclient"

// Env object should contain the environment used by the functions
// associated to REST API, such as a ConnectionPool to a DB or a
// Provider to OpenStackAPI
type Env struct {
	Client *openstackclient.OpenStackClient
}
