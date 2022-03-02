package config

import (
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Configurations exported
type Configurations struct {
	Log      string
	Server   ServerConfigurations
	Database DatabaseConfigurations
	Networks NetworkConfigurations
	Vim      []VimConfigurations
	Vpnaas   VpnaasConfigurations
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

// NetworkConfigurations exported
type NetworkConfigurations struct {
	Start string
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
	FloatingNetwork  string
}

// VpnaasConfigurations exported
type VpnaasConfigurations struct {
	VpnaasPort int
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

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func ReadConfigFile(configFileName string) *Configurations {
	var config Configurations

	log.Trace("Config file name: ", configFileName)
	base := path.Base(configFileName)
	dir := path.Dir(configFileName)
	ext := path.Ext(configFileName)
	log.Debug("Selected config file name: ", base, ", path: ", dir, ", extension: ", ext)
	ext = ext[1:]

	// Set the file name of the configurations file, the path and the type file
	viper.SetConfigName(fileNameWithoutExtension(base))
	viper.AddConfigPath(dir)
	viper.SetConfigType(ext)

	// Set default values
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("networks.start", "192.168.161.0/28")
	viper.SetDefault("vpnaas.port", 8181)

	// Read and initialize
	if err := viper.ReadInConfig(); err != nil {
		log.Error(err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Error(err)
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
	if vimConfig.FloatingNetwork == "" {
		return ErrNoFloatingDefined
	}
	return nil
}
