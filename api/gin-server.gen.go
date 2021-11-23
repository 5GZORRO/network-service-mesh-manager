// Package NsmmApi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version (devel) DO NOT EDIT.
package NsmmApi

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/gin-gonic/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (DELETE /gateway)
	DeleteGateway(c *gin.Context, params DeleteGatewayParams)

	// (GET /gateway)
	GetGateway(c *gin.Context, params GetGatewayParams)

	// (POST /gateway)
	PostGateway(c *gin.Context, params PostGatewayParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
}

type MiddlewareFunc func(c *gin.Context)

// DeleteGateway operation middleware
func (siw *ServerInterfaceWrapper) DeleteGateway(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params DeleteGatewayParams

	// ------------- Required query parameter "sliceId" -------------
	if paramValue := c.Query("sliceId"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument sliceId is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "sliceId", c.Request.URL.Query(), &params.SliceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter sliceId: %s", err)})
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.DeleteGateway(c, params)
}

// GetGateway operation middleware
func (siw *ServerInterfaceWrapper) GetGateway(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetGatewayParams

	// ------------- Required query parameter "sliceId" -------------
	if paramValue := c.Query("sliceId"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument sliceId is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "sliceId", c.Request.URL.Query(), &params.SliceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter sliceId: %s", err)})
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetGateway(c, params)
}

// PostGateway operation middleware
func (siw *ServerInterfaceWrapper) PostGateway(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostGatewayParams

	// ------------- Required query parameter "sliceId" -------------
	if paramValue := c.Query("sliceId"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument sliceId is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "sliceId", c.Request.URL.Query(), &params.SliceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter sliceId: %s", err)})
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostGateway(c, params)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL     string
	Middlewares []MiddlewareFunc
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router *gin.Engine, si ServerInterface) *gin.Engine {
	return RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router *gin.Engine, si ServerInterface, options GinServerOptions) *gin.Engine {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
	}

	router.DELETE(options.BaseURL+"/gateway", wrapper.DeleteGateway)

	router.GET(options.BaseURL+"/gateway", wrapper.GetGateway)

	router.POST(options.BaseURL+"/gateway", wrapper.PostGateway)

	return router
}
