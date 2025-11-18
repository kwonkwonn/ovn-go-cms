package externalmodel

import (
	"context"
	"sync"
)

type EXRList map[string]*ExternRouter
type EXSList map[string]*ExternSwitch

type Context struct{
	context.Context
	EXPList  EXRList
	EXSList  EXSList
	mutex  	*sync.RWMutex
}

// type Context 는 컨텍스트를 임베딩 함.
// 외부로 관리되는 인메모리 모델 리스트(EXRList, EXSList)의 동시성 문제 해결을 위해 사용됨.

func NewContext() *Context {
	return 	&Context{
		EXPList: make(EXRList, 0),
		EXSList: make(EXSList,0),
		mutex:  &sync.RWMutex{},
	}
}



func (cont *Context) IsInitialized() bool{
	if len(cont.EXPList) == 0 && len(cont.EXSList) == 0 {
		return false
	}
	return true
}
// After call all devices from DB,
// if not set, should make new maps


func (cont Context) GetRouter(uuid string) *ExternRouter {
	cont.mutex.RLock()
	defer cont.mutex.RUnlock()
	if router, ok := cont.EXPList[uuid]; ok {
		return router
	}
	return nil
}


func (cont Context) GetSwitch(uuid string) *ExternSwitch {
	cont.mutex.RLock()	
	defer cont.mutex.RUnlock()
	if switchDevice, ok := cont.EXSList[uuid]; ok {
		return switchDevice
	}
	return nil
}




