package main

import (
	"log"

	initialize "github.com/kwonkwonn/ovn-go-cms/initialize"
	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	"github.com/kwonkwonn/ovn-go-cms/ovs/operation"
	"github.com/kwonkwonn/ovn-go-cms/ovs/util"
	"github.com/kwonkwonn/ovn-go-cms/server"
	"github.com/kwonkwonn/ovn-go-cms/service"
)

const NB_DB string = "10.5.15.3"

func main(){

	
	ovnClient, err := initialize.InitializeOvnClient(NB_DB)
	if err != nil {
		log.Fatalf("Failed to initialize OVN client: %v", err)
	}

	Operator := &operation.Operator{
		Client: ovnClient,
	}
	Operator.ExternRouters = make(map[string]*externalmodel.ExternRouter, 0)
	Operator.ExternSwitchs = make(map[string]*externalmodel.ExternSwitch,0)

	Operator.IPMapping = make(map[string]string) // 항상 초기화
	err = util.ReadMapNetYaml(Operator.IPMapping)
	if err != nil {
		log.Printf("IPMapping YAML read error: %v", err)
	}
	Operator.InitializeLogicalDevices()
	if _,ok:=Operator.IPMapping["EXT_S"] ;!ok {
		Operator.InitialSettig()
	}
	// Chassis 초기화

	handler:=service.Handler{
		Operator: Operator,
	}
	server.InitServer(8081,handler)

	
	
	select{}
}

func init() {

}
