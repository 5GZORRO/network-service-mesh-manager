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
	nsm "nextworks/nsm/internal/nsm"
	"nextworks/nsm/internal/vim"

	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	log "github.com/sirupsen/logrus"
)

func NewGinServer(impl *nsm.ServerInterfaceImpl, port int) *http.Server {
	swagger, err := nsmapi.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// TODO set to production mode
	// gin.SetMode(gin.ReleaseMode)

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

func main() {
	// Config log
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	// log.SetReportCaller(true)
	log.SetLevel(log.ErrorLevel)

	//  Read config file
	configuration := config.ReadConfigFile()
	// log.Info(*configuration)

	// set log level
	level, err := config.LogLevel(configuration)
	if err != nil {
		log.Error(err, ", default log level is ", log.ErrorLevel)
	} else {
		log.SetLevel(level)
	}
	log.Info(configuration)

	// STEP DB
	// Connect to the DB
	dsn := configuration.Database.Username + ":" + configuration.Database.Password + "@tcp(" + configuration.Database.Host + ":" + configuration.Database.Port + ")/" + configuration.Database.DB + "?charset=utf8mb4&parseTime=True&loc=Local"
	log.Trace(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error("Error connecting to the database")
		return
	}
	// init tables
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&nsm.ResourceSet{}, &nsm.Network{}, &nsm.Sap{})
	// log.Trace(db)

	// STEP VIM
	// initizialize a driver for each vim,
	// TODO also reading from DB
	drivers := vim.InizializeVims(db, configuration.Vim)
	log.Trace(drivers)

	// wait SIG TERM
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("Received SIG TERM")
		for _, vim := range drivers.VimList {
			vim.Revoke()
		}
		os.Exit(1)
	}()

	// Create an instance of our server handler object, containing shared info (DB, driver)
	sii := nsm.NewServerInterfaceImpl(db, drivers, &configuration.Networks)
	s := NewGinServer(sii, configuration.Server.Port)

	log.Fatal(s.ListenAndServe())
}
