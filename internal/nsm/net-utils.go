package nsm

import (
	"errors"
	"strconv"
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
