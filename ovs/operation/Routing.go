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
func (o *Operator) DelRouterPort(netInt string)(error){
    // ports:= &NBModel.LogicalRouterPort{
    //     // Networks: []string{"20.20.22.1"+ "/24"},
    // }
    portss:= &[]NBModel.LogicalRouter{}
     o.Client.List(context.Background(),portss)

    fmt.Println("asfsfidsaofbadsfoisabfsodifsbfa", portss)

    lruuid:=o.IPMapToDev(string(ROUTER))
    connectedRouter:= o.ExternRouters[lruuid]
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

func (o *Operator) AddRouterPort(lruuid string ,lrpuuid string, ip string)(error){
    mac,_:= util.MacGenerator()
    
    newRP:=&NBModel.LogicalRouterPort{
        	UUID: lrpuuid,
        	Name: lrpuuid,
        	MAC: mac,
        	Networks: []string{ip+"/24"}, // 이미 초기화되어 있음
        }
        ops,err:= o.Client.Create(newRP)
        if err!=nil{
            return fmt.Errorf("creating logical router port create operation failed for newRP: %+v, error: %w", newRP, err)
        }

        latestRouterModel := o.ExternRouters[lruuid].InternalRouter

        lsMute, muteErr := o.Client.Where(latestRouterModel).Mutate(latestRouterModel, model.Mutation{
            Field: &latestRouterModel.Ports,  
            Mutator: ovsdb.MutateOperationInsert,
            Value: []string{lrpuuid},  
        })
        if muteErr != nil {
            return fmt.Errorf("failed to create mutate operation for router ports: %w", muteErr)
        }


        ops=append(ops, lsMute...)
        fmt.Printf("AddRouterPort: Total operations in transaction: %d\n", len(ops))

        result, err:= o.Client.Transact(context.Background(),ops...)
        if err!=nil{
            return fmt.Errorf("creating logical router port transaction error: %v, result: %+v",err, result)
        }
        fmt.Println("AddRouterPort Transact Result:", result)
        for i, res := range result {
            fmt.Printf("  Operation %d: UUID=%v\n", i, res.UUID) 
            // res.Err이 없다면 이렇게 확인
        // }
        }
        return nil
}




// func (o*Operator) AddRouterPort(lruuid string ,lrpuuid string, ip string)(error){
//     mac, err := util.MacGenerator()
//     if err != nil {
//         return fmt.Errorf("generating mac for router port error %v", err)
//     }

//     // `ovn-nbctl`의 절대 경로를 사용
//     // creatR := fmt.Sprintf("sudo ovn-nbctl lrp-add %s %s %s %s", lruuid, lrpuuid, mac, ip) // 기존 코드
    
//     // 변경된 코드 (예시: /usr/bin/ovn-nbctl 에 설치된 경우)
//     command := "ovn-nbctl" 
//     args := []string{
//         "lrp-add",
//         lruuid,
//         lrpuuid,
//         mac,
//         ip+"/24",
//     }

//     cmd := exec.Command(command, args...) // `exec.Command`는 명령어와 인자를 분리해서 받는 것이 더 안전합니다.
//     err = cmd.Run()
//     if err != nil {
//         return fmt.Errorf("error creating router command, %v", err)
//     }

//     return nil
// }


func (o *Operator) AddRouter(IP string) (string, error) {
	RtUUID, err := util.UUIDGenerator()
	if err != nil {
		return "", fmt.Errorf("generating error: transaction logical router %v", err)
	}
	newR := &NBModel.LogicalRouter{
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


	util.SaveMapYaml(o.IPMapping)

	return RtUUID.String(),nil

}