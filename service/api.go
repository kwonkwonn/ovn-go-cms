package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	"github.com/kwonkwonn/ovn-go-cms/ovs/operation"
	"github.com/kwonkwonn/ovn-go-cms/ovs/util"
)





func (h *Handler) CreateNewVm(w  http.ResponseWriter,r *http.Request ){
	body, err:= io.ReadAll(r.Body)
	if err!=nil{
		fmt.Println("add switch error")

		w.Write([]byte(err.Error()))
	}
	defer r.Body.Close()
	request:= &NewInstanceRequeset{}
	err= json.Unmarshal(body, request)
	if err!=nil{
		fmt.Println("add switch error")
		w.Write([]byte(err.Error()))
		return
	}

	var RtoSInterfaceIP = request.RequestSubnet + "1"
	RtoSInterface:= externalmodel.GetNetInt(h.Operator.ExternRouters, RtoSInterfaceIP)
	
	
	var swUUID string 
	newvifIP := externalmodel.FindRemainIP(h.Operator.ExternRouters, request.RequestSubnet, externalmodel.VIF)
	
	mac,err:=util.MacGenerator()
		if err!=nil{
			w.Write([]byte(fmt.Errorf("mac generating error, cleaning").Error()))
			return
		}
	InstUUID,err := util.UUIDGenerator()
		if err!=nil{
			fmt.Println("no such switch exist")
	}


	if len(RtoSInterface) == 0 {
		fmt.Println("request for new subnet, creating new router port")
		routerUUID := h.Operator.ExternRouters[string(operation.ROUTER)].UUID
			if routerUUID == "" {
			panic("router not exist, something went wrong, critical")
		}

		swUUID,err = h.Operator.AddSwitch()
			if err!=nil{
				fmt.Println("add switch error")
				fmt.Printf("%v",fmt.Errorf("http sending error, cleanning"))
				return
		}
		err=h.Operator.AddInterconnectR_S(swUUID,routerUUID, RtoSInterfaceIP)
			if err!=nil{
				fmt.Println(err)
		}

	}else{
		swUUID = RtoSInterface[0].(*externalmodel.RtoSwitchPort).ConnectedSwitch.UUID
	}

	

	
	err = h.Operator.SwitchesPortConnect([]string{swUUID}, newvifIP, InstUUID.String(), mac)

	if err != nil {
		fmt.Println(err)
	}




	result:= &NewInstanceResult{
		MacAddress: mac,
		IP: newvifIP,
		IfaceID: InstUUID.String(),
	}

	data,err:= json.Marshal(result)
	if err!=nil{
		fmt.Printf("%v",fmt.Errorf("http sending error, cleanning"))
		w.Write([]byte(fmt.Errorf("http sending error, cleanning").Error()))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)


}



// func (h *Handler) CreateNewNetVm(w http.ResponseWriter,r *http.Request ){
// 	body, err:= io.ReadAll(r.Body)
// 	if err!=nil{
// 		fmt.Println("add switch error")

// 		w.Write([]byte(err.Error()))
// 		return 
// 	}

// 	defer r.Body.Close()
// 	request:= &NewInstanceRequeset{}
// 	err= json.Unmarshal(body, request)
// 	if err!=nil{
// 		fmt.Println("add switch error")
// 		w.Write([]byte(err.Error()))
// 		return

// 	}


// 	routerUUID := h.Operator.ExternRouters[string(operation.ROUTER)].UUID
// 	if routerUUID == "" {
// 		panic("router not exist, something went wrong, critical")
// 	}



// 	VifIP := externalmodel.FindRemainIP(h.Operator.ExternRouters, request.RequestSubnet, externalmodel.VIF)
// 	fmt.Println("ip found" + VifIP)
// 	SwitchPortIP :=externalmodel.FindRemainIP(h.Operator.ExternRouters, request.RequestSubnet, externalmodel.SWITCH)
// 	fmt.Println("ip found" + SwitchPortIP)

//  	swUUID,err := h.Operator.AddSwitch()
// 	if err!=nil{
// 		fmt.Println("add switch error")
// 		fmt.Printf("%v",fmt.Errorf("http sending error, cleanning"))
// 		return
// 	}

// 	mac,err:=util.MacGenerator()
// 	if err!=nil{
// 		w.Write([]byte(fmt.Errorf("mac generating error, cleaning").Error()))
// 		return
// 	}
// 	InstUUID,err := util.UUIDGenerator()
// 	if err!=nil{
// 		fmt.Println("no such switch exist")
// 	}

// 	err=h.Operator.AddInterconnectR_S(swUUID,routerUUID,SwitchPortIP)
// 	if err!=nil{
// 		fmt.Println(err)
// 	}

