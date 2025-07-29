package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
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

	newIP := externalmodel.FindRemainIP(h.Operator.ExternRouters, request.RequestSubnet, externalmodel.VIF)
	fmt.Println(newIP)

	InstUUID, err := util.UUIDGenerator()
	if err != nil 	{
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

	fmt.Print(h.Operator.ExternRouters[string(operation.ROUTER)].SubNetworks)
	interconnectInt := externalmodel.GetNetInt(h.Operator.ExternRouters, request.RequestSubnet+"1")
	switchs := interconnectInt[0].(externalmodel.RtoSwitchPort).ConnectedSwitch
	err = h.Operator.SwitchesPortConnect([]string{switchs.UUID}, newIP, InstUUID.String(), mac)
	if err != nil {
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


	routerUUID := h.Operator.ExternRouters[string(operation.ROUTER)].UUID
	if routerUUID == "" {
		panic("router not exist, something went wrong, critical")
	}



	VifIP := externalmodel.FindRemainIP(h.Operator.ExternRouters, request.RequestSubnet, externalmodel.VIF)
	fmt.Println(VifIP)
	SwitchPortIP :=externalmodel.FindRemainIP(h.Operator.ExternRouters, request.RequestSubnet, externalmodel.SWITCH)
	

 	swUUID,err := h.Operator.AddSwitch()
	if err!=nil{
		fmt.Println("add switch error")
		fmt.Printf("%v",fmt.Errorf("http sending error, cleanning"))
		return
	}

	mac,err:=util.MacGenerator()
	if err!=nil{
		w.Write([]byte(fmt.Errorf("mac generating error, cleaning").Error()))
		return
	}
	InstUUID,err := util.UUIDGenerator()
	if err!=nil{
		fmt.Println("no such switch exist")
	}

	err=h.Operator.AddInterconnectR_S(swUUID,routerUUID,SwitchPortIP)
	if err!=nil{
		fmt.Println(err)
	}


	_, err = h.Operator.AddSwitchAPort(swUUID, VifIP, InstUUID.String(), mac)
	if err != nil {
		fmt.Println("add switch error")
		fmt.Printf("%v", fmt.Errorf("http sending error, cleanning"))
		return
	}




	command:= "ovn-nbctl" 
    args := []string{
        "lr-nat-add",
        h.Operator.ExternRouters[routerUUID].UUID,
        "snat",
        string(operation.ROUTER),
        request.RequestSubnet+"0/24",
    }


    cmd := exec.Command(command, args...) // `exec.Command`는 명령어와 인자를 분리해서 받는 것이 더 안전합니다.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
    err = cmd.Run()
    if err != nil {
        fmt.Println("error creating router command", err)
    }



	result := &NewInstanceResult{
		MacAddress: mac,
		IP: VifIP,
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




// func (h *Handler) DelNetVM(w http.ResponseWriter, r *http.Request) {
// 	body,err:= io.ReadAll(r.Body)
// 	if err!=nil{
// 		fmt.Println("del switch error")
// 		w.Write([]byte(err.Error()))
// 		return
// 	}
// 	defer r.Body.Close()
// 	request:= &DelInstanceRequest{}
// 	err= json.Unmarshal(body, request)
// 	if err!=nil{
// 		fmt.Println("del switch error")
// 		w.Write([]byte(err.Error()))
// 		return
// 	}
// 	NetInterface,err:=util.GetNetWorkInterface(request.IP)
// 	if err!=nil{
// 		fmt.Println("del switch error")
// 		w.Write([]byte(err.Error()))
// 		return
// 	}

// 	count :=0
// 	for i:= 11; i<=254; i++{
// 		IP := NetInterface+strconv.Itoa(i)
// 		_,ok:= h.Operator.IPMapping[IP]
// 		if !ok{
// 			continue
// 		}
// 		count++
// 		if count >=2{
// 			fmt.Println("more than 2 devices exist, cannot delete")
// 			w.Write([]byte(fmt.Errorf("more than 2 devices exist, cannot delete").Error()))
// 			return
// 		}
// 	}
// 	if count == 0{
// 		fmt.Println("no such device exist byungsn")
// 		w.Write([]byte(fmt.Errorf("no such device exist").Error()))
// 		return
// 	}

// 	err= h.Operator.DelSwitch(h.Operator.IPMapping[NetInterface+"1"])
// 	if err!=nil{
// 		fmt.Println("del switch error")
// 		w.Write([]byte(err.Error()))
// 		return
// 	}

// 	err= h.Operator.DelRouterPort(NetInterface+"1")
// 	if err!=nil{
// 		fmt.Println("del router port error")
// 		w.Write([]byte(err.Error()))
// 		return
// 	}


// 	result := &DelInstanceResult{
// 		Detail: fmt.Errorf("delete vm success"),
// 	}
// 	data,err:= json.Marshal(result)
// 	if err!=nil{
// 		fmt.Printf("%v",fmt.Errorf("http sending error, cleanning"))
// 		w.Write([]byte(fmt.Errorf("http sending error, cleanning").Error()))
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(data)
// 	fmt.Println("delete vm success")

// }


// func (h *Handler) DelVM(w http.ResponseWriter, r *http.Request) {
// 	body, err:= io.ReadAll(r.Body)
// 	if err!=nil{
// 		fmt.Println("del switch error")
// 		w.Write([]byte(err.Error()))
// 		return
// 	}
// 	defer r.Body.Close()
// 	request:= &DelInstanceRequest{}
// 	err= json.Unmarshal(body, request)
// 	if err!=nil{
// 		fmt.Println("del switch error")
// 		w.Write([]byte(err.Error()))
// 		return
// 	}


// 	err = h.Operator.DelSwitchPort(request.IP)
// 	if err!=nil{
// 		fmt.Println("del switch port error")
// 		w.Write([]byte(err.Error()))
// 		return
// 	}
	

// 	result := &DelInstanceResult{
// 		Detail: fmt.Errorf("delete vm success"),
// 	}
// 	data,err:= json.Marshal(result)
// 	if err!=nil{
// 		fmt.Printf("%v",fmt.Errorf("http sending error, cleanning"))
// 		w.Write([]byte(fmt.Errorf("http sending error, cleanning").Error()))
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(data)
// 	fmt.Println("delete vm success")
// }