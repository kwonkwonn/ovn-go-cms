package initialize

import (
	"errors"

	client "github.com/ovn-kubernetes/libovsdb/client"
	model "github.com/ovn-kubernetes/libovsdb/model"
)

 


func InitializeNBDBModel() (*model.ClientDBModel, error) {
	dbModelReq, _ := model.NewClientDBModel("OVN_Northbound", map[string]model.Model{
		"Logical_Switch": &Logical_Switch{},
})

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

	return ovnClient, nil
}