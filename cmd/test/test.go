package main

import (
	"nextworks/nsm/internal/nsm"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("test ")

	mng := nsm.NewNetworkManager("192.168.161.16/28", true)
	log.Info(*mng)
	nn := mng.NextSubnet()
	log.Info(*mng)
	log.Info(nn.String())

	nn = mng.NextSubnet()
	log.Info(*mng)
	log.Info(nn)

}
