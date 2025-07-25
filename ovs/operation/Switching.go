package operation

import (
	"context"
	"fmt"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
	"github.com/kwonkwonn/ovn-go-cms/ovs/util"
	"github.com/ovn-kubernetes/libovsdb/model"
	"github.com/ovn-kubernetes/libovsdb/ovsdb"
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
	for _,i:= range o.IPMapping{
		fmt.Println(i)
	}
	lsp = append(lsp, lsMute...)
	result,err := o.Client.Transact(context.Background(),lsp...)
	if err!=nil{
		fmt.Println("the problem is...", err)
	}
	fmt.Println(result)

	
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
		Addresses: []string{"router"},
		Options: map[string]string{"router-port":lrpuuid},
		}

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


func (o * Operator) DelSwitchPort(ip string)(error){
	lspuuid := o.CheckIPExistance(ip)
	if lspuuid == "" {
		return fmt.Errorf("no such switch port exist for ip %s", ip)
	}
	var lspDelOPS []ovsdb.Operation

	netAdd,err:=util.GetNetWorkInterface(ip)
	if err != nil {
		return fmt.Errorf("invalid IP format: %v %s %s", err, netAdd, ip)
	}
	
	dev, ok := o.IPMapping[netAdd+"1"]
	if !ok{
		return fmt.Errorf("no such switch port exist for ip %s", ip)
	}
	
	EXT,err := o.findDevByUUID(dev)
	if err != nil {
		return fmt.Errorf("no such device for uuid %s: %v", dev, err)
	}
	fmt.Println("EXT:", EXT)	
	s,ok := EXT.(*externalmodel.ExternSwitch)
		if !ok{
			return fmt.Errorf("no such switch exist for uuid %s", dev)
		}
		for j := 0; j < len(s.InternalSwitch.Ports); { 
			if s.InternalSwitch.Ports[j] == lspuuid {
				swMute ,_ := o.Client.Where(s.InternalSwitch).Mutate(s.InternalSwitch, model.Mutation{
				Field: &s.InternalSwitch.Ports,
				Mutator: ovsdb.MutateOperationDelete,
				Value: []string{lspuuid}, 

				})
				s.InternalSwitch.Ports = append(s.InternalSwitch.Ports[:j], s.InternalSwitch.Ports[j+1:]...)
				fmt.Println("swMute:", swMute)
				lspDelOPS = append(lspDelOPS, swMute...)
				break
						} 
			j++
	}
	fmt.Println("after processing EXT:", EXT)	
			
		
	

	// 포트 삭제
	lsp := &NBModel.LogicalSwitchPort{
		UUID: lspuuid,
	}
	lspDelOP, err := o.Client.Where(lsp).Delete()
	if err != nil {
		return fmt.Errorf("deleting switch port error %v", err)
	}
	lspDelOPS = append(lspDelOPS, lspDelOP...)

	result, err := o.Client.Transact(context.Background(), lspDelOPS...)
	if err != nil {
		return fmt.Errorf("deleting switch port error %v", err)
	}
	fmt.Println(result)

	delete(o.IPMapping, ip)
	util.SaveMapYaml(o.IPMapping)

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


func (o *Operator) DelSwitch(uuid string)(error){
	value,ok := o.ExternSwitchs[uuid]
	if !ok{
		return fmt.Errorf("no such switch exist")
	}
	
	lspDelOps:= make([]ovsdb.Operation, 0)


	for _,port := range value.InternalSwitch.Ports{
		lsp := &NBModel.LogicalSwitchPort{
			UUID: port,
		}
		lspDelOP, err := o.Client.Where(lsp).Delete()
		if err != nil {
			return fmt.Errorf("deleting switch port error %v", err)
		}
		lspDelOps = append(lspDelOps, lspDelOP...)

	}
	lsDelOp, err := o.Client.Where(value.InternalSwitch).Delete()
	if err != nil {
		return fmt.Errorf("deleting switch error %v", err)
	}
	lspDelOps = append(lspDelOps, lsDelOp...)

	result, err := o.Client.Transact(context.Background(), lspDelOps...)
	if err != nil {
		return fmt.Errorf("deleting switch port error %v", err)
	}
		fmt.Println(result)

	delete(o.ExternSwitchs, uuid)
	delete(o.IPMapping, value.IP)

	util.SaveMapYaml(o.IPMapping)

	return nil
}