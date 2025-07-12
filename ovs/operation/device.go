package operation

import (
	"context"
	"fmt"

	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
)




func (o *Operator) DeleteSwitch(uuid string){
	Ds:= &NBModel.LogicalSwitch{
		UUID:uuid,
	}
	YO ,_:=o.Client.Where(Ds).Delete()
	fmt.Println(YO)
	o.Client.Transact(context.Background(), YO...)
}

func (o *Operator) DeleteAll(){
	DS:= &[]NBModel.LogicalSwitch{
	}
	DR:= &[]NBModel.LogicalSwitch{
	}

	o.Client.List(context.Background(),DS)
	for i:= range (*DS){
		o.DeleteSwitch((*DS)[i].UUID)
	}
	o.Client.List(context.Background(),DR)
	for i:= range (*DR){
		o.DeleteRouter((*DR)[i].UUID)
	}
	
}


func (o *Operator) DeleteRouter(uuid string){
	Ds:= &NBModel.LogicalRouter{
		UUID:uuid,
	}
	YO ,_:=o.Client.Where(Ds).Delete()
	fmt.Println(YO)
	o.Client.Transact(context.Background(), YO...)
}