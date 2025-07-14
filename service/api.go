package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kwonkwonn/ovn-go-cms/ovs/operation"
	"github.com/kwonkwonn/ovn-go-cms/ovs/util"
)




func (h *Handler) CreateNewVm(w  http.ResponseWriter,r *http.Request ){
	fmt.Println("new REQUEST")
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

	newIP := h.Operator.AvailableIP_VM(request.RequestSubnet)
	fmt.Println(newIP)
 
	InstUUID,err := util.UUIDGenerator()
	if err!=nil{
		fmt.Println("no such switch exist")
	}


	mac,err:=util.MacGenerator()
	if err!=nil{
		fmt.Printf("%v",fmt.Errorf("http sending error, cleanning"))
		w.Write([]byte(fmt.Errorf("mac generating error, cleaning").Error()))
		//클린업 함수 만들어야 함 
		return
	}
	fmt.Println(mac)

	devsUUID:= h.Operator.FindExistdev(request.RequestSubnet)
	err= h.Operator.SwitchesPortConnect(devsUUID,newIP,InstUUID.String(),mac)
	if err!=nil{
		fmt.Println(err)
	}

	result:= &NewInstanceResult{
		MacAddress: mac,
		IP: newIP,
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


func (h *Handler) CreateNewNetVm(w http.ResponseWriter,r *http.Request ){
	fmt.Println("new REQUEST")
	body, err:= io.ReadAll(r.Body)
	if err!=nil{
		fmt.Println("add switch error")

		w.Write([]byte(err.Error()))
		return 
	}

	defer r.Body.Close()
	request:= &NewInstanceRequeset{}
	err= json.Unmarshal(body, request)
	if err!=nil{
		fmt.Println("add switch error")
		w.Write([]byte(err.Error()))
		return

	}
	routerUUID,ok := h.Operator.IPMapping[string(operation.ROUTER)]
	if !ok{
		panic("router not exist, something went wrong, critical")
	}



	newIP_VM := h.Operator.AvailableIP_VM(request.RequestSubnet)
	fmt.Println(newIP_VM)

	mac,err:=util.MacGenerator()
	if err!=nil{
		fmt.Printf("%v",fmt.Errorf("http sending error, cleanning"))
		w.Write([]byte(fmt.Errorf("mac generating error, cleaning").Error()))
		//클린업 함수 만들어야 함 
		return
	}
	InstUUID,err := util.UUIDGenerator()
	if err!=nil{
		fmt.Println("no such switch exist")
	}

	newIP_Dev := h.Operator.AvailableIP_Dev(request.RequestSubnet)
	fmt.Println(newIP_Dev)
 	swUUID,err := h.Operator.AddSwitch(newIP_Dev)
	if err!=nil{
		fmt.Println("add switch error")
		fmt.Printf("%v",fmt.Errorf("http sending error, cleanning"))
		return
	}
	err = h.Operator.AddSwitchAPort(swUUID,newIP_VM, InstUUID.String(),mac)
	if err!=nil{
		fmt.Println("add switch error")
		fmt.Printf("%v",fmt.Errorf("http sending error, cleanning"))
		return
	}

	err=h.Operator.AddInterconnectR_S(swUUID,routerUUID,newIP_Dev)
	if err!=nil{
		fmt.Println(err)
	}
 	result:= &NewInstanceResult{
		MacAddress: mac,
		IP: newIP_VM,
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

func (h *Handler) DeleteAll (w http.ResponseWriter, r* http.Request){

	h.Operator.DeleteAll()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("work done"))
}

// existUUID:=h.Operator.CheckIPExistance(request.RequestSubnet)
// if existUUID!=""{
// 	fmt.Println(request)

// 	w.Write([]byte(fmt.Errorf("network already exist").Error()))
// 	return
// }