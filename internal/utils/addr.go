package utils

import (
	"net"
)

func IpAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	ip := ""
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}

	return ip
}
