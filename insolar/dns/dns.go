package dns

import (
	"net"
	"strings"
)

// GetIPFromDomain returns IP address string from domain.
func GetIPFromDomain(domain string) (string, error) {
	woPort := strings.Split(domain, ":")
	address := woPort[0]
	var port string
	if len(woPort) > 1 {
		port = ":" + woPort[1]
	}

	ips, err := net.LookupIP(address)
	if err != nil {
		return "", err
	}
	return ips[0].String() + port, nil
}
