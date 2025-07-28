package operation

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	NBModel "github.com/kwonkwonn/ovn-go-cms/ovs/internalModel"
	"github.com/kwonkwonn/ovn-go-cms/ovs/util"
	"github.com/ovn-kubernetes/libovsdb/model"
	"github.com/ovn-kubernetes/libovsdb/ovsdb"
	"gopkg.in/yaml.v3"
)


func (o* Operator) ChassisInitializing(RouterUplinkPort string )(error){
	//Chassis 에 대한 정보는 ovn-sb에 저장되어 있기 때문에
	// 현재는 외부 파일에서 읽어오고 있습니다.
	filepath:="./.chassis.yaml"
	cfg:=  &externalmodel.Config{
		ChassisList: make([]externalmodel.Chassis, 0),
	}
	data, err:=os.ReadFile(filepath)
	if err!=nil{
		panic("reading chassis file error, terminating process")
	}
	yaml.Unmarshal(data,cfg)

	//router port 가 인식이 안되는 오류때문에, 커맨드로 대체하는중,,, 해결되면 고칠예정 

	for i:= range cfg.ChassisList{
		command:= "ovn-nbctl"
		args:= []string{
			"lrp-set-gateway-chassis",
			RouterUplinkPort,
			cfg.ChassisList[i].UUID,
			"100",
		}
		cmd := exec.Command(command, args...) // 수정된 command 사용

		// 명령어 실행 시 표준 출력과 표준 에러를 볼 수 있도록 연결
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	
		fmt.Printf("Executing command: %s %v\n", command, args)
	
		err:= cmd.Run()
		if err!=nil{
			// 더 구체적인 에러 메시지를 반환합니다.
			return fmt.Errorf("error setting router chassis priority for UUID %s: %w", cfg.ChassisList[i].UUID, err)
		}
	}


return nil
}

func (o* Operator)InitializeLogicalDevices (){
	o.ExternRouters = make(map[string]*externalmodel.ExternRouter)
	o.ExternSwitchs = make(map[string]*externalmodel.ExternSwitch)
	routerPorts :=make(map[string]externalmodel.NetInt)
	switchPorts :=make(map[string]externalmodel.NetInt)
	
	
	LR :=&[]NBModel.LogicalRouter{}
	LS :=&[]NBModel.LogicalSwitch{}
	RPort := &[]NBModel.LogicalRouterPort{}
	SPort:= &[]NBModel.LogicalSwitchPort{}

	err:= o.Client.List(context.Background(), LS )
	if err!=nil{
		fmt.Println(fmt.Errorf("%v", err))
	}
	err= o.Client.List(context.Background(), LR)
	if err!=nil{
		fmt.Println(fmt.Errorf("%v", err))
	}
	err = o.Client.List(context.Background(), RPort)
	if err != nil {
		fmt.Println(fmt.Errorf("error listing logical router ports: %v", err))
	}
	err = o.Client.List(context.Background(), SPort)
	if err != nil {
		fmt.Println(fmt.Errorf("error listing logical switch ports: %v", err))
	}

	for i:= range *SPort {
		if slices.Contains((*SPort)[i].Addresses,"router"){
			RtoS:= externalmodel.RtoSwitchPort{
				SwitchPort: &externalmodel.SwitchPort{
					UUID: (*SPort)[i].UUID,
				},
				RouterPort: &externalmodel.RouterPort{
					UUID: (*SPort)[i].Options["router-port"],
				},
			}
			switchPorts[RtoS.SwitchPort.UUID] = RtoS
			routerPorts[RtoS.RouterPort.UUID] = RtoS
		}else if ((*SPort)[i].Type=="vif"){
			StoVM := externalmodel.StoVMPort{
				SwitchPort: externalmodel.SwitchPort{
					UUID: (*SPort)[i].UUID,
				},
			}
			switchPorts[StoVM.SwitchPort.UUID] = StoVM
		}else{
			continue
		}
	}


	for i:=range *LR{
		o.AddExternRouter((*LR)[i])
	}
	for i:=range *LS{
		o.AddExternSwitch((*LS)[i])
	}
}

func (o* Operator)AddExternRouter (LR NBModel.LogicalRouter)error {
	exR:= &externalmodel.ExternRouter{
		UUID:LR.UUID,
		InternalRouter: &LR,
	}

	// R_Port:= &NBModel.LogicalRouterPort{

	// }

	o.ExternRouters[LR.UUID] = exR
	return nil
}

func (o* Operator)AddExternSwitch (LS NBModel.LogicalSwitch) error{
	exS:=&externalmodel.ExternSwitch{
		UUID: LS.UUID,
		//IP: yaml에서 읽어서 할당
		InternalSwitch: &LS,
	}


	o.ExternSwitchs[LS.UUID]=exS

	return nil
}




