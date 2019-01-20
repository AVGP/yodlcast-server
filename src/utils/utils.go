package utils

import (
	"errors"
	"net"
)

// DeviceUUID is the UPnP UUID we are broadcasting and announcing
const DeviceUUID = "c000ffee-cafe-c0c0-dead-c000ffffeeee"

// GetLocalIP returns the first non-loopback IPv4 address of this device
func GetLocalIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}
	return nil, errors.New("No local IPv4 address for local network found")
}
