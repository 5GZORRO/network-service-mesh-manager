package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Global variables
var IdentityEndpoint string = "http://10.30.7.10:5000/v3"
var TenantID string = "7953babdca974e7ab44cc6c69f093956"
var Username string = "timeo"
var Password string = "nextworks"
var DomainID string = "default"

func setupRouter() *gin.Engine {
	router := gin.Default()

	// routes
	router.GET("/network", retrieveNetwork)
	router.POST("/network", createNetwork)
	router.DELETE("/network", deleteNetwork)

	return router
}

func main() {
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	// log.SetReportCaller(true)  // add function name to logs

	// wait SIG TERM
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Received SIG TERM")
		Close()
		os.Exit(1)
	}()

	Init()

	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")

	Close()
}
