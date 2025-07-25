package operation

import (
	"context"
	"fmt"
	"os/exec"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
	"github.com/kwonkwonn/ovn-go-cms/ovs/util"
)

//lrpuuid string
func (o *Operator) DelRouterPort(netInt string)(error){
    // `ovn-nbctl`의 절대 경로를 사용
    // delR := fmt.Sprintf("sudo ovn-nbctl lrp-del %s %s", lruuid, lrpuuid) // 기존 코드
    lruuid:=o.CheckIPExistance(string(ROUTER))
    // 변경된 코드 (예시: /usr/bin/ovn-nbctl 에 설치된 경우)
    // command := "ovn-nbctl"
    // args := []string{
    //     "lrp-del",
    //     lruuid,
    //     lrpuuid,
    // }
    // 라우터 포트 구조체 문제 때문에, 당장 수정은 불가능
    // 추후에 라우터 포트 구조체를 수정하고, 그에 맞게 수정해야 함
    // cmd := exec.Command(command, args...) // `exec.Command`는 명령어와 인자를 분리해서 받는 것이 더 안전합니다.
    // err := cmd.Run()
    // if err != nil {
    //     return fmt.Errorf("error deleting router port command, %v", err)
    // }
    command:= "ovn-nbctl"
    args:= []string{
        "lr-nat-del",
        lruuid,
        "snat",
       netInt+"/24",
    }

   cmd := exec.Command(command, args...) // `exec.Command`는 명령어와 인자를 분리해서 받는 것이 더 안전합니다.
    err := cmd.Run()
    if err != nil {
        return fmt.Errorf("error deleting router port command, %v", err)
    }

    return nil
}


func (o*Operator) AddRouterPort(lruuid string ,lrpuuid string, ip string)(error){
    mac, err := util.MacGenerator()
    if err != nil {
        return fmt.Errorf("generating mac for router port error %v", err)
    }

    // `ovn-nbctl`의 절대 경로를 사용
    // creatR := fmt.Sprintf("sudo ovn-nbctl lrp-add %s %s %s %s", lruuid, lrpuuid, mac, ip) // 기존 코드
    
    // 변경된 코드 (예시: /usr/bin/ovn-nbctl 에 설치된 경우)
    command := "ovn-nbctl" 
    args := []string{
        "lrp-add",
        lruuid,
        lrpuuid,
        mac,
        ip+"/24",
    }

    cmd := exec.Command(command, args...) // `exec.Command`는 명령어와 인자를 분리해서 받는 것이 더 안전합니다.
    err = cmd.Run()
    if err != nil {
        return fmt.Errorf("error creating router command, %v", err)
    }



	// newRP:=&NBModel.LogicalRouterPort{
	// 	UUID: lrpuuid,
	// 	Name: lrpuuid,
	// 	MAC: mac,
	// 	Networks: []string{ip+"/24"}, // 이미 초기화되어 있음
	// }
    // ops,err:= o.Client.Create(context.Background(),newRP)
    // if err!=nil{
    //     return fmt.Errorf("creating logical router port create operation failed for newRP: %+v, error: %w", newRP, err)
    // }
    // lsMute, muteErr := o.Client.Where(&latestRouterModel).Mutate(&latestRouterModel, model.Mutation{
    //     Field: &latestRouterModel.Ports, // latestRouterModel의 Ports 필드 참조
    //     Mutator: ovsdb.MutateOperationInsert,
    //     Value: []string{lrpuuid}, // OVSDB Set 타입에 맞게 []string 전달 (model.NewSet이 없다면)
    // })
    // if muteErr != nil {
    //     return fmt.Errorf("failed to create mutate operation for router ports: %w", muteErr)
    // }
    ///// <---- 이 ship 놈이 스위치 포트 추가에서는 정상적으로 동작하는데 라우터 포트에서는 어떤 이유인지 동작안함
    // 시간이 부족해서 일단 임시방편으로 추가 후, 나중에 gdb등 사용해서 디버깅할 예정..(이미 하루 소모)


    // ops=append(ops, lsMute...)
    // fmt.Printf("AddRouterPort: Total operations in transaction: %d\n", len(ops))

    // result, err:= o.Client.Transact(context.Background(),ops...)
    // if err!=nil{
    //     return fmt.Errorf("creating logical router port transaction error: %v, result: %+v",err, result)
    // }
    // fmt.Println("AddRouterPort Transact Result:", result)
    // for i, res := range result {
        // fmt.Printf("  Operation %d: UUID=%v\n", i, res.UUID) // res.Err이 없다면 이렇게 확인
    // }

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


	util.SaveMapYaml(o.IPMapping)

	return RtUUID.String(),nil

}