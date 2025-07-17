package util

import (
	"fmt"
	"net"
)




func GetNetworkAddress(cidrIP string) (string, error) {
	_, ipNet, err := net.ParseCIDR(cidrIP)
	if err != nil {
		return "", fmt.Errorf("유효하지 않은 CIDR 형식입니다: %w", err)
	}

	return ipNet.IP.String(), nil
}