package operation

import (
	"context"
	"fmt"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	"github.com/kwonkwonn/ovn-go-cms/ovs/util"
	"github.com/ovn-kubernetes/libovsdb/ovsdb"
)

func (o *Operator) AddSwitchAPort(SWUUID string, InstanceIP string, uuid string, mac string) (*externalmodel.SwitchPort, error) {
	ops := make([]ovsdb.Operation, 0)
	Address := fmt.Sprintf("%s %s", mac, InstanceIP)
	SP := &externalmodel.SwitchPort{}
	CrOps, err := SP.Create(o.Client, uuid, "vif", Address, nil)
	if err != nil {
		return nil, fmt.Errorf("creating switch port error %v", err)
	}
	ops = append(ops, CrOps...)
	request := externalmodel.RequestControl{
		EXRList:    o.ExternRouters,
		EXSList:    o.ExternSwitchs,
		TargetUUID: SWUUID,
		Client:     o.Client,
	}

	ConOps, err := SP.Connect(request)
	if err != nil {
		return nil, fmt.Errorf("connecting switch port error %v", err)
	}
	ops = append(ops, ConOps...)

	result, err := o.Client.Transact(context.Background(), ops...)
	if err != nil {
		fmt.Println("the problem is...", err)
	}
	fmt.Println(result)

	switchs := o.ExternSwitchs.GetSwitch(SWUUID)
	VIF := &externalmodel.StoVMPort{
		SwitchPort:      SP,
		ConnectedSwitch: switchs,
	}

	externalmodel.AddNetIntToRouter(o.ExternRouters[string(ROUTER)], InstanceIP, VIF)
	externalmodel.AddNetInt(o.ExternRouters, InstanceIP, VIF)

	return SP, nil
}

func (o *Operator) AddSwitchAPort_Router(SWUUID string, lrpuuid string, uuid string) (*externalmodel.SwitchPort, error) {
	ops := make([]ovsdb.Operation, 0)
	SP := &externalmodel.SwitchPort{}
	CrOps, err := SP.Create(o.Client, uuid, "router", "router", map[string]string{"router-port": lrpuuid})
	// client, uuid, portType, Address, router-options-string
	if err != nil {
		return nil, fmt.Errorf("creating switch port error %v", err)
	}

	ops = append(ops, CrOps...)

	request := externalmodel.RequestControl{
		EXRList:    o.ExternRouters,
		EXSList:    o.ExternSwitchs,
		TargetUUID: SWUUID,
		Client:     o.Client,
	}

	ConOps, err := SP.Connect(request)
	if err != nil {
		return nil, fmt.Errorf("connecting switch port error %v", err)
	}

	ops = append(ops, ConOps...)

	result, err := o.Client.Transact(context.Background(), ops...)
	if err != nil {
		fmt.Println("the problem is...", err)
	}
	fmt.Println(result)

	return SP, nil
}

func (o *Operator) DelSwitchPort(ip string) error {
	NetInt := externalmodel.GetNetInt(o.ExternRouters, ip)
	if len(NetInt) == 0 {
		return fmt.Errorf("no such switch port exist for ip %s", ip)
	}
	request := externalmodel.RequestControl{
		EXRList: o.ExternRouters,
		EXSList: o.ExternSwitchs,
		Client:  o.Client,
	}
	ops := make([]ovsdb.Operation, 0)
	for _, v := range NetInt {
		var delPort externalmodel.Deleter
		switch Port := v.(type) {
		case *externalmodel.RtoSwitchPort:
			lsUUID := Port.ConnectedSwitch.UUID
			request.TargetUUID = lsUUID
			delPort = Port.GetDeletor(externalmodel.SWITCH)

		case *externalmodel.StoVMPort:
			lsUUID := Port.ConnectedSwitch.UUID
			request.TargetUUID = lsUUID
			delPort = Port.GetDeletor(externalmodel.SWITCH)

		default:
			return fmt.Errorf("no such switch port exist for ip %s", ip)
		}

		operation, err := delPort.Delete(request)
		if err != nil {
			return fmt.Errorf("deleting switch port error %v", err)
		}
		ops = append(ops, operation...)
	}

	result, err := o.Client.Transact(context.Background(), ops...)
	if err != nil {
		return fmt.Errorf("deleting switch port error %v", err)
	}
	fmt.Println(result)

	return nil

}

func (o *Operator) AddSwitch() (string, error) {
	uuid, err := util.UUIDGenerator()
	if err != nil {
		return "", fmt.Errorf("creating switch error %v", err)
	}
	newSwitch := externalmodel.ExternSwitch{}

	ops, err := newSwitch.Create(o.Client, uuid.String())
	if err != nil {
		return "", fmt.Errorf("creating switch error %v", err)
	}

	result, err := o.Client.Transact(context.Background(), ops...)
	if err != nil {
		return uuid.String(), fmt.Errorf("creating switch error %v", err)
	}
	fmt.Println(result)

	o.ExternSwitchs[uuid.String()] = &newSwitch
	return uuid.String(), nil

}

func (o *Operator) DelSwitch(uuid string) error {
	value, ok := o.ExternSwitchs[uuid]
	if !ok {
		return fmt.Errorf("no such switch exist")
	}

	ops := make([]ovsdb.Operation, 0)

	lsDelOp, err := o.Client.Where(value.InternalSwitch).Delete()
	if err != nil {
		return fmt.Errorf("deleting switch error %v", err)
	}
	ops = append(ops, lsDelOp...)

	result, err := o.Client.Transact(context.Background(), ops...)
	if err != nil {
		return fmt.Errorf("deleting switch port error %v", err)
	}
	fmt.Println(result)

	delete(o.ExternSwitchs, uuid)

	return nil
}
