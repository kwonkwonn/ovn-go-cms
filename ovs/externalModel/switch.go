package externalmodel

import (
	"fmt"

	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
	"github.com/kwonkwonn/ovn-go-cms/ovs/util"
)






func NewSwitch(ip string)(*ExternSwitch,error){
	uuid ,err:=util.UUIDGenerator()
	if err!=nil{
		return nil,fmt.Errorf("creating switch error:  %v",err)
	}
	newSwitch:=&ExternSwitch{
		UUID: uuid.String(),
		InternalSwitch: &NBModel.LogicalSwitch{
		Name: uuid.String(),
	},
}
	return newSwitch,nil
}




