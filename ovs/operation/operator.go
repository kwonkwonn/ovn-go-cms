package operation

import (
	"fmt"
	"strconv"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	"github.com/ovn-org/libovsdb/client"
)





type Operator struct{
	Client client.Client
	ExternRouters map[string]*externalmodel.ExternRouter
	ExternSwitchs map[string]*externalmodel.ExternSwitch
	IPMapping map[string]string // device uuid
}


func (o* Operator) CheckIPExistance(subnet string)(string){
	dev,ok := o.IPMapping[subnet+"1"]
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
	for i:=1; i<10;{
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
	for i:=11; i<= 254; {
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
	for i:=1; i<= 10; {
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
		dev,ok := o.ExternRouters[uuid]
		if !ok{ 
			return nil,fmt.Errorf("no such device for uuid")
		} 
		return dev,nil
	}
	return dev,nil
}