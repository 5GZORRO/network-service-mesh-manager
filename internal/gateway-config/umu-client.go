package gatewayconfig

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type UmuClient struct {
	ip   net.IP
	port string
}

func New(addr net.IP, port string) *UmuClient {
	return &UmuClient{
		ip:   addr,
		port: port,
	}
}

func (client *UmuClient) Start(ipRange string, netInterface string, port string) bool {
	log.Trace("Starting VPN Server... mgmt info {", client.ip.String(), " ", client.port, "}")
	c := http.Client{Timeout: time.Duration(1) * time.Second}

	bodyrequest := PostLaunch{
		IpRange:      ipRange,
		NetInterface: netInterface,
		Port:         port,
	}
	jsonBody, _ := json.Marshal(bodyrequest)
	log.Trace("Starting VPN Server... with body request ", bodyrequest)
	req, err := http.NewRequest("POST", "http://"+client.ip.String()+":"+client.port+"/launch", bytes.NewReader(jsonBody))
	if err != nil {
		log.Error(err)
		return false
	}
	req.Header.Add("Accept", `application/json`)
	// send request
	resp, err := c.Do(req)
	if err != nil {
		log.Error(err)
		return false
	}
	log.Debug("Starting VPN Server... Response status: ", resp.Status)
	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
}

// retrieves the current VPN Server configuration
func (client *UmuClient) GetCurrentConfiguration() *VpnInfo {
	log.Trace("Get VPN Server configuration... mgmt info {", client.ip.String(), " ", client.port, "}")
	resp, err := http.Get("http://" + client.ip.String() + ":" + client.port + "/get_configuration")
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Trace("Get VPN Server configuration... Response status: ", resp.Status)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil
	}
	var info VpnInfo
	json.Unmarshal(body, &info)
	log.Trace("Get VPN Server configuration... Body ", info)
	return &info
}

func (client *UmuClient) Connect(peerIp string, peerPort string, exposedNets string) bool {
	log.Trace("Connecting to peer... ")
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	requestConnect := PostConnect{
		IpAddressServer:   peerIp,
		PortServer:        peerPort,
		IPRangeToRedirect: exposedNets,
	}
	jsonBody2, _ := json.Marshal(requestConnect)
	req, err := http.NewRequest("POST", "http://"+client.ip.String()+":"+client.port+"/connect_to_VPN", bytes.NewReader(jsonBody2))
	if err != nil {
		log.Error(err)
		return false
	}
	req.Header.Add("Accept", `application/json`)
	// send request
	resp, err := c.Do(req)
	if err != nil {
		log.Error(err)
		return false
	}
	log.Debug("Connecting to peer... Response status: ", resp.Status)
	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
}

func (client *UmuClient) Disconnect(peerIP string, peerPort string) bool {
	log.Trace("Connecting to peer... ")
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	requestDisconnect := PostDisconnect{
		IpAddressServer: peerIP,
		PortServer:      peerPort,
	}

	jsonBody2, _ := json.Marshal(requestDisconnect)
	req, err := http.NewRequest("POST", "http://"+client.ip.String()+":"+client.port+"/disconnect_to_VPN", bytes.NewReader(jsonBody2))
	if err != nil {
		log.Error(err)
		return false
	}
	req.Header.Add("Accept", `application/json`)
	// send request
	resp, err := c.Do(req)
	if err != nil {
		log.Error(err)
		return false
	}
	log.Debug("Disconnecting from client... Response status: ", resp.Status)
	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
}
