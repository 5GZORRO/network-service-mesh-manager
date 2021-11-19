// Copyright 2019 DeepMap, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=types.cfg.yaml ../../petstore-expanded.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml ../../petstore-expanded.yaml

package NsmmApi

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type GatewayInterface struct {
	Gateways string
	Lock     sync.Mutex
}

func NewGatewayInterface() *GatewayInterface {
	return &GatewayInterface{
		Gateways: "ciao",
	}
}

// Here, we implement all of the handlers in the ServerInterface
func (p *GatewayInterface) GetGateway(ctx *gin.Context, params GetGatewayParams) {

	ctx.Status(http.StatusAccepted)
}

func (p *GatewayInterface) PostGateway(ctx *gin.Context, params PostGatewayParams) {
	ctx.Status(http.StatusNoContent)
}

func (p *GatewayInterface) DeleteGateway(ctx *gin.Context, params DeleteGatewayParams) {
	ctx.Status(http.StatusNoContent)
}
