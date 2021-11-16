package main

import (
	"nextworks/nsm/internal/nbi"
	"nextworks/nsm/internal/openstackclient"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Global variables
var identityEndpoint string = "http://10.30.7.10:5000/v3"
var username string = "timeo"
var password string = "nextworks"
var tenantID string = "7953babdca974e7ab44cc6c69f093956"
var domainID string = "default"

// TODO: how to manage API tokens and auth
// TODO: thread safe object?

func setupRouter(client *openstackclient.OpenStackClient) *gin.Engine {
	router := gin.Default()

	client.Init()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Build Env oject
	env := new(nbi.Env)
	env.Client = client

	// pre-provisioning routes (related to create required networks)
	router.GET("/network", env.RetrieveNetwork)
	router.POST("/network", env.CreateNetwork)
	router.DELETE("/network", env.DeleteNetwork)
	// provisioning routes (related to gateway connectivity)
	router.GET("gateway/connectivity", env.RetrieveGatewayConnectivity)
	router.POST("gateway/connectivity", env.CreateGatewayConnectivity)
	router.DELETE("gateway/connectivity", env.DeleteGatewayConnectivity)
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

	// Set up environment & object (OpenstackClient)
	openstackclient := openstackclient.NewOpenStackClient(identityEndpoint, username, password, tenantID, domainID)
	log.Info("Tenant" + openstackclient.TenantID)

	// wait SIG TERM
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Received SIG TERM")
		openstackclient.Close()
		os.Exit(1)
	}()

	r := setupRouter(openstackclient)
	// // Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
