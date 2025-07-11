package externalmodel

import (
	"github.com/google/uuid"
	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
)

/*
내부적으로 돌아가는 ovn-nb를 추상화하는 파일
**이유:내부 모든 로직의 파악이 끝나지 않음
**데이터베이스 조회를 최소한으로 진행(특히 읽기)
**아직 모든 함수를 사용하지 않기 때문에 최소한의 인터페이스+
필요한 정보만 저장하기 위해 ex)uuid -> 연결 디바이스 매핑

initialize 폴더에서 시스템 가동시 db에서 모든

*/
type subnet int8


type ExternRouter struct{
	UUID uuid.UUID
	IP string
	Subnet subnet //코어 컨트롤 ip를 net.Ip 로 바꿀때 함께 바꿀예정
	InternalRouter NBModel.LogicalRouter
}

type ExternSwitch struct{
	UUID uuid.UUID
	ParentRouter *ExternRouter
	IP    subnet  //코어 컨트롤 ip를 net.Ip 로 바꿀때 함께 바꿀예정
	InternalRouter NBModel.LogicalSwitch
}