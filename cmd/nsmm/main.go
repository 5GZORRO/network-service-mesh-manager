package main

import (
	"nextworks/nsm/internal/config"
	"nextworks/nsm/internal/nbi"
	"nextworks/nsm/internal/openstackclient"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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

// TODO: thread safe object?

func setupRouter(env *nbi.Env) *gin.Engine {
	router := gin.Default()

	env.Client.Init()

	// pre-provisioning routes (related to create required networks)
	router.GET("/network", env.RetrieveNetwork)
	router.POST("/network", env.CreateNetwork)
	router.DELETE("/network", env.DeleteNetwork)
	// provisioning routes (related to gateway connectivity)
	router.GET("gateway/connectivity", env.RetrieveGatewayConnectivity)
	router.POST("gateway/connectivity", env.CreateGatewayConnectivity)
	router.DELETE("gateway/connectivity", env.DeleteGatewayConnectivity)
	router.GET("gateway/connectivity/ip", env.RetriveGatewayFloatingIP)
	// procedure to clean up
	router.GET("/management/db", env.GetDB)
	router.POST("/management/db", env.AddDBEntry)
	router.DELETE("/management/db", env.DeleteDBEntry)
	router.GET("/management/clean", env.CleanUp)
	router.POST("/management/clean", env.CleanUpGateway)
	// test
	router.GET("/test", env.Test)

	return router
}

func main() {
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	// log.SetReportCaller(true)  // add function name to logs

	//  Read config file
	config := readConfigFile()
	log.Info(*config)

	// Build Env oject
	env := new(nbi.Env)

	// Initialize Env with OpenstackClient
	log.Info("Vim configs", config.Vim)
	openstackclient := openstackclient.NewOpenStackClient(config.Vim.IdentityEndpoint, config.Vim.Username, config.Vim.Password, config.Vim.TenantID, config.Vim.DomainID)
	env.Client = openstackclient

	// wait SIG TERM
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Received SIG TERM")
		openstackclient.Close()
		os.Exit(1)
	}()

	r := setupRouter(env)
	// // Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
