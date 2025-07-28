package externalmodel






type EXRList map[string]*ExternRouter
type EXSList map[string]*ExternSwitch




func (EXR EXRList) GetRouter(uuid string) *ExternRouter {
	if router, ok := EXR[uuid]; ok {
		return router
	}
	return nil
}


func (EXS EXSList) GetSwitch(uuid string) *ExternSwitch {
	if switchDevice, ok := EXS[uuid]; ok {
		return switchDevice
	}
	return nil
}


