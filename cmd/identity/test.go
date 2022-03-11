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

	ip := "172.28.3.153"
	// port1 := "6800"
	port2 := "6600"
	// client1 := identityclient.New(net.ParseIP(ip), port1, "5gzorroidportalnsmm")

	// res, _ := client1.CreateKeyPair()
	// log.Debug(res)

	client2 := identityclient.New(net.ParseIP(ip), port2, "5gzorroidportalnsmm")
	input := identityclient.VerifyKeyPairBody{
		Did:       "VdjbRWzuTwSj84twf92q5S",
		PubKey:    "72KkcAZHocXxMUyAnBJ0i5C0VM/phNuoCU7KF8s3brw=",
		Timestamp: "1646923005",
	}

	output := client2.VerifyKeyPair(&input)
	log.Debug(output)
}
