package operation

import (
	"fmt"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"

	"github.com/ovn-kubernetes/libovsdb/client"
)

type isOCcupied string

const (
	Occupied   isOCcupied = "occupied"
	UnOccupied isOCcupied = "unoccupied"
) // subnet 검사를 위해 사용, xxx.xxx.xxx.0에 저장 됨


type KNOWN_DEVICES string

const ( 
 UPLINK KNOWN_DEVICES = "UPLINK"
 //ip 가 할당 되어 있지 않지만 한번씩 필요한 놈들,
 // o.ipmap map[string]string 에 스트링으로 저장 
 DEFAULT_GATEWAY KNOWN_DEVICES = "10.5.15.1"
 ROUTER KNOWN_DEVICES = "10.5.15.4" // 추후에 getenv등으로 숨김 , const 라서 그렇게 초기화 될지는 몰?루
)


type Operator struct{
	Client client.Client
	ExternRouters externalmodel.EXRList
	ExternSwitchs externalmodel.EXSList
	// IPMapping map[string]string// device uuid
}



func (o* Operator) IPMapToDev(IP string)(externalmodel.NetInt){
	list:= externalmodel.GetNetInt(o.ExternRouters, IP)
 if len(list) > 0 {
		return list[0]
	}
	return nil
}

func (o* Operator) SwitchesPortConnect(uuids []string,IP string ,VMUUID string, VMMac string)(error){
	for _,uuid:=range uuids{
		
		o.AddSwitchAPort(uuid,IP,VMUUID,VMMac)
	}

	return nil
}


func (o* Operator) findDevByUUID(uuid string) (any, error){
	dev,ok := o.ExternRouters[uuid]
	if !ok{
		dev,ok := o.ExternSwitchs[uuid]
		if !ok{ 
			return nil,fmt.Errorf("no such device for uuid")
		} 
		return dev,nil
	}
	return dev,nil
}