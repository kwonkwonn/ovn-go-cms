package main

import (
	"context"
	"log"

	initialize "github.com/kwonkwonn/ovn-go-cms/initialize"
)


func main(){

	ovnClient, err := initialize.InitializeOvnClient("127.0.0.1")
	if err != nil {
		log.Fatalf("Failed to initialize OVN client: %v", err)
	}

	ovnClient.Connect(context.Background())


	select{}
}
