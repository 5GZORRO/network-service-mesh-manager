package config

import (
	"errors"
)

var (
	//
	ErrMissingVimName = errors.New("missing Name in config file for vim")
	//
	ErrMissingVimType = errors.New("missing Type in config file for vim")
	//
	ErrMissingVimEndpoint = errors.New("missing Endpoint in config file for vim")
	//
	ErrMissingVimUsername = errors.New("missing Username in config file for vim")
	//
	ErrMissingVimPassoword = errors.New("missing Password in config file for vim")
	//
	ErrMissingVimTenant = errors.New("missing TenantID in config file for vim")
	//
	ErrMissingVimDomain = errors.New("missing DomainID in config file for vim")
	//
	ErrWrongVimType = errors.New("vim type not found")
)
