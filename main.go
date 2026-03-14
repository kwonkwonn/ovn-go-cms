package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	initialize "github.com/kwonkwonn/ovn-go-cms/initialize"
	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	"github.com/kwonkwonn/ovn-go-cms/ovs/operation"
	"github.com/kwonkwonn/ovn-go-cms/server"
	"github.com/kwonkwonn/ovn-go-cms/service"
)

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	nbDB := getEnvOrDefault("OVN_NB_DB", "10.5.15.3")
	port, err := strconv.Atoi(getEnvOrDefault("OVN_PORT", "8081"))
	if err != nil {
		log.Fatalf("invalid OVN_PORT: %v", err)
	}

	ovnClient, err := initialize.InitializeOvnClient(nbDB)
	if err != nil {
		log.Fatalf("Failed to initialize OVN client: %v", err)
	}

	Operator := &operation.Operator{
		Client: ovnClient,
	}
	Operator.ExternRouters = make(map[string]*externalmodel.ExternRouter, 0)
	Operator.ExternSwitchs = make(map[string]*externalmodel.ExternSwitch, 0)

	Operator.InitializeLogicalDevices()
	if len(Operator.ExternRouters) == 0 && len(Operator.ExternSwitchs) == 0 {
		err := Operator.InitialSetting()
		if err != nil {
			log.Fatalf("InitialSetting failed: %v", err)
		}
	}

	fmt.Println("ExternRouters: ", Operator.ExternRouters)
	fmt.Println("ExternSwitchs: ", Operator.ExternSwitchs)
	handler := service.Handler{
		Operator: Operator,
	}

	server.InitServer(port, handler)

	select {}
}

func init() {

}
