package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	
	// existUUID:=h.Operator.CheckIPExistance(request.RequestSubnet)
	// if existUUID!=""{
	// 	fmt.Println(request)

	// 	w.Write([]byte(fmt.Errorf("network already exist").Error()))
	// 	return
	// }
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
	
	}





}


func (h *Handler) CreateNewNetVm(w http.ResponseWriter,r *http.Request ){
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
	
	// existUUID:=h.Operator.CheckIPExistance(request.RequestSubnet)
	// if existUUID!=""{
	// 	fmt.Println(request)

	// 	w.Write([]byte(fmt.Errorf("network already exist").Error()))
	// 	return
	// }
	newIP_VM := h.Operator.AvailableIP_VM(request.RequestSubnet)
	fmt.Println(newIP_VM)
	mac,err:=util.MacGenerator()
	if err!=nil{
		fmt.Printf("%v",fmt.Errorf("http sending error, cleanning"))
		w.Write([]byte(fmt.Errorf("mac generating error, cleaning").Error()))
		//클린업 함수 만들어야 함 
		return
	}
	newIP_Dev := h.Operator.AvailableIP_Dev(request.RequestSubnet)
	fmt.Println(newIP_Dev)
 
	fmt.Println(mac)
	swUUID1,err := h.Operator.AddSwitch(newIP_Dev)
	if err!=nil{
		fmt.Println("add switch error")
		fmt.Printf("%v",fmt.Errorf("http sending error, cleanning"))
		return
	}
 
	InstUUID,err := util.UUIDGenerator()
	if err!=nil{
		fmt.Println("no such switch exist")
	}

	err = h.Operator.AddSwitchAPort(swUUID1,newIP_Dev, InstUUID.String(),mac)
	if err!=nil{
		fmt.Println("add switch error")
		fmt.Printf("%v",fmt.Errorf("http sending error, cleanning"))

		return
	}
 

	result:= &NewInstanceResult{
		MacAddress: mac,
		IP: newIP_Dev,
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