package main

import (
	"nextworks/nsm/internal/nsm"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("test ")

	subs := []string{"1", "2", "3"}
	s := nsm.SubnetsToString(subs)
	log.Info(s)
	ss := nsm.SubnetsToArray(s)
	log.Info(ss)
}
