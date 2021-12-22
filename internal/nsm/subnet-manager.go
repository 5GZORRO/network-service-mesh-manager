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

func (this *Network_manager) Next() *net.IPNet {
	if this.subnets_init != nil {
		ip := this.next_subnet
		net2, v := cidr.NextSubnet(ip, 28)
		if !v {
			this.next_subnet = net2
			return ip
		} else {
			log.Error("Error")
			this.next_subnet = nil
			return ip
		}
	} else {
		log.Error("Error")
		return nil
	}
}

func (this *Network_manager) NextSubnet(lastUsed string) *net.IPNet {
	_, net, err := net.ParseCIDR(lastUsed)
	if err != nil {
		log.Error("Error")
		return nil
	}
	net1, v := cidr.NextSubnet(net, 28)
	if v {
		log.Error("Error")
		this.next_subnet = nil
		return nil
	}
	net2, v := cidr.NextSubnet(net1, 28)
	if !v {
		this.next_subnet = net2
		return net1
	} else {
		log.Error("Error")
		this.next_subnet = nil
		return net1
	}

}

func NewNetworkManager(start string) *Network_manager {
	_, net, err := net.ParseCIDR(start)
	if err != nil {
		log.Error("Error")
		return nil
	}
	return &Network_manager{
		subnets_init: net,
		next_subnet:  net,
	}
}
