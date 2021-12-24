package openstackdriver

import (
	"errors"
)

var (
	// Creation
	ErrFloatingNotFound = errors.New("floating network name not found")
	ErrNetworkCreation  = errors.New("error creating network")
	ErrRouterCreation   = errors.New("error creating router")
	ErrPortCreation     = errors.New("error creating router interface")
	// Delete
	ErrPortRemoval    = errors.New("error deleting port of router")
	ErrRouterRemoval  = errors.New("error deleting router")
	ErrNetworkRemoval = errors.New("error deleting network")
)
