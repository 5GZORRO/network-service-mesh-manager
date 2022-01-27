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

// VPNaaSClient contains the VPNaaS Server information such as its IP and port
// its methods are the API offered by the server (Launch(), Connect_to_VPN(), Disconnect_to_VPN())
type VPNaaSClient struct {
	ip   net.IP
	port string
}

func New(addr net.IP, port string) *VPNaaSClient {
	return &VPNaaSClient{
		ip:   addr,
		port: port,
	}
}

// Start the VPN service calling the /launch() endpoint
func (client *VPNaaSClient) Launch(ipRange string, netInterface string, port string) bool {
	c := http.Client{Timeout: time.Duration(10) * time.Second}

	bodyrequest := PostLaunch{
		IpRange:      ipRange,
		NetInterface: netInterface,
		Port:         port,
		Environment:  "local",
	}
	jsonBody, _ := json.Marshal(bodyrequest)
	log.Trace("VPNaaS {", client.ip.String(), " ", client.port, "} -- Starting with body ", bodyrequest)
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
	log.Debug("VPNaaS {", client.ip.String(), " ", client.port, "} -- Response status: ", resp.Status)
	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
}

// Retrieve the current VPN Server configuration
func (client *VPNaaSClient) GetCurrentConfiguration() *VpnInfo {
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

// Connect to a client (peer2), calling the /connect_to_VPN
func (client *VPNaaSClient) Connect(peerIp string, peerPort string, remoteIPs string, localIPs string) bool {
	c := http.Client{Timeout: time.Duration(10) * time.Second}
	requestConnect := PostConnect{
		IpAddressServer: peerIp,
		PortServer:      peerPort,
		RemoteSubnet:    remoteIPs,
		LocalSubnet:     localIPs,
		Environment:     "local",
	}
	jsonBody2, _ := json.Marshal(requestConnect)
	log.Trace("VPNaaS {", client.ip.String(), " ", client.port, "} -- Connecting to peer... with body ", requestConnect)
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
	log.Debug("VPNaaS {", client.ip.String(), " ", client.port, "} -- Connecting to peer... Response status: ", resp.Status)
	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
}

// Disconnect from a client (peer), calling the /disconnect_to_VPN
func (client *VPNaaSClient) Disconnect(peerIP string, peerPort string) bool {
	c := http.Client{Timeout: time.Duration(10) * time.Second}
	requestDisconnect := PostDisconnect{
		IpAddressServer: peerIP,
		PortServer:      peerPort,
	}

	jsonBody2, _ := json.Marshal(requestDisconnect)
	log.Trace("VPNaaS {", client.ip.String(), " ", client.port, "} -- Disconnecting to peer... with body ", requestDisconnect)
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
	log.Debug("VPNaaS {", client.ip.String(), " ", client.port, "} -- Disconnecting from client... Response status: ", resp.Status)
	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
}
