package util

import (
	"crypto/rand"
	"fmt"
	"net"
)

//mac prefix(kvm)=  52:54:00   // 01010010 : 01010100 : 00000000
//mac 중복성은 아직 고려하지 않음


func MacGenerator() (string, error) {
	// MAC 주소는 6바이트입니다.
	mac := make(net.HardwareAddr, 6)

	mac[0] = 0x52
	mac[1] = 0x54
	mac[2] = 0x00

	_, err := rand.Read(mac[3:])
	if err != nil {
		return "", fmt.Errorf("MAC 주소의 나머지 부분을 생성하는 데 실패했습니다: %w", err)
	}
	
	return mac.String(), nil
}

