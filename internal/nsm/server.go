// Package to implement all of the handlers in the ServerInterface
package nsm

import (
	"nextworks/nsm/internal/config"
	"nextworks/nsm/internal/vim"

	"gorm.io/gorm"
)

// Shared object between different HTTP REST handlers
// it should contain the DBconnection
type ServerInterfaceImpl struct {
	DB           *gorm.DB
	Netconfig    *config.NetworkConfigurations
	VpnaasConfig *config.VpnaasConfigurations
	Vims         *vim.VimDriverList
	// Lock     sync.Mutex
}

func NewServerInterfaceImpl(DBconnection *gorm.DB, drivers *vim.VimDriverList,
	net *config.NetworkConfigurations, vpaas *config.VpnaasConfigurations) *ServerInterfaceImpl {
	return &ServerInterfaceImpl{
		DB:           DBconnection,
		Netconfig:    net,
		VpnaasConfig: vpaas,
		Vims:         drivers,
	}
}
