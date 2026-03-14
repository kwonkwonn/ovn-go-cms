package externalmodel

import (
	"fmt"

	"github.com/ovn-kubernetes/libovsdb/model"
	"github.com/ovn-kubernetes/libovsdb/ovsdb"
)

type RouterPort_test struct {

}

func (LP *RouterPort_test) Connect(request RequestControl) ([]ovsdb.Operation, error) {

	return transactions, nil
}

//

func (sp SwitchPort) Connect(request RequestControl) ([]ovsdb.Operation, error) {
	targetSwitch := request.EXSList.GetSwitch(request.TargetUUID)
	if targetSwitch == nil {
		return nil, fmt.Errorf("no such switch exist")
	}

	targetSwitch.InternalSwitch.Ports = append(targetSwitch.InternalSwitch.Ports, sp.UUID)
	lsMute, err := request.Client.Where(targetSwitch.InternalSwitch).Mutate(targetSwitch.InternalSwitch, model.Mutation{
		Field:   &targetSwitch.InternalSwitch.Ports,
		Mutator: ovsdb.MutateOperationInsert,
		Value:   []string{sp.UUID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to mutate switch ports: %w", err)
	}
	return lsMute, nil
}

func (p RtoSwitchPort) GetConnector(intType portType) Connector {
	if intType == ROUTER {
		return p.RouterPort
	} else if intType == SWITCH {
		return p.SwitchPort
	}
	return nil
}

func (p StoVMPort) GetConnector(intType portType) Connector {
	if intType == "switch" {
		return p.SwitchPort
	}
	return nil
}
