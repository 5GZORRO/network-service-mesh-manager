// Package Nsmm provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version (devel) DO NOT EDIT.
package Nsmm

// Connection defines model for Connection.
type Connection struct {
	// Unique identifier of the connection-id
	Id int `json:"id"`

	// Public key of the remote peer
	PubKey string `json:"pub-key"`

	// Public IP of the remote peer VPN
	RemotePeerIp string `json:"remote-peer-ip"`

	// Remote peer VPN port
	RemotePeerPort string `json:"remote-peer-port"`

	// Subnet to expose
	SubnetsToExpose []string `json:"subnets-to-expose"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Actual gateway configuration of a VPN server gateway
type Gateway struct {
	// External IP of the gateway
	ExternalIp string `json:"external-ip"`

	// Gateway VM management IP
	MgmtIp string `json:"mgmt-ip"`

	// Gateway VM management port
	MgmtPort string `json:"mgmt-port"`

	// Private VPN IP of the server
	PrivateVpnIp string `json:"private-vpn-ip"`

	// Private VPN subnet
	PrivateVpnRange string `json:"private-vpn-range"`

	// Public key of the server
	PubKey string `json:"pub-key"`

	// Subnet to expose
	SubnetToExpose []string `json:"subnet-to-expose"`
}

// Network defines model for Network.
type Network struct {
	// Network name specified in the NSD
	NetworkName string `json:"network-name"`

	// Network address used for this subnet
	SubnetCidr string `json:"subnet-cidr"`
}

// Body of the POST to create a new VPN connection for a gateway. It contains the information of the remote peer (pub-key, ip, port) and the networks to expose
type PostConnection struct {
	// Subnets exposed by VPN server
	ExposedSubnets []string `json:"exposed-subnets"`

	// Public key of the remote peer
	PubKey string `json:"pub-key"`

	// Public IP of the remote peer VPN
	RemotePeerIp string `json:"remote-peer-ip"`

	// Remote peer VPN port
	RemotePeerPort string `json:"remote-peer-port"`
}

// Configuration of a gateway
type PostGateway struct {
	// External IP of the gateway
	ExternalIp string `json:"external-ip"`

	// Gateway VM management IP
	MgmtIp string `json:"mgmt-ip"`

	// Gateway VM management port
	MgmtPort string `json:"mgmt-port"`

	// Private VPN Ip of the peer on the private VPN subnet
	PrivateVpnPeerIp *string `json:"private-vpn-peer-ip,omitempty"`

	// Private VPN subnet
	PrivateVpnRange string `json:"private-vpn-range"`

	// Public key of the server
	PubKey string `json:"pub-key"`

	// Subnet to expose
	SubnetToExpose []string `json:"subnet-to-expose"`
}

// PostNetwork defines model for PostNetwork.
type PostNetwork struct {
	// Network name specified in the NSD
	NetworkName string `json:"network-name"`
}

// PostSap defines model for PostSap.
type PostSap struct {
	// Network name of the floating network, specified in the SAP information of the NSD
	FloatingNetworkName string `json:"floating-network-name"`

	// Network name specified in the NSD
	NetworkName string `json:"network-name"`
}

// POST to create all the network resources of a slice on a vim
type PostSliceResources struct {
	// Last subnet used by the peer
	ExcludeSubnet *string `json:"exclude-subnet,omitempty"`

	// Name of the networks specified in the NSD
	Networks []PostNetwork `json:"networks"`

	// SAP specified in the NSD
	ServiceAccessPoints []PostSap `json:"service-access-points"`

	// Id of the network slice owning the network resources, assigned by the Slicer
	SliceId string `json:"slice-id"`

	// Name of the VIM where to create the requested resources
	VimName string `json:"vim-name"`
}

// Sap defines model for Sap.
type Sap struct {
	// Network name of the floating network, specified in the SAP information of the NSD
	FloatingNetworkName string `json:"floating-network-name"`

	// Network name specified in the NSD
	NetworkName string `json:"network-name"`

	// Network address used for this subnet
	SubnetCidr string `json:"subnet-cidr"`
}

// SliceResources defines model for SliceResources.
type SliceResources struct {
	// Unique identifier of the set of network resources of the slice
	Id int `json:"id"`

	// Name of the networks specified in the NSD
	Networks []Network `json:"networks"`

	// SAP specified in the NSD
	ServiceAccessPoints []Sap `json:"service-access-points"`

	// Unique identifier assigned the Slicer
	SliceId string `json:"slice-id"`

	// Status of the resources
	Status string `json:"status"`

	// Name of the VIM where to create the requested resources
	VimName string `json:"vim-name"`
}

// DeleteNetResourcesParams defines parameters for DeleteNetResources.
type DeleteNetResourcesParams struct {
	// Id of the network slice owning the network resources
	SliceId string `json:"slice-id"`
}

// GetNetResourcesParams defines parameters for GetNetResources.
type GetNetResourcesParams struct {
	// Id of the network slice owning the network resources
	SliceId *string `json:"slice-id,omitempty"`
}

// PostNetResourcesJSONBody defines parameters for PostNetResources.
type PostNetResourcesJSONBody PostSliceResources

// PutNetResourcesIdGatewayJSONBody defines parameters for PutNetResourcesIdGateway.
type PutNetResourcesIdGatewayJSONBody PostGateway

// PostNetResourcesIdGatewayConnectionsJSONBody defines parameters for PostNetResourcesIdGatewayConnections.
type PostNetResourcesIdGatewayConnectionsJSONBody PostConnection

// PostNetResourcesJSONRequestBody defines body for PostNetResources for application/json ContentType.
type PostNetResourcesJSONRequestBody PostNetResourcesJSONBody

// PutNetResourcesIdGatewayJSONRequestBody defines body for PutNetResourcesIdGateway for application/json ContentType.
type PutNetResourcesIdGatewayJSONRequestBody PutNetResourcesIdGatewayJSONBody

// PostNetResourcesIdGatewayConnectionsJSONRequestBody defines body for PostNetResourcesIdGatewayConnections for application/json ContentType.
type PostNetResourcesIdGatewayConnectionsJSONRequestBody PostNetResourcesIdGatewayConnectionsJSONBody

