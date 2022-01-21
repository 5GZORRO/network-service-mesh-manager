package nsm

import (
	"errors"
	"strconv"
	"strings"
)

func parsePort(port string) (uint16, error) {
	portInt, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return 0, err
	}
	if portInt == 0 {
		return 0, errors.New("0 is not a valid port number")
	}
	result := uint16(portInt)
	return result, nil
}

func parsePortToString(port uint16) string {
	return strconv.Itoa(int(port))
}

func SubnetsToString(subnets []string) string {
	stringSubs := ""
	for i, sub := range subnets {
		if i == 0 {
			stringSubs = stringSubs + sub
		} else {
			stringSubs = stringSubs + "," + sub
		}
	}
	return stringSubs
}

func SubnetsToArray(subnets string) []string {
	return strings.Split(subnets, ",")
}

// parseExposedSubnets returns a string to be used in WireGuard file and VPNaaS API
func ParseExposedSubnets(subnets []string) string {
	stringSubs := ""
	for i, sub := range subnets {
		if i == 0 {
			stringSubs = stringSubs + sub
		} else {
			stringSubs = stringSubs + ", " + sub
		}
	}
	return stringSubs
}
