package nsm

import (
	"errors"
)

var (
	// A gateway for slice already exists
	ErrGatewayExists = errors.New("gateway for slice already exists")
	// A gateway for slice already exists
	ErrMissingInfo = errors.New("body request is missing info")
)
