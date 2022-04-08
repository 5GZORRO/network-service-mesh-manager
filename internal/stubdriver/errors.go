package stubdriver

import (
	"errors"
)

var (
	//
	ErrVimNotFound           = errors.New("Selected VIM does not exists")
	ErrStaticGatewayNotFound = errors.New("static gateway for selected VIM does not exist")
)
