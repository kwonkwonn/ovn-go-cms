package main

import (
	"fmt"
	"log"

	initialize "github.com/kwonkwonn/ovn-go-cms/initialize"
	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	"github.com/kwonkwonn/ovn-go-cms/ovs/operation"
	"github.com/kwonkwonn/ovn-go-cms/server"
	"github.com/kwonkwonn/ovn-go-cms/service"
)

const NB_DB string = "10.5.15.3"

func main(){


	ovnClient, err := initialize.InitializeOvnClient(NB_DB)
	if err != nil {
		log.Fatalf("Failed to initialize OVN client: %v", err)
	}
	listCon := externalmodel.NewContext()
	operator := &operation.Operator{
		Client: ovnClient,
		ListCon: listCon,
	}


	operator.InitializeLogicalDevices()
	if operator.ListCon.IsInitialized() == false {
		err:= operator.InitialSetting()
		if err != nil {
			panic("initialize error: " + err.Error())
		}
	}

	fmt.Println("ExternRouters: ", operator.ListCon.EXPList)
	fmt.Println("ExternSwitchs: ", operator.ListCon.EXSList)
	handler:=service.Handler{
		Operator: operator,
	}



	server.InitServer(8081,handler)

	
	
	select{}
}

func init() {

}
