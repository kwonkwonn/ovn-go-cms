package operation

import (
	"context"
	"fmt"

	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
	"github.com/kwonkwonn/ovn-go-cms/ovs/util"
	"github.com/ovn-kubernetes/libovsdb/model"
	"github.com/ovn-kubernetes/libovsdb/ovsdb"
)

// AddSwitchACLs attaches isolation ACLs to a switch so that only same-subnet
// traffic and established outbound return traffic are allowed in.
// subnetGW is the router-port IP (e.g. "10.5.20.1"); subnet is derived from it.
func (o *Operator) AddSwitchACLs(swUUID string, subnetGW string) error {
	prefix, err := util.GetNetWorkSignifier(subnetGW)
	if err != nil {
		return fmt.Errorf("parsing subnet from %s: %w", subnetGW, err)
	}
	subnet := prefix + "0/24"

	controlPrefix, err := util.GetNetWorkSignifier(string(ROUTER))
	if err != nil {
		return fmt.Errorf("parsing control subnet from router IP: %w", err)
	}
	controlSubnet := controlPrefix + "0/24"

	sw := o.ExternSwitchs[swUUID]
	if sw == nil {
		return fmt.Errorf("switch %s not found in ExternSwitchs", swUUID)
	}

	rules := []struct {
		priority int
		match    string
		action   string
	}{
		{1000, fmt.Sprintf("ip4.src == %s", subnet), NBModel.ACLActionAllowRelated},
		// allow inbound from the control/proxy network (e.g. Traefik)
		{950, fmt.Sprintf("ip4.src == %s", controlSubnet), NBModel.ACLActionAllowRelated},
		{900, "ct.est", NBModel.ACLActionAllow},
		{800, "ct.rel", NBModel.ACLActionAllow},
		{500, "ip4", NBModel.ACLActionDrop},
	}

	ops := []ovsdb.Operation{}
	aclUUIDs := []string{}

	for _, rule := range rules {
		id, err := util.UUIDGenerator()
		if err != nil {
			return fmt.Errorf("generating ACL uuid: %w", err)
		}
		acl := &NBModel.ACL{
			UUID:      id.String(),
			Direction: NBModel.ACLDirectionToLport,
			Priority:  rule.priority,
			Match:     rule.match,
			Action:    rule.action,
		}
		createOps, err := o.Client.Create(acl)
		if err != nil {
			return fmt.Errorf("creating ACL (match=%q): %w", rule.match, err)
		}
		ops = append(ops, createOps...)
		aclUUIDs = append(aclUUIDs, id.String())
	}

	mutOps, err := o.Client.Where(sw.InternalSwitch).Mutate(sw.InternalSwitch, model.Mutation{
		Field:   &sw.InternalSwitch.ACLs,
		Mutator: ovsdb.MutateOperationInsert,
		Value:   aclUUIDs,
	})
	if err != nil {
		return fmt.Errorf("attaching ACLs to switch %s: %w", swUUID, err)
	}
	ops = append(ops, mutOps...)

	if _, err = o.Client.Transact(context.Background(), ops...); err != nil {
		return fmt.Errorf("ACL transact for switch %s: %w", swUUID, err)
	}
	return nil
}
