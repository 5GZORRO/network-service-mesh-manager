package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	nsmapi "nextworks/nsm/api"

	log "github.com/sirupsen/logrus"

	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
)

func NewGinServer(petStore *nsmapi.ServerInterfaceImpl, port int) *http.Server {
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
	r = nsmapi.RegisterHandlers(r, petStore)

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

	// Connect to the DB
	// TODO read them from a config file
	dsn := "root:root@tcp(127.0.0.1:3306)/nsmm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error("Error connecting to the database")
		return
	}

	var port = flag.Int("port", 8080, "Port for test HTTP server")
	flag.Parse()
	// Create an instance of our handler which satisfies the generated interface
	sii := nsmapi.NewServerInterfaceImpl(db)
	s := NewGinServer(sii, *port)
	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
