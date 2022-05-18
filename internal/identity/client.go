package identityclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var timout = 30 // Timout for HTTP client, before returning: context deadline exceeded (Client.Timeout exceeded while awaiting headers)

// VPNaaSClient contains the VPNaaS Server information such as its IP and port
// its methods are the API offered by the server (Launch(), Connect_to_VPN(), Disconnect_to_VPN())
type IdentityClient struct {
	ip     net.IP
	port   string
	secret string
}

func New(addr net.IP, port string, secret string) *IdentityClient {
	return &IdentityClient{
		ip:     addr,
		port:   port,
		secret: secret,
	}
}

// Start the VPN service calling the /launch() endpoint
func (client *IdentityClient) CreateKeyPair() (string, error) {
	c := http.Client{Timeout: time.Duration(timout) * time.Second}

	log.Trace("IdentityClient {", client.ip.String(), " ", client.port, "} -- Requesting a key pair ")
	req, err := http.NewRequest("GET", "http://"+client.ip.String()+":"+client.port+"/authentication/operator_key_pair", nil)
	if err != nil {
		log.Error(err)
		return "", err
	}

	// req.Header.Add("Accept", `application/json`)
	req.Header.Add("shared-secret", client.secret)
	resp, err := c.Do(req)
	// send request
	if err != nil {
		log.Error(err)
		return "", err
	}
	log.Debug("IdentityClient {", client.ip.String(), " ", client.port, "} -- Response status: ", resp.Status)
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
			return "", err
		}
		// var info KeyPair
		// json.Unmarshal(body, &info)
		bodyString := string(body)
		log.Trace("IdentityClient {", client.ip.String(), " ", client.port, "} -- Response Body ", bodyString)
		jwt := strings.Trim(bodyString, "\"")
		return jwt, nil
	} else {
		return "", err
	}
}

func (client *IdentityClient) VerifyKeyPair(in *VerifyKeyPairBody) bool {
	c := http.Client{Timeout: time.Duration(timout) * time.Second}

	jsonBody, _ := json.Marshal(in)
	log.Trace("IdentityClient {", client.ip.String(), " ", client.port, "} -- Verifying key pair ", in)
	req, err := http.NewRequest("POST", "http://"+client.ip.String()+":"+client.port+"/authentication/operator_key_pair/verify", bytes.NewReader(jsonBody))
	if err != nil {
		log.Error(err)
		return false
	}
	req.Header.Add("shared-secret", client.secret)
	// send request
	resp, err := c.Do(req)
	if err != nil {
		log.Error(err)
		return false
	}
	log.Debug("IdentityClient {", client.ip.String(), " ", client.port, "} -- Response status: ", resp.Status)
	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
}
