package externalmodel

import (
	"fmt"

	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
	"github.com/ovn-kubernetes/libovsdb/client"
	"github.com/ovn-kubernetes/libovsdb/model"
	"github.com/ovn-kubernetes/libovsdb/ovsdb"
)

type RequestControl struct {
	EXRList    EXRList
	EXSList    EXSList
	TargetUUID string
	Client     client.Client
}

func (RP *RouterPort) Delete(request RequestControl) ([]ovsdb.Operation, error) {
	Router := request.EXRList.GetRouter(request.TargetUUID).InternalRouter

	Router = &NBModel.LogicalRouter{
		UUID: Router.UUID,
	}
	ops := []ovsdb.Operation{}

	lrp := NBModel.LogicalRouterPort(*RP)
	transactions, err := request.Client.Where(&lrp).Delete()
	if err != nil {
		return nil, fmt.Errorf("failed to create delete operation for router port: %w", err)
	}
	ops = append(ops, transactions...)

	transactions, err = request.Client.Where(Router).Mutate(Router, model.Mutation{
		Field:   &Router.Ports,
		Mutator: ovsdb.MutateOperationDelete,
		Value:   []string{RP.UUID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create mutate operation for router ports: %w", err)
	}
	ops = append(ops, transactions...)

	return ops, nil

}

func (SP *SwitchPort) Delete(request RequestControl) ([]ovsdb.Operation, error) {
	targetSwitch := request.EXSList.GetSwitch(request.TargetUUID)
	if targetSwitch == nil {
		return nil, nil // No switch to delete from
	}
	ops := []ovsdb.Operation{}

	lsp := NBModel.LogicalSwitchPort(*SP)
	transaction, err := request.Client.Where(&lsp).Delete()
	if err != nil {
		return nil, fmt.Errorf("failed to create delete operation for switch port: %w", err)
	}
	ops = append(ops, transaction...)

	transaction, err = request.Client.Where(targetSwitch.InternalSwitch).Mutate(targetSwitch.InternalSwitch, model.Mutation{
		Field:   &targetSwitch.InternalSwitch.Ports,
		Mutator: ovsdb.MutateOperationDelete,
		Value:   []string{SP.UUID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create mutate operation for switch ports: %w", err)
	}
	ops = append(ops, transaction...)
	return ops, nil
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
	if intType == SWITCH {
		return p.SwitchPort
	}
	return nil
}
