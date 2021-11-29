package nbi

import "github.com/gin-gonic/gin"

// TODO implement connection logic
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