func (o * Operator) AddInterconnectR_S(lsUUID string, lrUUID string, ip string)(error){
    lrpuuid,err:=util.UUIDGenerator()
    if err!=nil{
        panic("lrpuuid generating error" )
    }
    lspuuid,err:=util.UUIDGenerator()
    if err!=nil{
        panic("lrpuuid generating error" )
    }

	InterPort:= externalmodel.RtoSwitchPort{
		ConnectedRouter: o.ExternRouters[lrUUID],
		ConnectedSwitch: o.ExternSwitchs[lsUUID],
	}

    err = o.AddSwitchAPort_Router(lsUUID, lrpuuid.String(), lspuuid.String())
    if err != nil {
        fmt.Printf("AddInterconnectR_S ERROR: Error in AddSwitchAPort_Router: %v\n", err)
        return err
    }


    routerPort, err := o.AddRouterPort(lrUUID, lrpuuid.String(),ip)
    if err != nil {
        fmt.Printf("AddInterconnectR_S ERROR: Error in AddRouterPort: %v\n", err)
        return err
    }
	InterPort.RouterPort = routerPort
    return nil
}

func (o* Operator) InitialSetting()(error){

		EXTS_uuid,err:= o.AddSwitch("EXT_S")
		if (err!=nil){
			panic("bootstraping failed, creating external Switch")
		}
		fmt.Printf("InitialSettig: Created EXTS_uuid: %s\n", EXTS_uuid)
	
		EXTR_uuid,err:=o.AddRouter(string(ROUTER))
		if (err!=nil){
			panic("bootstraping failed, creating external Switch")
		}

	
{
	lrpuuid,err:=util.UUIDGenerator()
	if err!=nil{
		panic("lrpuuid generating error" )
	}
	lspuuid,err:=util.UUIDGenerator()
	if err!=nil{
			panic("lrpuuid generating error" )
	}
	
	InterPort:= externalmodel.RtoSwitchPort{
		ConnectedRouter: o.ExternRouters[lrUUID],
		ConnectedSwitch: o.ExternSwitchs[lsUUID],
	}
	err = o.AddSwitchAPort_Router(EXTS_uuid, lrpuuid.String(), lspuuid.String())
	if err != nil {
		fmt.Printf("AddInterconnectR_S ERROR: Error in AddSwitchAPort_Router: %v\n", err)
		return err
	}

	routerPort, err := o.AddRouterPort(EXTR_uuid, lrpuuid.String(), string(ROUTER))
	if err != nil {
			fmt.Printf("AddInterconnectR_S ERROR: Error in AddRouterPort: %v\n", err)
			return err
	}
	InterPort.RouterPort = routerPort


	err= o.ChassisInitializing(lrpuuid.String())
	if err!= nil{
		fmt.Printf("error adding chassis priority %v", err)
	}
}
	
	br_EXTS_UUID,err := util.UUIDGenerator()
	if err!=nil{
		return fmt.Errorf("generating uuid error: br_exts_uuid")
	}
	
	//Addswitch port 함수의 복사본, 함수화가 아직 부족해서 복붙함
	// 리팩토링 대상 1번 
	value ,ok := o.ExternSwitchs[EXTS_uuid]; 
	if !ok{
		return fmt.Errorf("no such switch exist")
	}

	newSP:= &NBModel.LogicalSwitchPort{
		UUID: br_EXTS_UUID.String(),
		Name: br_EXTS_UUID.String(),
		Type: "localnet",
		Options: map[string]string{"network_name":"UPLINK"},
		}
	Address :="unknown"
	newSP.Addresses=append(newSP.Addresses, Address)
	
	lsp , err := o.Client.Create(newSP)
	if err!=nil{
		return fmt.Errorf("%v", err)
	}
	
	value.InternalSwitch.Ports = append(value.InternalSwitch.Ports, newSP.UUID)
	lsMute,_  := o.Client.Where(value.InternalSwitch).Mutate( value.InternalSwitch, model.Mutation{
		Field: &value.InternalSwitch.Ports,
		Mutator: ovsdb.MutateOperationInsert,
		Value: value.InternalSwitch.Ports,	
		
	}) 
	o.IPMapping[string(UPLINK)]= newSP.UUID
	 
	lsp = append(lsp, lsMute...)
	result,err := o.Client.Transact(context.Background(),lsp...)
	if err!=nil{
		fmt.Println("the problem is...", err)
	}
	fmt.Println(result)

	util.SaveMapYaml(o.IPMapping)
	
	//ip가 할당되는 순간 Map 에 저장

	for _,i:= range o.IPMapping{
		fmt.Println(i)
	}
	command := "ovn-nbctl" 
    args := []string{
		"lr-route-add",
        EXTR_uuid,
		"0.0.0.0/0",
		string(DEFAULT_GATEWAY),
    }

    cmd := exec.Command(command, args...) 
    err = cmd.Run()
    if err != nil {
        return fmt.Errorf("error creating router command, %v", err)
    }


	return nil
}