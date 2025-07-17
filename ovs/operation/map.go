package operation

import (
	"fmt"
	"strconv"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
)



func (o* Operator) SwitchMapAdd(EXS *externalmodel.ExternSwitch)(error){
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.ExternSwitchs[EXS.UUID]= EXS

	return nil
}
func (o* Operator) RouterMapAdd(EXR *externalmodel.ExternRouter)(error){
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.ExternRouters[EXR.UUID]= EXR
	return nil
}

func (o* Operator) findRouterByUUID(uuid string) (*externalmodel.ExternRouter, error){
	o.mutex.Lock()
	defer o.mutex.Unlock()
	dev,ok := o.ExternRouters[uuid]
	if !ok{
		return nil,fmt.Errorf("no such device for uuid")
	}
	return dev,nil
}

func (o* Operator) findSwitchByUUID(uuid string) (*externalmodel.ExternSwitch, error){
	o.mutex.Lock()
	defer o.mutex.Unlock()
	dev,ok := o.ExternSwitchs[uuid]
	if !ok{ 
		return nil,fmt.Errorf("no such device for uuid")
	} 
	return dev,nil
}

func (o* Operator)findDeviceByIP(ip string)(string,error){
	o.mutex.Lock()
	defer o.mutex.Unlock()
	dev,ok := o.IPMapping[ip]
	if !ok{
		return "",fmt.Errorf("no such device correspond that ip")
	}
	return dev,nil
}

func (o*Operator) assignDevByIp( ip string ,uuid string)(error){
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.IPMapping[ip]=uuid
	return nil
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
		uuid,_:= o.findDeviceByIP(IP)
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
		uuid,_ := o.findDeviceByIP(IP)
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