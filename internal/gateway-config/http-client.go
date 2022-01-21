package gatewayconfig

type GatewayHttpClient interface {
	GetCurrentConfiguration()
	Launch()
	Connect()
	Disconnect()
}
