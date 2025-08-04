package externalmodel

import (
	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
	"github.com/kwonkwonn/ovn-go-cms/ovs/util"
	"github.com/ovn-kubernetes/libovsdb/client"
	"github.com/ovn-kubernetes/libovsdb/ovsdb"
)

func (ES *ExternSwitch) Create(client client.Client, uuid string) ([]ovsdb.Operation, error) {
	ES.UUID = uuid
	ES.InternalSwitch = &NBModel.LogicalSwitch{}
	ES.InternalSwitch.UUID = uuid
	ES.InternalSwitch.Name = uuid
	ES.InternalSwitch.Ports = []string{}

	transactions, err := client.Create(ES.InternalSwitch)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (SP *SwitchPort) Create(client client.Client, uuid string,  portType string, Address string, Options map[string]string) ([]ovsdb.Operation, error) {
	SP.UUID = uuid
	SP.Name = uuid
	SP.Addresses = []string{Address}

	if portType == "" || portType == "vif" {
		SP.Type = ""
	}else{
		SP.Type = portType // "vif" or "router"
	}
	

	if Options != nil {
		SP.Options = make(map[string]string)
		for k,v := range Options {
			SP.Options[k] = v
		}
	}
	Internal := NBModel.LogicalSwitchPort(*SP)
	transactions, err := client.Create(&Internal)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (R *ExternRouter) Create(client client.Client, uuid string) ([]ovsdb.Operation, error) {
	R.UUID = uuid
	R.SubNetworks = make(map[string]NetInt)
	R.InternalRouter = &NBModel.LogicalRouter{}
	R.InternalRouter.UUID = uuid
	R.InternalRouter.Name = uuid
	R.InternalRouter.Ports = []string{}

	transactions, err := client.Create(R.InternalRouter)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}


func (RP *RouterPort)Create(client client.Client,uuid string, ip string) ([]ovsdb.Operation, error) {
    mac,_:= util.MacGenerator()

	RP.UUID = uuid
	RP.Name = uuid
	RP.MAC = mac
	RP.Networks = []string{ip + "/24"}
	
	internal := NBModel.LogicalRouterPort(*RP)

	transactions, err := client.Create(&internal	)
	if err != nil {
		return nil, err
	}
	return transactions, nil

}
	