//     fmt.Println("zxvavasdvdsas", h.Operator.ExternRouters[string(operation.ROUTER)].SubNetworks)

// 	_, err = h.Operator.AddSwitchAPort(swUUID, VifIP, InstUUID.String(), mac)
// 	if err != nil {
// 		fmt.Println("add switch error")
// 		fmt.Printf("%v", fmt.Errorf("http sending error, cleanning"))
// 		return
// 	}



// 	result := &NewInstanceResult{
// 		MacAddress: mac,
// 		IP: VifIP,
// 		IfaceID: InstUUID.String(),
// 	}
// 	data,err:= json.Marshal(result)
// 	if err!=nil{
// 		fmt.Printf("%v",fmt.Errorf("http sending error, cleanning"))
// 		w.Write([]byte(fmt.Errorf("http sending error, cleanning").Error()))
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(data)
	
// }	
//테스트 완전 될때까지만 살려둠...

func (h *Handler) DeleteAll (w http.ResponseWriter, r* http.Request){

	h.Operator.DeleteAll()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("work done"))
}




func (h *Handler) DelNet(w http.ResponseWriter, r *http.Request) {
	body,err:= io.ReadAll(r.Body)
	if err!=nil{
		fmt.Println("del switch error")
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()
	request:= &DelInstanceRequest{}
	err= json.Unmarshal(body, request)
	if err!=nil{
		fmt.Println("del switch error")
		w.Write([]byte(err.Error()))
		return
	}


	NetSignifier,err:=util.GetNetWorkSignifier(request.IP)
	if err!=nil{
		fmt.Println("del switch error")
		w.Write([]byte(err.Error()))
		return
	}

	NetInt := externalmodel.GetNetInt(h.Operator.ExternRouters,request.IP)
	if len(NetInt) == 0 {
		fmt.Println("no such switch exist")
		w.Write([]byte(fmt.Errorf("no such switch exist").Error()))
		return
	}

	_,ok := NetInt[0].(*externalmodel.RtoSwitchPort)
	if ok {
		w.Write([]byte(fmt.Errorf("cannot delete switch port, connected to router").Error()))
		return
	}
	Port,ok := NetInt[0].(*externalmodel.StoVMPort)
	if !ok {
		fmt.Println("no such switch port exist")
		w.Write([]byte(fmt.Errorf("no such switch port exist").Error()))
		return
	}

	SwitchPort := Port.ConnectedSwitch
	if SwitchPort == nil {
		fmt.Println("switch port not connected")
		w.Write([]byte(fmt.Errorf("switch port not connected").Error()))
		return
	}


	err= h.Operator.DelSwitchPort(request.IP)
	if err!=nil{
		fmt.Println("del switch port error")
		w.Write([]byte(err.Error()))
		return
	}


	result := &DelInstanceResult{Detail: fmt.Errorf("delete switch port operation")}

	delete(h.Operator.ExternRouters[string(operation.ROUTER)].SubNetworks, request.IP) 
	//여러 라우터가 생성될 필요가 있을때 수정,,
	// 1라우터 - 1서브넷(스위치) 구조이므로 굳이 복잡한 순회를 하지 않아도 됨
	// 서브넷이 삭제되면 연결된 스위치도 삭제
	// 서브넷이 삭제되면 연결된 라우터 포트도 삭제
	// 서브넷이 삭제되면 연결된 NAT도 삭제
	nets:= externalmodel.GetAllVIF(h.Operator.ExternRouters, NetSignifier)
	if len(nets) == 0 {
		fmt.Println("Deleting Connected Switch")
		err= h.Operator.DelSwitch(SwitchPort.UUID)
		if err!=nil{
			fmt.Println("del switch error")
			w.Write([]byte(err.Error()))
			return
		}
		
		err= h.Operator.DelRouterPort(NetSignifier+"1")
		if err!=nil{
			fmt.Println("del router port error %w", err)
			w.Write([]byte(fmt.Errorf("del router port error %w", err).Error()))
			return
		}
		result.Detail = fmt.Errorf("delete switch and router port success")
		delete(h.Operator.ExternRouters[string(operation.ROUTER)].SubNetworks, NetSignifier+"1")
	}

		// nat 삭제 도입 할 예정
	result.Detail = fmt.Errorf("%v switch port deleted", result.Detail)
	data,err:= json.Marshal(result)
	if err!=nil{
		fmt.Printf("%v",fmt.Errorf("http sending error, cleanning"))
		w.Write([]byte(fmt.Errorf("http sending error, cleanning").Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	fmt.Println("delete vm success")

}


