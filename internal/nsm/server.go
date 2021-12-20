// Package to implement all of the handlers in the ServerInterface
package nsm

import (
	"nextworks/nsm/internal/vim"

	"gorm.io/gorm"
)

// Shared object between different HTTP REST handlers
// it should contain the DBconnection
type ServerInterfaceImpl struct {
	DB *gorm.DB
	// Config file info
	// TODO handle more than one VIMs
	Vims *vim.VimDriverList
	// Lock     sync.Mutex
}

func NewServerInterfaceImpl(DBconnection *gorm.DB, drivers *vim.VimDriverList) *ServerInterfaceImpl {
	return &ServerInterfaceImpl{
		DB:   DBconnection,
		Vims: drivers,
	}
}
