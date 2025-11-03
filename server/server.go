package server

import (
	"log"
	"net" // net 패키지 추가
	"net/http"
	"strconv"

	"github.com/kwonkwonn/ovn-go-cms/service"
)

func InitServer(portNum int, handler service.Handler) {
	addr := "0.0.0.0:" + strconv.Itoa(portNum)

	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}
	defer listener.Close()

	// http.HandleFunc("POST /New/Instance", handler.CreateNewNetVm)
	// 새로운 가상머신을 생성, 새로운 네트워크를 만드는 과정까지 추상화 됨
	// http.HandleFunc("POST /New/Net",)
	// 새로운 네트워크를 생성, 새로운 서브넷
	http.HandleFunc("POST /New/Instance", handler.CreateNewVm)
	//새로운 가상 머신을 생성, 기존 네트워크에 붙
	//http.HandleFunc("POST /Add/Net")
	// 네트워크를 생성, 기존 서브넷에 붙임, 아직 안씀
	http.HandleFunc("DELETE /New/Instance", handler.DelNet)
	// 분해는 조립의 역순...
	// 가상 머신을 삭제, 연결되어 있는 포트(서브넷)도 삭제
	// 해당 연결 스위치를 삭제, 그 스위치와 연결 되어 있던 라우터 포트를 삭제, nat 삭제

	http.HandleFunc("DELETE /ALL", handler.DeleteAll)
	// //**테스트용으로만 사용, 모든 가상 디바이스 삭제**

	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
