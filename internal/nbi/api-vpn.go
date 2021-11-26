package nbi

import "github.com/gin-gonic/gin"

// (DELETE /gateways/{id}/vpn/configuration)
func (obj *ServerInterfaceImpl) DeleteGatewaysIdVpnConfiguration(c *gin.Context, id int) {}

// (GET /gateways/{id}/vpn/configuration)
func (obj *ServerInterfaceImpl) GetGatewaysIdVpnConfiguration(c *gin.Context, id int) {}

// (PUT /gateways/{id}/vpn/configuration)
func (obj *ServerInterfaceImpl) PutGatewaysIdVpnConfiguration(c *gin.Context, id int) {}

// (GET /gateways/{id}/vpn/connections)
func (obj *ServerInterfaceImpl) GetGatewaysIdVpnConnections(c *gin.Context, id int) {}

// (POST /gateways/{id}/vpn/connections)
func (obj *ServerInterfaceImpl) PostGatewaysIdVpnConnections(c *gin.Context, id int) {}

// (DELETE /gateways/{id}/vpn/connections/{connectionid}/)
func (obj *ServerInterfaceImpl) DeleteGatewaysIdVpnConnectionsConnectionid(c *gin.Context, id int, connectionid int) {
}

// (GET /gateways/{id}/vpn/connections/{connectionid}/)
func (obj *ServerInterfaceImpl) GetGatewaysIdVpnConnectionsConnectionid(c *gin.Context, id int, connectionid int) {
}
