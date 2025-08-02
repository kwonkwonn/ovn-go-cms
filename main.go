package main

import (
	"context"
	"fmt"
	"log"

	initialize "github.com/kwonkwonn/ovn-go-cms/initialize"
	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
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

	Operator := &operation.Operator{
		Client: ovnClient,
	}
	Operator.ExternRouters = make(map[string]*externalmodel.ExternRouter, 0)
	Operator.ExternSwitchs = make(map[string]*externalmodel.ExternSwitch,0)

	Operator.InitializeLogicalDevices()
	if len(Operator.ExternRouters) == 0 && len(Operator.ExternSwitchs) == 0 {
		err:= Operator.InitialSetting()
		if err != nil {
			panic("initialize error: " + err.Error())
		}
	}

	fmt.Println("ExternRouters: ", Operator.ExternRouters)
	fmt.Println("ExternSwitchs: ", Operator.ExternSwitchs)
	handler:=service.Handler{
		Operator: Operator,

	}
	//     nat:= &NBModel.NAT{
    //     Type: "snat",
    //     LogicalIP: "20.20.22.1" + "/24",
    //     ExternalIP: "10.5.15.4",
    // }
		// ports:= &NBModel.LogicalRouterPort{
		// 	// Networks: []string{"20.20.22.1"+ "/24"},
		// }
    portss:= &[]NBModel.NAT{}
     Operator.Client.List(context.Background(),portss)

    fmt.Println("asfsfidsaofbadsfoisabfsodifsbfa", portss)
    fmt.Println("zxvavasdvdsas", Operator.ExternRouters[string(operation.ROUTER)].SubNetworks)	


	server.InitServer(8081,handler)

	
	
	select{}
}

func init() {

}
