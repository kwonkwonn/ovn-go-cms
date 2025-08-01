package externalmodel

import "strconv"




func GetNetInt(routers EXRList, ip string)[]NetInt {
	netInst:= make([]NetInt, 0)
	for _, router := range routers {
		if netint , ok := router.subNetworks[ip]; ok {
			netInst = append(netInst, netint)
		}
	}

	return netInst
}


func AddNetInt(routers EXRList, ip string, netint NetInt) {
	for _, router := range routers {
		if _, ok := router.subNetworks[ip]; !ok {
			router.subNetworks[ip] = netint
		}
	}
}

func FindRemainIP(routers EXRList, subnet string, travertype portType ) string {

	if travertype == SWITCH {
		for i := 1; i <= 10; i++ {
		IP := subnet + strconv.Itoa(i)
		if len(GetNetInt(routers, IP)) == 0 {
			return IP
		}
		}
	}else{
		for i := 11; i <= 254; i++ {
			IP := subnet + strconv.Itoa(i)
			if len(GetNetInt(routers, IP)) == 0 {
				return IP
			}
		}
	}
	

	return ""
}