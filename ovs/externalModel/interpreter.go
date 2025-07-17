package externalmodel




func (ER ExternRouter) ReturnUUID() string{
	return ER.UUID
}

func (ES ExternSwitch)ReturnUUID() string{
	return ES.UUID
}

//func()string 은 디바이스의 uuid 만을 배출함
// func (ER ExternRouter) AddPort(client.Client, ...func()(string))([]ovsdb.Operation,error){
	
// }

// func (ES ExternSwitch)AddPort(client.Client, ...func()(string))([]ovsdb.Operation,error){

// }