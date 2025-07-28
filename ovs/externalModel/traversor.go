package externalmodel






func GetNetInt(routers EXRList, ip string)[]NetInt {
	netInst:= make([]NetInt, 0)
	for _, router := range routers {
		if netint , ok := router.ports[ip]; ok {
			netInst = append(netInst, netint)
		}
	}
	return netInst
}