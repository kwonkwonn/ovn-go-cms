package main

import (
	"context"
	"fmt"
	"log"
	"time"

	initialize "github.com/kwonkwonn/ovn-go-cms/initialize"
	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
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
	err=util.ReadMapNetYaml(Operator.IPMapping)
	if err!= nil{
		Operator.IPMapping= make(map[string]string,0)
	}
	Operator.InitializeLogicalDevices()

	time.Sleep(2 * time.Second) // 2초 대기

	{

	handler:=service.Handler{
		Operator: Operator,
	}
	server.InitServer(8081,handler)

	ls:= &[]NBModel.LogicalSwitch{}
	Operator.Client.List(context.Background(),ls)
	
	fmt.Println("List of logical devices:",ls)
	}
	select{}
}
