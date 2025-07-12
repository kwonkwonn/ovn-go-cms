package main

import (
	"fmt"
	"log"
	"time"

	initialize "github.com/kwonkwonn/ovn-go-cms/initialize"
	"github.com/kwonkwonn/ovn-go-cms/ovs/operation"
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
 
	// ops, err := ovnClient.Create(&NBModel.LogicalRouter{
	// 	Name: "foo",
	// 	UUID: "3a3456b8-ab16-4602-a352-0bd9db372c97",
	// })
	// if err!=nil{
	// 	fmt.Println(err)
	// }
	
	// _,err = ovnClient.Transact(context.Background(),ops...)
	// if(err!=nil){
	// 	fmt.Println(err)
	// }

	// Ors,_:=ovnClient.Create(&NBModel.LogicalRouter{})

// YO ,_:=Operator.Client.Where(Ors).Delete()
// fmt.Println(YO)
// 	ovnClient.Transact(context.Background(), YO[0])
	uuid ,err :=Operator.AddSwitch() 	
	if err!=nil{
		fmt.Println(uuid)
	}
	Operator.AddSwitchAPort(uuid, "20.10.15.24")

time.Sleep(2 * time.Second) // 2초 대기


	Operator.InitializeLogicalDevices()

	time.Sleep(2 * time.Second) // 2초 대기
	

	
	fmt.Println("List of logical devices:",Operator.ExternRouters)
	select{}
}
