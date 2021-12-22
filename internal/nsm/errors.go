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
	ErrGatewayNotConfigured        = errors.New("gateway is not yet configured")
	ErrRequestConfigurationGateway = errors.New("configuration body request contains wrong info")
	ErrGatewayConfigMgmtPort       = errors.New("configuration gateway body request contains wrong info: management port")
	ErrGatewayConfigMgmtIp         = errors.New("configuration gateway body request contains wrong info: management IP")
	ErrGatewayConfigExternalIp     = errors.New("configuration gateway body request contains wrong info: external IP")
	ErrGatewayConfigSubnet         = errors.New("configuration gateway body request contains wrong info: subnet to expose")
	ErrConfiguringGateway          = errors.New("gateway can't be configured")
	ErrDeleteConfigurationGateway  = errors.New("configuration of gateway can't be removed")
	// VIM Errors
	ErrVimNotExists = errors.New("vim does not exist")
)
