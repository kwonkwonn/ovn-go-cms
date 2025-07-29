package externalmodel

import (
	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
	"github.com/ovn-kubernetes/libovsdb/ovsdb"
)

// router

type portType string
const (
	ROUTER portType = "router"
	SWITCH portType = "switch"
	VIF     portType = "vif"
)


type RouterPort NBModel.LogicalRouterPort
type SwitchPort NBModel.LogicalSwitchPort

type ExternRouter struct {
	UUID            string
	InternalRouter  *NBModel.LogicalRouter
	subNetworks     map[string]NetInt // uuid -> port
}
// 간선 형태의 자료구조..
// 현재는 기본적인 3-tier 형태로 구성되어 있음
// 추후에 더 복잡한 형태로 확장 가능   (1개)router -> 다수 port -> 다수 switch -> 또 다수의 vm 연결
// ovn 에서는 포트가 각 연결 지점 마다 다른 uuid 를 가짐. (관리가 매우 빡셈)
// ip map을 통해 해당 연결 지점의 포트를 관리함

type NetInt interface {
	GetConnector (portType) Connector
	// GetDeletor (portType) Deleter
}
// 모든 연결 지점은 해당 인터페이스를 구현해야 함

type Connector interface {
	Connect(RequestControl) ([]ovsdb.Operation,error)
}

// type Deleter interface {
// 	Delete(RequestControl) ([]ovsdb.Operation,error)
// }
type RtoSwitchPort struct {
	ConnectedRouter *ExternRouter
	ConnectedSwitch *ExternSwitch
	SwitchPort      *SwitchPort
	RouterPort      *RouterPort
}

type StoVMPort struct {
	ConnectedSwitch *ExternSwitch
	SwitchPort      *SwitchPort
}

type ExternSwitch struct {
	UUID             string
	InternalSwitch   *NBModel.LogicalSwitch
}
/*
내부적으로 돌아가는 ovn-nb를 추상화하는 파일
**이유:내부 모든 로직의 파악이 끝나지 않음
**데이터베이스 조회를 최소한으로 진행(특히 읽기)
**아직 모든 함수를 사용하지 않기 때문에 최소한의 인터페이스+
필요한 정보만 저장하기 위해 ex)uuid -> 연결 디바이스 매핑

initialize 폴더에서 시스템 가동시 db에서 모든
*/

type Chassis struct {
	UUID string `yaml:"uuid"`
	IP   string `yaml:"ip"`
	Tag  string `yaml:"tag"`
}

type Config struct {
	ChassisList []Chassis `yaml:"chassis"`
}


