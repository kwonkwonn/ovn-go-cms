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


func (o*Operator) AddRouterPort(lruuid string ,lrpuuid string, ip string)(error){
	//새로운 서브넷을 추가하는 것과 동일한 기능
	mac,err:=util.MacGenerator()
	if err!=nil{
		return fmt.Errorf("generating mac for router port error %v", err)
	}

	externRouter,ok := o.ExternRouters[lruuid]
	if !ok{
		return fmt.Errorf("router not exist")
	}
    var latestRouter =NBModel.LogicalRouter{ UUID: externRouter.UUID}
    err = o.Client.Get(context.Background(), &latestRouter)
    if err != nil {
        return fmt.Errorf("failed to get latest LogicalRouter from DB for UUID %s: %v", externRouter.UUID, err)
    }

	newRP:=&NBModel.LogicalRouterPort{
		UUID: lrpuuid,
		Name: lrpuuid,
		MAC: mac,
		Networks: []string{ip},
	}
	ops,err:= o.Client.Create(context.Background(),newRP)
	if err!=nil{
		return fmt.Errorf("creating logical router error %v",err)
	}
	lsMute,_  := o.Client.Where(&latestRouter).Mutate( latestRouter, model.Mutation{
		Field: &latestRouter.Ports,
		Mutator: ovsdb.MutateOperationInsert,
		Value: []string{newRP.UUID},
	})
	ops=append(ops, lsMute...)
	result, err:= o.Client.Transact(context.Background(),ops...)
	if err!=nil{
		return fmt.Errorf("creating logical router error: transaction error %v",err)
	}
	fmt.Println(result)

	return nil
}



func (o * Operator)AddRouter( IP string) (string, error){
	RtUUID , err:= util.UUIDGenerator()
	if err!=nil{
		return "",fmt.Errorf("generating error: transaction logical rotuer %v",err)
	}
	newR:= &NBModel.LogicalRouter{
		UUID: RtUUID.String(),
		Name: RtUUID.String(),
	}

	ops, err:= o.Client.Create(newR)
	if err!=nil{
		return "",fmt.Errorf("creaing operations for Router failed %v",err)
	}
	result,err:=o.Client.Transact(context.Background(),ops...)
	if err!=nil{
		return "",fmt.Errorf("creaing operations for Router failed %v",err)
	}
	fmt.Println(result)
	
	o.ExternRouters[RtUUID.String()]=&externalmodel.ExternRouter{
		InternalRouter: newR,
		IP: IP,
		UUID: RtUUID.String(),
	}
	o.IPMapping[IP]=RtUUID.String()



	return RtUUID.String(),nil

}