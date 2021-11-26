package nsm

import (
	"errors"
)

var (
	// A gateway for slice already exists
	ErrGatewayExists = errors.New("gateway for slice already exists")
	//
	ErrElementAlreadyExists = errors.New("gateway already exists")
	//
	ErrGatewayNotFound = errors.New("gateway does not exist")
	//
	ErrGatewayCantBeDeleted = errors.New("gateway can't be deleted")
	//
	ErrBodyMissingInfo = errors.New("body request is missing info")
	//
	ErrBodyWrongInfo = errors.New("body request contains wrong info")
	//
	ErrGatewayConfiguration = errors.New("configuration body request contains wrong info")
	//
	ErrConfiguringGateway = errors.New("gateway can't be configured")
)
