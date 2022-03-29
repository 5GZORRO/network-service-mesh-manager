package nsm

import (
	"errors"
)

var (
	// GENERAL
	ErrBodyMissingInfo       = errors.New("body request is missing info")
	ErrBodyWrongInfo         = errors.New("body request contains wrong info")
	ErrMissingQueryParameter = errors.New("missing or empty query parameter")
	ErrGeneral               = errors.New("internal error")
	// NET-RESOURCES
	ErrResourcesNotExists                  = errors.New("network resources for requested ID do not exist")
	ErrSliceNotExists                      = errors.New("network resources for requested slice-id do not exist")
	ErrSliceExists                         = errors.New("network resources for slice-id already exist")
	ErrResourcesCantBeDeleted              = errors.New("network resources of the slice cannot be canceled")
	ErrNetResourcesExcludeSubnetsWrongInfo = errors.New("body request contains wrong info in exclude-subnet field")
	// External-IP
	ErrDeleteExternalIP      = errors.New("external-IP of gateway can't be removed. Gateway should be reset first")
	ErrAssociatingExternalIP = errors.New("external-IP cannot be associated to gateway")
	ErrExternalIPExists      = errors.New("an external-IP is already associated to the gateway")
	ErrDeletingExternalIP    = errors.New("error deleting external-IP of gateway")
	ErrGatewayNoNetworkFound = errors.New("no SAP network found")
	ErrNoExternalIP          = errors.New("external-IP of gateway does not exists")
	ErrExternalIPWrongState  = errors.New("external-IP of gateway does not exists. Wrong state")
	// CONFIGURATION ERRORS
	ErrGatewayNotConfigured        = errors.New("gateway is not yet configured")
	ErrRequestConfigurationGateway = errors.New("configuration body request contains wrong info")
	ErrGatewayConfigMgmtIp         = errors.New("configuration gateway body request contains wrong info: management IP")
	ErrConfiguringGateway          = errors.New("gateway can't be configured")
	ErrDeleteConfigurationGateway  = errors.New("configuration of gateway can't be removed")
	ErrUpdatingGatewayInDB         = errors.New("error updating gateway in DB")
	ErrGatewayNotRunning           = errors.New("error gateway is not running")
	// VIM Errors
	ErrVimNotExists       = errors.New("vim does not exist")
	ErrVimCreatingNetwork = errors.New("error creating network on vim")
	ErrVimCreatingSAP     = errors.New("error creating SAP on vim")
	// Connections
	ErrConnectionParameters    = errors.New("connection parameters in body request are wrong")
	ErrNoConnection            = errors.New("no connection found")
	ErrCreatingConnection      = errors.New("error creating VPN connection")
	ErrClosingConnection       = errors.New("error closing a VPN connection")
	ErrSavingConnectionDB      = errors.New("error saving connection in DB")
	ErrConnectionCanBeCanceled = errors.New("error deleting a VPN connection, gateway is not running")
)
