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
	// CONFIGURATION ERRORS
	ErrGatewayNotConfigured         = errors.New("gateway is not yet configured")
	ErrRequestConfigurationGateway  = errors.New("configuration body request contains wrong info")
	ErrGatewayConfigMgmtPort        = errors.New("configuration gateway body request contains wrong info: management port")
	ErrGatewayConfigMgmtIp          = errors.New("configuration gateway body request contains wrong info: management IP")
	ErrGatewayConfigExternalIp      = errors.New("configuration gateway body request contains wrong info: external IP")
	ErrGatewayConfigSubnet          = errors.New("configuration gateway body request contains wrong info: subnet to expose")
	ErrGatewayVpnPrivateRange       = errors.New("configuration gateway body request contains wrong info: VPN private range")
	ErrGatewayVpnPeerPrivateIp      = errors.New("configuration gateway body request contains wrong info: VPN peer private ip")
	ErrGatewayVpnPeerPrivateIpRange = errors.New("configuration gateway body request contains wrong info: VPN peer private ip does not belong to private-vpn-range")
	ErrConfiguringGateway           = errors.New("gateway can't be configured")
	ErrDeleteConfigurationGateway   = errors.New("configuration of gateway can't be removed")
	ErrUpdatingGatewayInDB          = errors.New("error updating gateway in DB")
	ErrGatewayNotRunning            = errors.New("error gateway is not running")
	// VIM Errors
	ErrVimNotExists       = errors.New("vim does not exist")
	ErrVimCreatingNetwork = errors.New("error creating network on vim")
	ErrVimCreatingSAP     = errors.New("error creating SAP on vim")
	// Connections
	ErrConnectionParameters    = errors.New("connection parameters in body request are wrong")
	ErrNoConnection            = errors.New("no connection found")
	ErrCreatingConnection      = errors.New("error creating VPN connection")
	ErrSavingConnectionDB      = errors.New("error saving connection in DB")
	ErrConnectionCanBeCanceled = errors.New("error deleting a VPN connection, gateway is not running")
)
