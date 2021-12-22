package nsm

import (
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	log "github.com/sirupsen/logrus"
)

type Network_manager struct {
	subnets_init *net.IPNet
	next_subnet  *net.IPNet
}

func (ob *Network_manager) NextSubnet() *net.IPNet {
	if ob.subnets_init != nil {
		ip := ob.next_subnet
		net2, v := cidr.NextSubnet(ip, 28)
		if !v {
			ob.next_subnet = net2
			return ip
		} else {
			log.Error("Error")
			ob.next_subnet = nil
			return ip
		}
	} else {
		log.Error("Error")
		return nil
	}
}

// NewNetworkManager returns a Network_Manager to allocate subnets.
// start is the starting subnet
// skip indicates if the first subnet should be avoid (in case of an exclude)
func NewNetworkManager(start string, skip bool) *Network_manager {
	if !skip {
		_, net, err := net.ParseCIDR(start)
		if err != nil {
			log.Error("Error")
			return nil
		}
		return &Network_manager{
			subnets_init: net,
			next_subnet:  net,
		}
	} else {
		_, net, err := net.ParseCIDR(start)
		if err != nil {
			log.Error("Error")
			return nil
		}
		net1, v := cidr.NextSubnet(net, 28)
		if v {
			log.Error("Error")
			return nil
		}
		return &Network_manager{
			subnets_init: net1,
			next_subnet:  net1,
		}
	}
}
