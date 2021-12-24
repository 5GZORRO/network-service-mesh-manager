package gatewayconfig

type GatewayHttpClient interface {
	GetCurrentConfiguration()
	Start()
	Connect()
	Disconnect()
}
