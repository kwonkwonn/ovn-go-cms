package initialize

import (
	"context"
	"errors"
	"time"

	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
	client "github.com/ovn-org/libovsdb/client"
	model "github.com/ovn-org/libovsdb/model"
)

 


func InitializeNBDBModel() (*model.ClientDBModel, error) {
	dbModelReq, _ := NBModel.FullDatabaseModel()
	return &dbModelReq,nil
}


func InitializeOvnClient( IPAddressNB string ) (client.Client, error) {
	dbModel, err := InitializeNBDBModel()
	if err != nil {
		return nil, err
	}
	//put valid address and address for northdb 
	
	ovnClient, err := client.NewOVSDBClient(*dbModel, client.WithEndpoint("tcp:"+IPAddressNB+":6641"))
	if err != nil {
		panic(errors.New("initial connection failed booting ovn-cms, check if ovn-northdb is on"))
	}

	for {
		
		if ovnClient.Connected(){
			break
		}
		time.Sleep(10000)
	}
	ovnClient.Connect(context.Background())
	ovnClient.MonitorAll(context.Background())

	return ovnClient, nil
}

