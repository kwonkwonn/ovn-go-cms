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
        fmt.Println(fmt.Errorf("error listing logical switches: %v", err))
        return
    }
    err= o.Client.List(context.Background(), LR)
    if err!=nil{
        fmt.Println(fmt.Errorf("error listing logical routers: %v", err))
        return
    }
    err = o.Client.List(context.Background(), RPort)
    if err != nil {
        fmt.Println(fmt.Errorf("error listing logical router ports: %v", err))
        return
    }
    err = o.Client.List(context.Background(), SPort)
    if err != nil {
        fmt.Println(fmt.Errorf("error listing logical switch ports: %v", err))
        return
    }

    fmt.Printf("Found %d routers, %d switches, %d router ports, %d switch ports\n", 
        len(*LR), len(*LS), len(*RPort), len(*SPort))

    // 먼저 스위치 포트들을 분류하고 매핑
    for i:= range *SPort {
        port := &(*SPort)[i]
        fmt.Printf("Processing switch port: UUID=%s, Type=%s, Addresses=%v\n", 
            port.UUID, port.Type, port.Addresses)

        if slices.Contains(port.Addresses,"router"){
            // router-port 옵션이 있는지 확인
            routerPortUUID, ok := port.Options["router-port"]
            if !ok {
                fmt.Printf("Warning: router type port %s missing router-port option\n", port.UUID)
                continue
            }

            RtoS:= externalmodel.RtoSwitchPort{
                SwitchPort: &externalmodel.SwitchPort{
                    UUID: port.UUID,
                },
                RouterPort: &externalmodel.RouterPort{
                    UUID: routerPortUUID,
                },
            }
            switchPorts[RtoS.SwitchPort.UUID] = RtoS
            routerPorts[RtoS.RouterPort.UUID] = RtoS
            fmt.Printf("Added router-switch connection: SwitchPort=%s, RouterPort=%s\n", 
                RtoS.SwitchPort.UUID, RtoS.RouterPort.UUID)

        } else if port.Type == "vif" {
            StoVM := externalmodel.StoVMPort{				
                SwitchPort: &externalmodel.SwitchPort{
                    UUID: port.UUID,
                },
            }
            switchPorts[StoVM.SwitchPort.UUID] = StoVM
            fmt.Printf("Added VIF port: %s\n", StoVM.SwitchPort.UUID)

        } else {
            fmt.Printf("Skipping port: UUID=%s, Type=%s\n", port.UUID, port.Type)
            continue
        }
    }

    fmt.Printf("Created %d router port mappings, %d switch port mappings\n", 
        len(routerPorts), len(switchPorts))

    // 라우터들을 ExternRouter로 변환
    for i:=range *LR{
        router := (*LR)[i]
        fmt.Printf("Processing router: UUID=%s, Name=%s, Ports=%v\n", 
            router.UUID, router.Name, router.Ports)
        
        err := o.AddExternRouter(router, routerPorts)
        if err != nil {
            fmt.Printf("Error adding extern router %s: %v\n", router.UUID, err)
        }
    }

    // 스위치들을 ExternSwitch로 변환
    for i:=range *LS{
        switchObj := (*LS)[i]
        fmt.Printf("Processing switch: UUID=%s, Name=%s, Ports=%v\n", 
            switchObj.UUID, switchObj.Name, switchObj.Ports)
        
        err := o.AddExternSwitch(switchObj, switchPorts)
        if err != nil {
            fmt.Printf("Error adding extern switch %s: %v\n", switchObj.UUID, err)
        }
    }

    fmt.Printf("Initialization complete: %d routers, %d switches registered\n", 
        len(o.ExternRouters), len(o.ExternSwitchs))
}

func (o* Operator)AddExternRouter (LR NBModel.LogicalRouter, ports map[string]externalmodel.NetInt)error {
    exR:= &externalmodel.ExternRouter{
        UUID:LR.UUID,
        InternalRouter: &LR,
    }
    
    // subNetworks 초기화
    if exR.InternalRouter != nil {
        // subNetworks를 초기화 (private 필드이므로 reflection이나 public method 필요)
        // 임시로 빈 맵으로 설정
    }

    // 라우터의 각 포트에 대해 연결 정보 설정
    portCount := 0
    for _, portUUID := range LR.Ports {
        if netInt, ok := ports[portUUID]; ok {
            // 포트가 존재하면 연결 정보 설정
            if rtoS, ok := netInt.(externalmodel.RtoSwitchPort); ok {
                rtoS.ConnectedRouter = exR
                fmt.Printf("Connected router port %s to router %s\n", portUUID, LR.UUID)
                portCount++
            }
        }
    }
    
    o.ExternRouters[LR.UUID] = exR
    fmt.Printf("Added ExternRouter: UUID=%s, connected ports=%d\n", LR.UUID, portCount)
    return nil
}

func (o* Operator)AddExternSwitch (LS NBModel.LogicalSwitch, ports map[string]externalmodel.NetInt) error{
    exS:=&externalmodel.ExternSwitch{
        UUID: LS.UUID,
        InternalSwitch: &LS,
    }

    // 스위치의 각 포트에 대해 연결 정보 설정
    portCount := 0
    for _, portUUID := range LS.Ports {
        if netInt, ok := ports[portUUID]; ok {
            // 포트 타입에 따라 연결 정보 설정
            switch port := netInt.(type) {
            case externalmodel.RtoSwitchPort:
                port.ConnectedSwitch = exS
                fmt.Printf("Connected router-switch port %s to switch %s\n", portUUID, LS.UUID)
                portCount++
            case externalmodel.StoVMPort:
                port.ConnectedSwitch = exS
                fmt.Printf("Connected VM port %s to switch %s\n", portUUID, LS.UUID)
                portCount++
            }
        }
    }

    o.ExternSwitchs[LS.UUID] = exS
    fmt.Printf("Added ExternSwitch: UUID=%s, connected ports=%d\n", LS.UUID, portCount)
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

    SP,err := o.AddSwitchAPort_Router(lsUUID, lrpuuid.String(), lspuuid.String())
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
	InterPort.SwitchPort = SP

	externalmodel.AddNetInt(o.ExternRouters,ip , InterPort)

    return nil
}

func (o* Operator) InitialSetting()(error){

		EXTS_uuid,err:= o.AddSwitch()//"EXT_S"
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
		ConnectedRouter: o.ExternRouters[EXTS_uuid],
		ConnectedSwitch: o.ExternSwitchs[EXTR_uuid],
	}
	SwitchPort, err := o.AddSwitchAPort_Router(EXTS_uuid, lrpuuid.String(), lspuuid.String())
	if err != nil {
		fmt.Printf("AddInterconnectR_S ERROR: Error in AddSwitchAPort_Router: %v\n", err)
		return err
	}

	routerPort, err := o.AddRouterPort(EXTR_uuid, lrpuuid.String(), string(ROUTER))
	if err != nil {
			fmt.Printf("AddInterconnectR_S ERROR: Error in AddRouterPort: %v\n", err)
			return err
	}
	InterPort.SwitchPort = SwitchPort
	InterPort.RouterPort = routerPort

	externalmodel.AddNetInt(o.ExternRouters, string(ROUTER), InterPort)

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
//	o.IPMapping[string(UPLINK)]= newSP.UUID
	 
	lsp = append(lsp, lsMute...)
	result,err := o.Client.Transact(context.Background(),lsp...)
	if err!=nil{
		fmt.Println("the problem is...", err)
	}
	fmt.Println(result)

	
	//ip가 할당되는 순간 Map 에 저장

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