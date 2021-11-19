// Package NsmmApi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version (devel) DO NOT EDIT.
package NsmmApi

// Subnet of the private net
type PostGatewayBody struct {
	Subnet string `json:"subnet"`
}

// DeleteGatewayParams defines parameters for DeleteGateway.
type DeleteGatewayParams struct {
	// Slice unique identifier
	SliceId string `json:"sliceId"`
}

// GetGatewayParams defines parameters for GetGateway.
type GetGatewayParams struct {
	// Slice unique identifier
	SliceId string `json:"sliceId"`
}

// PostGatewayJSONBody defines parameters for PostGateway.
type PostGatewayJSONBody PostGatewayBody

// PostGatewayParams defines parameters for PostGateway.
type PostGatewayParams struct {
	// Slice unique identifier
	SliceId string `json:"sliceId"`
}

// PostGatewayJSONRequestBody defines body for PostGateway for application/json ContentType.
type PostGatewayJSONRequestBody PostGatewayJSONBody

