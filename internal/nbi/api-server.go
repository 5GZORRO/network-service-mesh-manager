// Package to implement all of the handlers in the ServerInterface
package nbi

import (
	"nextworks/nsm/internal/drivers"

	"gorm.io/gorm"
)

// Shared object between different HTTP REST handlers
// it should contain the DBconnection
type ServerInterfaceImpl struct {
	DB *gorm.DB
	// Config file info
	Vim drivers.VimDriver
	// Lock     sync.Mutex
}

func NewServerInterfaceImpl(DBconnection *gorm.DB, openstackclient *drivers.OpenStackDriver) *ServerInterfaceImpl {
	return &ServerInterfaceImpl{
		DB:  DBconnection,
		Vim: openstackclient,
	}
}
