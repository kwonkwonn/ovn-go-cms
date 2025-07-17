package operation

import (
	"context"
	"fmt"
	"sync"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"

	"github.com/ovn-kubernetes/libovsdb/client"
	"github.com/ovn-kubernetes/libovsdb/ovsdb"
)

//
type KNOWN_DEVICES string

const ( 
 UPLINK KNOWN_DEVICES = "UPLINK"
 //ip 가 할당 되어 있지 않지만 한번씩 필요한 놈들,
 // o.ipmap map[string]string 에 스트링으로 저장 
 ROUTER KNOWN_DEVICES = "10.5.15.4" // 추후에 getenv등으로 숨김 , const 라서 그렇게 초기화 될지는 몰?루
)

type Operator struct{
	Client client.Client
	ExternRouters map[string]*externalmodel.ExternRouter
	ExternSwitchs map[string]*externalmodel.ExternSwitch
	IPMapping map[string]string // device uuid
	Transaction []Transaction 
	mutex sync.Mutex
}

type Transaction struct{
	DBTransact []ovsdb.Operation
	SideEffect func(...any)(error)   // 
	Undo func(...any)(error)
}


func (o*Operator)AddTransaction(operation ...ovsdb.Operation, )(error){
	newTrans := Transaction{ 
		DBTransact:operation,
	}
}


// 명령들을 operations 에 있는 명령들을 수행합니다 
func (o*Operator) Transact()(error){
	result, err:= o.Client.Transact(context.Background(),o.operations...)
	if err!=nil{
		return fmt.Errorf("error on transactioning: %v")
	}
	fmt.Println(result)
}

// operation 들을 제거 합니다
func (o*Operator) Flush(){
	o.Transaction = make([]Transaction, 0)
}


