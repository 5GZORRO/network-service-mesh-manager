package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Configurations exported
type Configurations struct {
	Log      string
	Server   ServerConfigurations
	Database DatabaseConfigurations
	Vim      []VimConfigurations
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
	Name             string
	Type             string
	IdentityEndpoint string
	Username         string
	Password         string
	TenantID         string
	DomainID         string
}

func LogLevel(c *Configurations) (log.Level, error) {
	switch c.Log {
	case "Trace":
		return log.TraceLevel, nil
	case "Debug":
		return log.DebugLevel, nil
	case "Info":
		return log.InfoLevel, nil
	case "Error":
		return log.ErrorLevel, nil
	case "Panic":
		return log.PanicLevel, nil
	case "Fatal":
		return log.FatalLevel, nil
	default:
		return 0, ErrLogLevel
	}
}

func ReadConfigFile() *Configurations {
	var config Configurations

	// Set the file name of the configurations file, the path and the type file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	// Set default values
	viper.SetDefault("server.port", 8080)

	// Read and initialize
	if err := viper.ReadInConfig(); err != nil {
		log.Error("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Error("Unable to decode into struct, %v", err)
	}

	return &config
}

func CheckVimParams(vimConfig VimConfigurations) error {
	if vimConfig.Name == "" {
		return ErrMissingVimName
	}
	if vimConfig.Type == "" {
		return ErrMissingVimType
	}
	if vimConfig.IdentityEndpoint == "" {
		return ErrMissingVimEndpoint
	}
	if vimConfig.DomainID == "" {
		return ErrMissingVimDomain
	}
	if vimConfig.TenantID == "" {
		return ErrMissingVimTenant
	}
	if vimConfig.Password == "" {
		return ErrMissingVimPassoword
	}
	if vimConfig.Username == "" {
		return ErrMissingVimUsername
	}
	return nil
}
