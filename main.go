package main

import (
	"context"
	"fmt"
	"log"
	"time"

	initialize "github.com/kwonkwonn/ovn-go-cms/initialize"
)


func main(){

	ovnClient, err := initialize.InitializeOvnClient("127.0.0.1")
	if err != nil {
		log.Fatalf("Failed to initialize OVN client: %v", err)
	}
	ovnClient.Connect(context.Background())
	ovnClient.MonitorAll(context.Background())
	ops, err := ovnClient.Create(&initialize.Logical_Switch{
		Name: "foo",
	})
	if err!=nil{
		fmt.Println(err)
	}
	
	_,err = ovnClient.Transact(context.Background(),ops...)
	if(err!=nil){
		fmt.Println(err)
	}

	// OVN이 내부적으로 동기화할 시간을 줍니다.
// 이는 임시 방편이며, 실제로는 waitForCacheConsistent가 이 역할을 해야 합니다.
time.Sleep(2 * time.Second) // 2초 대기

	ls := &[]initialize.Logical_Switch{}
	// test := &initialize.Logical_Switch{Name:"foo"}
	ovnClient.List(context.Background(), ls )

	time.Sleep(2 * time.Second) // 2초 대기
	fmt.Println("List of logical devices:",ls)


	for i, _ := range *ls{
		ovs_to_del:=&initialize.Logical_Switch{UUID: (*ls)[i].UUID}
		ops,_:= ovnClient.Where(ovs_to_del).Delete()
		ovnClient.Transact(context.Background(),ops...)
	}

	ls = &[]initialize.Logical_Switch{}
	// test := &initialize.Logical_Switch{Name:"foo"}
	ovnClient.List(context.Background(), ls )

	time.Sleep(2 * time.Second) // 2초 대기
	fmt.Println("List of logical devices:",ls)
	select{}
}
