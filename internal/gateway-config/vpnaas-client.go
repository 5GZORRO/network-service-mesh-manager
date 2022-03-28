package gatewayconfig

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	identityclient "nextworks/nsm/internal/identity"

	log "github.com/sirupsen/logrus"
)

var timout = 30 // Timout for HTTP client, before returning: context deadline exceeded (Client.Timeout exceeded while awaiting headers)

// VPNaaSClient contains the VPNaaS Server information such as its IP and port
// its methods are the API offered by the server (Launch(), Connect_to_VPN(), Disconnect_to_VPN())
type VPNaaSClient struct {
	ip          net.IP
	port        string
	environment string
}

func New(addr net.IP, port string, env string) *VPNaaSClient {
	return &VPNaaSClient{
		ip:          addr,
		port:        port,
		environment: env,
	}
}

// Function to Launch the VPNaaS in Test mode (local)
// In case of test environment (local), keys to be used are local to the VM, not passed in the launch
func (client *VPNaaSClient) LaunchTest(ipRange string, netInterface string, port string) bool {
	c := http.Client{Timeout: time.Duration(timout) * time.Second}
	bodyrequest := PostLaunch{
		IpRange:      ipRange,
		NetInterface: netInterface,
		Port:         port,
		Environment:  client.environment,
	}

	jsonBody, _ := json.Marshal(bodyrequest)
	log.Trace("VPNaaS {", client.ip.String(), " ", client.port, "} -- Starting with body (local-test) ", bodyrequest)
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

// Start the VPN service calling the /launch() endpoint
func (client *VPNaaSClient) Launch(ipRange string, netInterface string, port string, keyPair *identityclient.KeyPair, idependpoint string) bool {
	c := http.Client{Timeout: time.Duration(timout) * time.Second}
	// In case of Prod environment (Zorro testbed), keys to be used are the one retrieved from the ID&P
	bodyrequest := PostLaunch{
		IpRange:      ipRange,
		NetInterface: netInterface,
		Port:         port,
		Environment:  client.environment,
		IDMEndpoint:  idependpoint,
		Did:          keyPair.Did,
		PubKey:       keyPair.PubKey,
		PrivKey:      keyPair.PrivKey,
		Timestamp:    keyPair.Timestamp,
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
	c := http.Client{Timeout: time.Duration(timout) * time.Second}
	requestConnect := PostConnect{
		IpAddressServer: peerIp,
		PortServer:      peerPort,
		RemoteSubnet:    remoteIPs,
		LocalSubnet:     localIPs,
		Environment:     client.environment,
	}
	jsonBody2, _ := json.Marshal(requestConnect)
	log.Trace("VPNaaS {", client.ip.String(), " ", client.port, "} -- Requesting connection to peer... with body ", requestConnect)
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
	log.Debug("VPNaaS {", client.ip.String(), " ", client.port, "} -- Response status: ", resp.Status)
	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
}

// Disconnect from a client (peer), calling the /disconnect_to_VPN
func (client *VPNaaSClient) Disconnect(peerIP string, peerPort string) bool {
	c := http.Client{Timeout: time.Duration(timout) * time.Second}
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
