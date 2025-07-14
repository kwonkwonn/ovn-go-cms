package operation

import (
	"context"
	"fmt"
	"time"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
	"github.com/kwonkwonn/ovn-go-cms/ovs/util"
	"github.com/ovn-org/libovsdb/model"
	"github.com/ovn-org/libovsdb/ovsdb"
)


func (o* Operator)InitializeLogicalDevices (){
	o.ExternRouters = make(map[string]*externalmodel.ExternRouter)
	o.ExternSwitchs = make(map[string]*externalmodel.ExternSwitch)

	LR :=&[]NBModel.LogicalRouter{}
	LS :=&[]NBModel.LogicalSwitch{}

	err:= o.Client.List(context.Background(), LS )
	if err!=nil{
		fmt.Println(fmt.Errorf("%v", err))
	}
	err= o.Client.List(context.Background(), LR)
	if err!=nil{
		fmt.Println(fmt.Errorf("%v", err))
	}
	time.Sleep(2 * time.Second) // 2초 대기
	for i:=range *LR{
		o.AddExternRouter((*LR)[i])
	}
	for i:=range *LS{
		o.AddExternSwitch((*LS)[i])
	}
}

func (o* Operator)AddExternRouter (LR NBModel.LogicalRouter)error {
	exR:= &externalmodel.ExternRouter{
		UUID:LR.UUID,
		InternalRouter: &LR,
	}

	o.ExternRouters[LR.UUID] = exR
	// if len(exR.InternalRouter.Ports)!=0{
	// 	ports:= &[]NBModel.LogicalRouterPort{}
	// 	o.Client.List(context.Background(),ports)

	// }
	return nil
}

func (o* Operator)AddExternSwitch (LS NBModel.LogicalSwitch) error{
	exS:=&externalmodel.ExternSwitch{
		UUID: LS.UUID,
		//IP: yaml에서 읽어서 할당
	}
	o.IPMapping[exS.IP] = exS.UUID 
	o.ExternSwitchs[LS.UUID]=exS

	return nil
	// switch메소드에 필요한 필드의 유무를 찾고 추가하는 함수를 넣을 예정
}


func (o *Operator) AddInterconnectR_S(lsUUID string, lrUUID string, ip string)(error){
	lrpuuid,err:=util.UUIDGenerator()
	if err!=nil{
		panic("lrpuuid generating error" )
	}
	lspuuid,err:=util.UUIDGenerator()
	if err!=nil{
		panic("lrpuuid generating error" )
	}
	
	o.AddSwitchAPort_Router(lsUUID, lspuuid.String(), lspuuid.String())
	o.AddRouterPort(lrUUID, lrpuuid.String(),ip)


	return nil
}


func (o* Operator) InitialSettig()(error){
	//ls-ext 를 만듬
	//lr-ext를 만듬
	//lr-ext와 연결, lr-extsms 10.5.15.4에 저장되어 새로운 스위치가 생길때마다 연결해줘야 함

	EXTS_uuid,err:= o.AddSwitch("EXT_S")
	if (err!=nil){
		panic("bootstraping failed, creating external Switch")
	}
	EXTR_uuid,err:=o.AddRouter(string(ROUTER))
	if (err!=nil){
		panic("bootstraping failed, creating external Switch")
	}


	o.AddInterconnectR_S(EXTS_uuid, EXTR_uuid, string(ROUTER)) // 새로운 가상 라우터는 10.5.15.4 할당 받음
	
	br_EXTS_UUID,err := util.UUIDGenerator()
	if err!=nil{
		return fmt.Errorf("generating uuid error: br_exts_uuid")
	}
	
	//Addswitch port 함수의 복사본, 함수화가 아직 부족해서 복붙함
	// 리팩토링 대상 1번 
	value ,ok := o.ExternSwitchs[EXTS_uuid]; 
	if !ok{
		return fmt.Errorf("no such switch exist")
	}

	newSP:= &NBModel.LogicalSwitchPort{
		UUID: br_EXTS_UUID.String(),
		Name: br_EXTS_UUID.String(),
		Type: "localnet",
		Options: map[string]string{"network_name":"UPLINK"},
		}
	Address :="unknown"
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
	o.IPMapping[string(UPLINK)]= newSP.UUID
	 
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