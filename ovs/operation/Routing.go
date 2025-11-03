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

// lrpuuid string
func (o *Operator) DelRouterPort(ip string) error {
	NetInt := externalmodel.GetNetInt(o.ExternRouters, ip)
	if len(NetInt) == 0 {
		return fmt.Errorf("no such router port exist for ip %s", ip)
	}

	Port, ok := NetInt[0].(*externalmodel.RtoSwitchPort)
	if !ok {
		return fmt.Errorf("the netInt is not a router port: %v", NetInt[0])
	}

	request := externalmodel.RequestControl{
		Client:     o.Client,
		EXRList:    o.ExternRouters,
		EXSList:    o.ExternSwitchs,
		TargetUUID: Port.ConnectedRouter.UUID,
	}

	lrp := NetInt[0].GetDeletor(externalmodel.ROUTER)
	if lrp == nil {
		return fmt.Errorf("no such router port exist for ip %s", ip)
	}

	ops := make([]ovsdb.Operation, 0)

	operation, err := lrp.Delete(request)
	if err != nil {
		return fmt.Errorf("deleting router port error %v", err)
	}
	ops = append(ops, operation...)

	nat := &NBModel.NAT{
		UUID: Port.NatConnected,
	}
	operation, err = o.Client.Where(nat).Delete()
	if err != nil {
		return fmt.Errorf("error deleting router nat: %v", err)
	}

	ops = append(ops, operation...)

	operation, err = o.DeleteNAT(Port.ConnectedRouter.UUID, Port.NatConnected)
	if err != nil {
		fmt.Printf("AddInterconnectR_S ERROR: deleting NAT from router error %v\n", err)
		return err
	}
	ops = append(ops, operation...)

	result, err := o.Client.Transact(context.Background(), ops...)
	if err != nil {
		return fmt.Errorf("deleting router port transaction error: %v, result: %+v", err, result)
	}
	fmt.Println("DelRouterPort Transact Result:", result)

	return nil
}

func (o *Operator) AddRouterPort(lruuid string, lrpuuid string, natuuid string, ip string) (*externalmodel.RouterPort, error) {

	operations := make([]ovsdb.Operation, 0)
	newRP := externalmodel.RouterPort{}

	ops, err := newRP.Create(o.Client, lrpuuid, ip)
	if err != nil {
		return nil, fmt.Errorf("creating logical router port error: %v", err)
	}
	operations = append(operations, ops...)

	request := externalmodel.RequestControl{
		Client:     o.Client,
		EXRList:    o.ExternRouters,
		EXSList:    o.ExternSwitchs,
		TargetUUID: lruuid,
	}
	ops, err = newRP.Connect(request)
	if err != nil {
		return nil, fmt.Errorf("connecting router port error: %v", err)
	}

	operations = append(operations, ops...)

	fmt.Println("IPADDRESS RECEIVED:", ip)
	if (ip) != string(ROUTER) {

		netmask, err := util.GetNetWorkSignifier(ip)
		if err != nil {
			fmt.Printf("AddInterconnectR_S ERROR: error getting network address: %v\n", err)
			return nil, err
		}

		fmt.Println("AddInterconnectR_S: netmask:", netmask)

		newNat := &NBModel.NAT{
			UUID:       natuuid,
			Type:       NBModel.NATTypeSNAT,
			ExternalIP: string(ROUTER),
			LogicalIP:  netmask + "0/24",
		}
		operation, err := o.Client.Create(newNat)
		if err != nil {
			fmt.Printf("AddInterconnectR_S ERROR: creating NAT error %v\n", err)
			return nil, err
		}
		operations = append(operations, operation...)

		ops, err = o.AddNAT(lruuid, natuuid)
		if err != nil {
			fmt.Printf("AddInterconnectR_S ERROR: adding NAT to router error %v\n", err)
			return nil, err
		}
	}
	operations = append(operations, ops...)

	result, err := o.Client.Transact(context.Background(), operations...)
	if err != nil {
		return nil, fmt.Errorf("creating logical router port transaction error: %v, result: %+v", err, result)
	}

	return &newRP, nil
}

func (o *Operator) AddRouter(IP string) (string, error) {
	//보통 라우터는 IP 가 닉네임으로 지정되어 있음. operator 참조
	RtUUID, err := util.UUIDGenerator()
	if err != nil {
		return "", fmt.Errorf("generating error: transaction logical router %v", err)
	}
	router := externalmodel.ExternRouter{}

	createOP, err := router.Create(o.Client, RtUUID.String())
	if err != nil {
		return "", fmt.Errorf("creating operations for Router failed %v", err)
	}

	result, err := o.Client.Transact(context.Background(), createOP...)
	if err != nil {
		return "", fmt.Errorf("creaing operations for Router failed %v", err)
	}
	fmt.Println(result)

	o.ExternRouters[IP] = &router
	o.ExternRouters[RtUUID.String()] = &router

	return RtUUID.String(), nil

}

func (o *Operator) AddNAT(routerUUID string, natuuid string) ([]ovsdb.Operation, error) {
	Router := o.ExternRouters.GetRouter(routerUUID)

	newRouter := &NBModel.LogicalRouter{
		UUID: Router.InternalRouter.UUID,
	}

	transaction, err := o.Client.Where(newRouter).Mutate(newRouter, model.Mutation{
		Field:   &newRouter.Nat,
		Mutator: ovsdb.MutateOperationInsert,
		Value:   []string{natuuid},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create mutate operation for router nat: %w", err)
	}

	return transaction, nil
}

func (o *Operator) DeleteNAT(routerUUID string, natuuid string) ([]ovsdb.Operation, error) {
	Router := o.ExternRouters.GetRouter(routerUUID)

	newRouter := &NBModel.LogicalRouter{
		UUID: Router.InternalRouter.UUID,
	}

	transaction, err := o.Client.Where(newRouter).Mutate(newRouter, model.Mutation{
		Field:   &newRouter.Nat,
		Mutator: ovsdb.MutateOperationDelete,
		Value:   []string{natuuid},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create mutate operation for router nat: %w", err)
	}

	return transaction, nil
}
