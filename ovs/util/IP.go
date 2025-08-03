package util

import (
	"fmt"
	"net"
	"strings"
)




func GetNetworkAddress(cidrIP string) (string, error) {
	_, ipNet, err := net.ParseCIDR(cidrIP)
	if err != nil {
		return "", fmt.Errorf("유효하지 않은 CIDR 형식입니다: %w", err)
	}

	return ipNet.IP.String(), nil
}// 20.20.20.12/24 란 ip가 주어지면 20.20.20.0로 파싱

func GetNetWorkSignifier(ip string) (string,  error) {

	ips:=strings.Split(ip,".")
	return ips[0]+"."+ips[1]+"."+ips[2]+".", nil
}