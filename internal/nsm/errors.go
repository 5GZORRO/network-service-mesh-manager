package nsm

import (
	"errors"
)

var (
	// NET-RESOURCES
	//
	ErrSliceExists = errors.New("network resources for slice-id already exist")
	//
	ErrSliceNotExists = errors.New("network resources for slice-id do not exist")
	//
	ErrGatewayNotFound = errors.New("gateway does not exist")
	//
	ErrResourcesCantBeDeleted = errors.New("network resources of the slice cannot be canceled")
	//
	ErrBodyMissingInfo = errors.New("body request is missing info")
	//
	ErrBodyWrongInfo = errors.New("body request contains wrong info")
	// CONFIGURATION ERRORS
	//
	ErrRequestConfigurationGateway = errors.New("configuration body request contains wrong info")
	//
	ErrConfiguringGateway = errors.New("gateway can't be configured")
	//
	ErrDeleteConfigurationGateway = errors.New("configuration of gateway can't be removed")
	//
	ErrConfigurationGatewayNotExists = errors.New("gateway is not yet configured")
	// GENERAL -
	//
	ErrMissingQueryParameter = errors.New("missing or empty query parameter")
	//
	ErrGeneral = errors.New("internal error")
)
