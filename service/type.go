package service

import "github.com/kwonkwonn/ovn-go-cms/ovs/operation"



type NewInstanceResult struct{
	IP 	string 	`json:"IP"`
	MacAddress string  `json:"macAddr"`
 	IfaceID string 	`json:"sdnUUID"`
	Detail error `json:"error"`
}


type NewInstanceRequeset struct {
	RequestSubnet string `json:"Subnet"` //새로운 서브넷 생성, 기존 서브넷 사용시 동일하게 사용됩니다
	// 새로운 서브넷을 생성할 시 컨트롤에서 할당할 서브넷을, 아닐 시 현재 할당 된 서브넷을 제공합니다
}







/////////////////////////////////////////////
type Handler struct{
	Operator *operation.Operator
}