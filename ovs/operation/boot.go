package operation

import (
	"context"
	"fmt"
	"time"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
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
	}

	o.ExternSwitchs[LS.UUID]=exS

	return nil
	// switch메소드에 필요한 필드의 유무를 찾고 추가하는 함수를 넣을 예정
}