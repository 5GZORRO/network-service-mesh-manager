package main

// Executable to test the VPNaaS Client and the interaction with the VPNaaS module (UMU)

import (
	"net"
	identityclient "nextworks/nsm/internal/identity"

	log "github.com/sirupsen/logrus"
)

func main() {
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	log.SetLevel(log.TraceLevel)

	// "172.28.3.153" 6600
	ip := "172.28.3.153"
	// ip := "127.0.0.1"
	port1 := "6800"
	// port2 := "8083"
	client1 := identityclient.New(net.ParseIP(ip), port1, "5gzorroidportalnsmm")

	res, _ := client1.CreateKeyPair()
	// Now the createKeyPair returns a string, which should be passes as it is to the VPNaaS
	log.Debug(res)

}
