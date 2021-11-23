package config

// Configurations exported
type Configurations struct {
	Server   ServerConfigurations
	Database DatabaseConfigurations
	Vim      VimConfigurations
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Port int
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	Host     string
	Port     string
	DB       string
	Username string
	Password string
}

// VimConfigurations exported
type VimConfigurations struct {
	IdentityEndpoint string
	Username         string
	Password         string
	TenantID         string
	DomainID         string
}
