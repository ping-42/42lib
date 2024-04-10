package helpers

import (
	"net"
	"net/url"
)

// Check if there are any IPv6 addresses
func HasIPv6(targetIPs []net.IP) (hasIPv6 bool) {
	for _, ip := range targetIPs {
		if ip.To4() == nil {
			hasIPv6 = true
			break
		}
	}
	return
}

func ExtractDomainFromUrl(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	return parsed.Host, nil
}
