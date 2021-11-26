package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	nsmapi "nextworks/nsm/api"
	"nextworks/nsm/internal/config"
	"nextworks/nsm/internal/drivers"
	"nextworks/nsm/internal/nbi"

	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewGinServer(impl *nbi.ServerInterfaceImpl, port int) *http.Server {
	swagger, err := nsmapi.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// This is how you set up a basic chi router
	r := gin.Default()

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidator(swagger))

	// We now register our petStore above as the handler for the interface
	r = nsmapi.RegisterHandlers(r, impl)

	log.Info(r.Routes())

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
	}
	return s
}

func readConfigFile() *config.Configurations {
	var config config.Configurations

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

func main() {
	// Config log
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	//  Read config file
	configuration := readConfigFile()
	log.Info(*configuration)

	// Connect to the DB
	dsn := configuration.Database.Username + ":" + configuration.Database.Password + "@tcp(" + configuration.Database.Host + ":" + configuration.Database.Port + ")/" + configuration.Database.DB + "?charset=utf8mb4&parseTime=True&loc=Local"
	log.Info(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error("Error connecting to the database")
		return
	}
	log.Info(db)

	// create driver for OS
	// TODO, should handle different vim
	openstackclient := drivers.NewOpenStackDriver(configuration.Vim.IdentityEndpoint, configuration.Vim.Username, configuration.Vim.Password, configuration.Vim.TenantID, configuration.Vim.DomainID)

	// wait SIG TERM
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Received SIG TERM")
		openstackclient.Revoke()
		os.Exit(1)
	}()

	// Create an instance of our handler object, containing shared info (DB, driver)
	sii := nbi.NewServerInterfaceImpl(db, openstackclient)
	s := NewGinServer(sii, configuration.Server.Port)

	log.Fatal(s.ListenAndServe())
}
