package operation

import (
	"context"
	"fmt"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	"github.com/kwonkwonn/ovn-go-cms/ovs/util"
	"github.com/ovn-kubernetes/libovsdb/ovsdb"
)


func (o * Operator) AddSwitchAPort(SWUUID string, InstanceIP string, uuid string , mac string)(*externalmodel.SwitchPort,error){
	ops:= make([]ovsdb.Operation, 0)
	Address := fmt.Sprintf("%s %s",mac , InstanceIP	)
	SP := &externalmodel.SwitchPort{}
	CrOps,err:= SP.Create(o.Client, uuid,  "vif", Address, nil)
	if err != nil {
		return nil, fmt.Errorf("creating switch port error %v", err)
	}
	ops = append(ops, CrOps...)
	request := externalmodel.RequestControl{
		EXRList: o.ExternRouters,
		EXSList: o.ExternSwitchs,
		TargetUUID: SWUUID,
		Client: o.Client,
	}

	ConOps,err:=SP.Connect(request)
	if err != nil {
		return nil, fmt.Errorf("connecting switch port error %v", err)
	}
	ops = append(ops, ConOps...)


	result,err := o.Client.Transact(context.Background(),ops...)
	if err!=nil{
		fmt.Println("the problem is...", err)
	}
	fmt.Println(result)

	
	switchs:= o.ExternSwitchs.GetSwitch(SWUUID)
	VIF := externalmodel.StoVMPort{
		SwitchPort: SP,
		ConnectedSwitch: switchs,
	}

	externalmodel.AddNetInt(o.ExternRouters, InstanceIP, VIF)



	// util.SaveMapYaml(o.IPMapping)
	

	return SP,nil
}


func (o * Operator) AddSwitchAPort_Router(SWUUID string, lrpuuid string , uuid string)(*externalmodel.SwitchPort, error){
	ops:= make([]ovsdb.Operation, 0)
	SP := &externalmodel.SwitchPort{}
	CrOps,err := SP.Create(o.Client, uuid,  "router", "router", map[string]string{"router-port": lrpuuid})
	// client, uuid, portType, Address, router-options-string
	if err != nil {
		return nil,fmt.Errorf("creating switch port error %v", err)
	}

	ops = append(ops, CrOps...)

	request := externalmodel.RequestControl{
		EXRList: o.ExternRouters,
		EXSList: o.ExternSwitchs,
		TargetUUID: SWUUID,
		Client: o.Client,
	}


	ConOps,err := SP.Connect(request)
	if err != nil {
		return nil,fmt.Errorf("connecting switch port error %v", err)
	}

	ops = append(ops, ConOps...)

	result,err := o.Client.Transact(context.Background(),ops...)
	if err!=nil{
		fmt.Println("the problem is...", err)
	}
	fmt.Println(result)

	return SP,nil
}


// func (o * Operator) DelSwitchPort(ip string)(error){
// 	lspuuid := o.IPMapToDev(ip)
// 	if lspuuid == "" {
// 		return fmt.Errorf("no such switch port exist for ip %s", ip)
// 	}
// 	var lspDelOPS []ovsdb.Operation

// 	netAdd,err:=util.GetNetWorkInterface(ip)
// 	if err != nil {
// 		return fmt.Errorf("invalid IP format: %v %s %s", err, netAdd, ip)
// 	}
	
// 	dev, ok := o.IPMapping[netAdd+"1"]
// 	if !ok{
// 		return fmt.Errorf("no such switch port exist for ip %s", ip)
// 	}
	
// 	EXT,err := o.findDevByUUID(dev)
// 	if err != nil {
// 		return fmt.Errorf("no such device for uuid %s: %v", dev, err)
// 	}
// 	fmt.Println("EXT:", EXT)	
// 	s,ok := EXT.(*externalmodel.ExternSwitch)
// 		if !ok{
// 			return fmt.Errorf("no such switch exist for uuid %s", dev)
// 		}
// 		for j := 0; j < len(s.InternalSwitch.Ports); { 
// 			if s.InternalSwitch.Ports[j] == lspuuid {
// 				swMute ,_ := o.Client.Where(s.InternalSwitch).Mutate(s.InternalSwitch, model.Mutation{
// 				Field: &s.InternalSwitch.Ports,
// 				Mutator: ovsdb.MutateOperationDelete,
// 				Value: []string{lspuuid}, 

// 				})
// 				s.InternalSwitch.Ports = append(s.InternalSwitch.Ports[:j], s.InternalSwitch.Ports[j+1:]...)
// 				fmt.Println("swMute:", swMute)
// 				lspDelOPS = append(lspDelOPS, swMute...)
// 				break
// 						} 
// 			j++
// 	}
// 	fmt.Println("after processing EXT:", EXT)	
			
// 	// 포트 삭제
// 	lsp := &NBModel.LogicalSwitchPort{
// 		UUID: lspuuid,
// 	}
// 	lspDelOP, err := o.Client.Where(lsp).Delete()
// 	if err != nil {
// 		return fmt.Errorf("deleting switch port error %v", err)
// 	}
// 	lspDelOPS = append(lspDelOPS, lspDelOP...)

// 	result, err := o.Client.Transact(context.Background(), lspDelOPS...)
// 	if err != nil {
// 		return fmt.Errorf("deleting switch port error %v", err)
// 	}
// 	fmt.Println(result)

// 	delete(o.IPMapping, ip)
// 	util.SaveMapYaml(o.IPMapping)

// 	return nil

// }





func (o *Operator) AddSwitch () (string, error) {
	uuid ,err:=util.UUIDGenerator()
	if err!=nil{
		return "",fmt.Errorf("creating switch error %v",err)
	}
	newSwitch := externalmodel.ExternSwitch{}

	ops,err :=newSwitch.Create(o.Client, uuid.String())
	if err!=nil{
		return "",fmt.Errorf("creating switch error %v",err)
	}


	result , err := o.Client.Transact(context.Background(),ops...)
	if err!=nil{
		return uuid.String() ,fmt.Errorf("creating switch error %v",err)
	}
	fmt.Println(result)

	

	o.ExternSwitchs[uuid.String()]= &newSwitch
	return uuid.String(),nil

}


// func (o *Operator) DelSwitch(uuid string)(error){
// 	value,ok := o.ExternSwitchs[uuid]
// 	if !ok{
// 		return fmt.Errorf("no such switch exist")
// 	}
	
// 	lspDelOps:= make([]ovsdb.Operation, 0)


// 	for _,port := range value.InternalSwitch.Ports{
// 		lsp := &NBModel.LogicalSwitchPort{
// 			UUID: port,
// 		}
// 		lspDelOP, err := o.Client.Where(lsp).Delete()
// 		if err != nil {
// 			return fmt.Errorf("deleting switch port error %v", err)
// 		}
// 		lspDelOps = append(lspDelOps, lspDelOP...)

// 	}
// 	lsDelOp, err := o.Client.Where(value.InternalSwitch).Delete()
// 	if err != nil {
// 		return fmt.Errorf("deleting switch error %v", err)
// 	}
// 	lspDelOps = append(lspDelOps, lsDelOp...)

// 	result, err := o.Client.Transact(context.Background(), lspDelOps...)
// 	if err != nil {
// 		return fmt.Errorf("deleting switch port error %v", err)
// 	}
// 		fmt.Println(result)

// 	delete(o.ExternSwitchs, uuid)
// 	delete(o.IPMapping, value.IP)

// 	util.SaveMapYaml(o.IPMapping)

// 	return nil
// }