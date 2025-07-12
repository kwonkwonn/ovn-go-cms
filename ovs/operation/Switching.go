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


func (o * Operator) AddSwitchAPort(SWUUID string, ip string )(error){
	value ,ok := o.ExternSwitchs[SWUUID]; 
	if !ok{
		return fmt.Errorf("no such switch exist")
	}

	
	uuid,err := util.UUIDGenerator()
	if err!=nil{
		return fmt.Errorf("no such switch exist")
	}
	mac,err := util.MacGenerator()
	if err!=nil{
		return fmt.Errorf("no such switch exist")
	}

	newSP:= &NBModel.LogicalSwitchPort{
		UUID: string(uuid.String()),
		}

	newSP.Addresses=append(newSP.Addresses, mac)

	lsp , err := o.Client.Create(newSP)
	if err!=nil{
		return fmt.Errorf("%v", err)
	}
	
	value.InternalSwitch.Ports = append(value.InternalSwitch.Ports, newSP.UUID)
	lsMute,_  := o.Client.Where(value.InternalSwitch).Mutate( value.InternalSwitch, model.Mutation{
		Field: value.InternalSwitch.Ports,
		Mutator: ovsdb.MutateOperationAdd,
		Value: value.InternalSwitch.Ports,	
	}) 
	lsp = append(lsp, lsMute...)
	o.Client.Transact(context.Background(),lsp...)

	return nil
}

func (o * Operator) AddSwitch () (string ,error){
	return  o.addSwitch()
	
}


func (o *Operator) addSwitch () (string, error) {
	uuid ,err:=util.UUIDGenerator()
	if err!=nil{
		return "",fmt.Errorf("creating switch error %v",err)
	}
	newSwitch:=&NBModel.LogicalSwitch{
		UUID:uuid.String(),
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

	o.ExternSwitchs[uuid.String()]= &externalmodel.ExternSwitch{
		UUID: uuid.String(),
		InternalSwitch: newSwitch,
	}
	return uuid.String(),nil

}