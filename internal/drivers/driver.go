package drivers

type VimDriver interface {
	CreateGatewayConnectivity(sliceId string, subnet string) (string, string, string, error)
	DeleteGatewayConnectivity(networkId string, subnetId string, routerId string) error
	GetGatewayConnectivity(networkId string, subnetId string, routerId string)
}
