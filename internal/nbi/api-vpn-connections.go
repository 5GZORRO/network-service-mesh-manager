package nbi

import (
	"github.com/gin-gonic/gin"
	// log "github.com/sirupsen/logrus"
)

// TODO implement connection logic
// (GET /gateways/{id}/vpn/connections)
func (obj *ServerInterfaceImpl) GetNetResourcesIdGatewayConnections(c *gin.Context, id int) {}

// (POST /gateways/{id}/vpn/connections)
func (obj *ServerInterfaceImpl) PostNetResourcesIdGatewayConnections(c *gin.Context, id int) {}

// (DELETE /gateways/{id}/vpn/connections/{connectionid}/)
func (obj *ServerInterfaceImpl) DeleteNetResourcesIdGatewayConnectionsCid(c *gin.Context, id int, connectionid int) {
}

// (GET /gateways/{id}/vpn/connections/{connectionid}/)
func (obj *ServerInterfaceImpl) GetNetResourcesIdGatewayConnectionsCid(c *gin.Context, id int, connectionid int) {
}
