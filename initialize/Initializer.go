package initialize

import (
	"context"
	"fmt"

	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
	"github.com/ovn-kubernetes/libovsdb/client"
	model "github.com/ovn-kubernetes/libovsdb/model"
)

func InitializeNBDBModel() (*model.ClientDBModel, error) {
	dbModelReq, _ := NBModel.FullDatabaseModel()
	return &dbModelReq, nil
}

func InitializeOvnClient(IPAddressNB string) (client.Client, error) {
	dbModel, err := InitializeNBDBModel()
	if err != nil {
		return nil, err
	}
	//put valid address and address for northdb

	ovnClient, err := client.NewOVSDBClient(*dbModel, client.WithEndpoint("tcp:"+IPAddressNB+":6641"))
	if err != nil {
		return nil, fmt.Errorf("initial connection failed, check if ovn-northdb is on: %w", err)
	}

	ovnClient.Connect(context.Background())
	ovnClient.MonitorAll(context.Background())

	return ovnClient, nil
}
