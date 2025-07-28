package externalmodel

import (
	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
	"github.com/ovn-kubernetes/libovsdb/client"
	"github.com/ovn-kubernetes/libovsdb/model"
	"github.com/ovn-kubernetes/libovsdb/ovsdb"
)






type RequestControl struct {
	EXRList  EXRList
	EXSList  EXSList
	TargetUUID string
	Client     client.Client
}


func (RP *RouterPort) Delete(request RequestControl) ([]ovsdb.Operation, error) {
	Router := request.EXRList.GetRouter(request.TargetUUID).InternalRouter
	Router = &NBModel.LogicalRouter{
		UUID: Router.UUID,
	}

	transactions, err := request.Client.Where(Router).Mutate(Router, model.Mutation{
		Field: &Router.Ports,
		Mutator: ovsdb.MutateOperationDelete,
		Value: []string{RP.UUID},
	})
	if err != nil {
		return nil, err
	}
	return transactions, nil
	


}


func (p RtoSwitchPort) GetDeletor(intType portType) Deleter {
	if intType == ROUTER {
		return p.RouterPort
	} else if intType == SWITCH {
		return p.SwitchPort
	}
	return nil
}

func (p StoVMPort) GetDeletor(intType portType) Deleter {
	if intType == "switch" {
		return p.SwitchPort
	}
	return nil
}