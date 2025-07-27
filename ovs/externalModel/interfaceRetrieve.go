package externalmodel






func (p port) GetUUID() string {
	return p.UUID
}


func (p port) GetBottomPorts() string {
	if p.ConnectedSwitch != nil {
		return p.ConnectedSwitch.UUID
	}
	return ""
}
func (p port) GetTopPorts() string {
	if p.ConnectedRouter != nil {
		return p.ConnectedRouter.UUID
	}
	return ""
}