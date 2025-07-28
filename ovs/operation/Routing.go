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

//lrpuuid string
func (o *Operator) DelRouterPort(network string)(error){
    lruuid:=o.IPMapToDev(string(ROUTER))
    lrp:= o.IPMapToDev()
    connectedRouter:= o.ExternRouters.GetRouter(lruuid)
    if connectedRouter == nil {
        return fmt.Errorf("no such router exist")
    }
    ops:= make([]ovsdb.Operation, 0)
    
    
    
    
    operation,err :=o.Client.Where(connectedRouter.InternalRouter).Mutate(connectedRouter.InternalRouter, model.Mutation{
        Field: &connectedRouter.InternalRouter.Nat,
        Mutator: ovsdb.MutateOperationDelete,
        Value: []string{netInt+"/24"},   
    })    
    if err != nil {
        return fmt.Errorf("error mutating router nat: %v", err)
    }
    ops = append(ops, operation...)

    nat:= &NBModel.NAT{
        Type: "snat",
        LogicalIP: netInt + "/24",
        ExternalIP: connectedRouter.IP,
    }
    operation,err = o.Client.WhereAny(nat).Delete()
    if err != nil {
        return fmt.Errorf("error deleting router nat: %v", err)
    }

    ops = append(ops, operation...)

    result, err:= o.Client.Transact(context.Background(), ops...)
    if err != nil {
        return fmt.Errorf("deleting router port transaction error: %v, result: %+v", err, result)
    }
    fmt.Println("DelRouterPort Transact Result:", result)

    return nil
}

func (o *Operator) AddRouterPort(lruuid string ,lrpuuid string, ip string)(*externalmodel.RouterPort,error){
        
        operations := make([]ovsdb.Operation, 0)
        newRP := externalmodel.RouterPort{}


        ops,err:= newRP.Create(o.Client,lrpuuid, ip)
        if err != nil {
            return nil,fmt.Errorf("creating logical router port error: %v", err)
        }
        operations=append(operations, ops...)

        request := externalmodel.RequestControl{
            Client: o.Client,
            EXRList: o.ExternRouters,
            EXSList: o.ExternSwitchs,
            TargetUUID: lruuid,
        }
        ops , err= newRP.Connect(request)
        if err != nil {
            return nil,fmt.Errorf("connecting router port error: %v", err)
        }

        operations=append(operations, ops...)


        result, err:= o.Client.Transact(context.Background(),ops...)
        if err!=nil{
            return nil,fmt.Errorf("creating logical router port transaction error: %v, result: %+v",err, result)
        }

        return &newRP, nil
}


func (o *Operator) AddRouter(IP string) (string, error) {
    //보통 라우터는 IP 가 닉네임으로 지정되어 있음. operator 참조
	RtUUID, err := util.UUIDGenerator()
	if err != nil {
		return "", fmt.Errorf("generating error: transaction logical router %v", err)
	}
    router := externalmodel.ExternRouter{}

	createOP,err:= router.Create(o.Client, RtUUID.String())
    if err != nil {
		return "", fmt.Errorf("creating operations for Router failed %v", err)
	}   


	result,err:=o.Client.Transact(context.Background(),createOP...)
	if err!=nil{
		return "",fmt.Errorf("creaing operations for Router failed %v",err)
	}
	fmt.Println(result)
	
	o.ExternRouters[RtUUID.String()]=&router
	o.IPMapping[IP]=RtUUID.String()


	util.SaveMapYaml(o.IPMapping)
    // IPMapping에 IP와 UUID를 저장
	return RtUUID.String(),nil

}