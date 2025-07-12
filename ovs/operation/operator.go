package operation

import (
	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	"github.com/ovn-org/libovsdb/client"
)





type Operator struct{
	Client client.Client
	ExternRouters map[string]*externalmodel.ExternRouter
	ExternSwitchs map[string]*externalmodel.ExternSwitch
	IPMapping map[string]string // device uuid
}