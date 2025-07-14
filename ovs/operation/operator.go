package operation

import (
	"fmt"
	"strconv"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	"github.com/ovn-org/libovsdb/client"
)

//
type KNOWN_DEVICES string

const ( 
 UPLINK KNOWN_DEVICES = "UPLINK"
 //ip 가 할당 되어 있지 않지만 한번씩 필요한 놈들,
 // o.ipmap map[string]string 에 스트링으로 저장 
 ROUTER KNOWN_DEVICES = "10.5.15.4" // 추후에 getenv등으로 숨김 , const 라서 그렇게 초기화 될지는 몰?루
)

type Operator struct{
	Client client.Client
	ExternRouters map[string]*externalmodel.ExternRouter
	ExternSwitchs map[string]*externalmodel.ExternSwitch
	IPMapping map[string]string // device uuid
}


func (o* Operator) CheckIPExistance(IP string)(string){
	dev,ok := o.IPMapping[IP]
	if !ok{
		return ""
	}
	return dev
}

func (o* Operator) SwitchesPortConnect(uuids []string,IP string ,VMUUID string, VMMac string)(error){
	for _,uuid:=range uuids{
		
		o.AddSwitchAPort(uuid,IP,VMUUID,VMMac)
	}

	return nil
}


func (o * Operator)FindExistdev(subnet string )([]string){
	var devs =make([]string,0)
	for i:=1; i<10; i++{
		dev, ok:= o.IPMapping[subnet + strconv.Itoa(i)]
		if !ok{
			return devs
		}
		devs=append(devs, dev)
	}	
	return devs
}


func (o* Operator) AvailableIP_VM(subnet string)(ip string){
	//   /24로 가정 , 1~10번은 가상 디바이스에 우선적으로 할당
	for i:=11; i<= 254; i++{
		IP := subnet+ strconv.Itoa(i)
		uuid := o.CheckIPExistance(IP)
		if uuid ==""{
			return IP
		}
	}
	return ""
}
func (o* Operator) AvailableIP_Dev(subnet string)(ip string){
	//   /24로 가정 , 1~10번은 가상 디바이스에 우선적으로 할당
	for i:=1; i<= 10; i++{
		IP := subnet+ strconv.Itoa(i)
		uuid := o.CheckIPExistance(IP)
		if uuid ==""{
			return IP
		}
	}
	return ""
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