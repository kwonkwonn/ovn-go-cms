package initialize

import (
	_ "github.com/ovn-kubernetes/libovsdb/client"
	model "github.com/ovn-kubernetes/libovsdb/model"
)



var MainDBModel *model.ClientDBModel


type Logical_Switch struct {
    UUID   string            `ovsdb:"_uuid"` // _uuid tag is mandatory
    Name   string            `ovsdb:"name"`
    Ports  []string          `ovsdb:"ports"`
    Config map[string]string `ovsdb:"other_config"`
}



