package operation

import (
	"context"
	"fmt"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
	"github.com/kwonkwonn/ovn-go-cms/ovs/util"
	"github.com/ovn-org/libovsdb/model"
	"github.com/ovn-org/libovsdb/ovsdb"
)


func (o * Operator) AddSwitchAPort(SWUUID string, ip string, uuid string , mac string)(error){
	value ,ok := o.ExternSwitchs[SWUUID]; 
	if !ok{
		return fmt.Errorf("no such switch exist")
	}

	newSP:= &NBModel.LogicalSwitchPort{
		UUID: string(uuid),
		Name: string(uuid),
		Type: "vif",
		}
	Address := fmt.Sprintf("%s %s",mac , ip)
	newSP.Addresses=append(newSP.Addresses, Address)
	
	lsp , err := o.Client.Create(newSP)
	if err!=nil{
		return fmt.Errorf("%v", err)
	}
	
	value.InternalSwitch.Ports = append(value.InternalSwitch.Ports, newSP.UUID)
	lsMute,_  := o.Client.Where(value.InternalSwitch).Mutate( value.InternalSwitch, model.Mutation{
		Field: &value.InternalSwitch.Ports,
		Mutator: ovsdb.MutateOperationInsert,
		Value: value.InternalSwitch.Ports,	
		
	}) 
	o.IPMapping[ip]= newSP.UUID
	 
	lsp = append(lsp, lsMute...)
	result,err := o.Client.Transact(context.Background(),lsp...)
	if err!=nil{
		fmt.Println("the problem is...", err)
	}
	fmt.Println(result)


	o.IPMapping[ip]= SWUUID
	util.SaveMapYaml(o.IPMapping)
	
	//ip가 할당되는 순간 Map 에 저장

	return nil
}


func (o * Operator) AddSwitchAPort_Router(SWUUID string, lrpuuid string , uuid string)(error){
	value ,ok := o.ExternSwitchs[SWUUID]; 
	if !ok{
		return fmt.Errorf("no such switch exist")
	}

	newSP:= &NBModel.LogicalSwitchPort{
		UUID: string(uuid),
		Name: string(uuid),
		Type: "router",
		Options: map[string]string{"router-port":lrpuuid},
		}
	Address := "router"
	newSP.Addresses=append(newSP.Addresses, Address)
	
	lsp , err := o.Client.Create(newSP)
	if err!=nil{
		return fmt.Errorf("%v", err)
	}
	value.InternalSwitch.Ports = append(value.InternalSwitch.Ports, newSP.UUID)
	lsMute,_  := o.Client.Where(value.InternalSwitch).Mutate( value.InternalSwitch, model.Mutation{
		Field: &value.InternalSwitch.Ports,
		Mutator: ovsdb.MutateOperationInsert,
		Value: value.InternalSwitch.Ports,	
		
	}) 
	 
	lsp = append(lsp, lsMute...)
	result,err := o.Client.Transact(context.Background(),lsp...)
	if err!=nil{
		fmt.Println("the problem is...", err)
	}
	fmt.Println(result)
	
	//ip가 할당되는 순간 Map 에 저장

	return nil
}



func (o * Operator) AddSwitch (ip string) (uuid string ,error error){
	return  o.addSwitch(ip)
	
}


func (o *Operator) addSwitch (ip string) (string, error) {
	uuid ,err:=util.UUIDGenerator()
	if err!=nil{
		return "",fmt.Errorf("creating switch error %v",err)
	}
	newSwitch:=&NBModel.LogicalSwitch{
		UUID:uuid.String(),
		Name: uuid.String(),
	}
	fmt.Println(newSwitch)
	ls,err:= o.Client.Create(newSwitch)	
	if err!=nil{
		return "",fmt.Errorf("creating switch error %v",err)
	}
	
	result , err := o.Client.Transact(context.Background(),ls...)
	if err!=nil{
		return uuid.String() ,fmt.Errorf("creating switch error %v",err)
	}
	fmt.Println(result)

	o.IPMapping[ip] = uuid.String()

	o.ExternSwitchs[uuid.String()]= &externalmodel.ExternSwitch{
		UUID: uuid.String(),
		InternalSwitch: newSwitch,
		IP: ip,
	}
	return uuid.String(),nil

}