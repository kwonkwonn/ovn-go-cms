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
    portsPool :=make(map[string]externalmodel.NetInt)
    routerPort :=make(map[string]*NBModel.LogicalRouterPort)
    
    LR :=&[]NBModel.LogicalRouter{}
    LS :=&[]NBModel.LogicalSwitch{}
    RPort := &[]NBModel.LogicalRouterPort{}
    SPort:= &[]NBModel.LogicalSwitchPort{}

    // List 코드들...
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

    // 포트 매핑 생성
    for i:= range *RPort { 
        port := &(*RPort)[i]
        routerPort[port.UUID] = port
    }

    for i:= range *SPort {
        port := &(*SPort)[i]
        switchPort:= externalmodel.SwitchPort(*port)
        
        if slices.Contains(port.Addresses,"router"){
            routerPortUUID, ok := port.Options["router-port"]  
            if !ok {
                fmt.Printf("Warning: router type port %s missing router-port option\n", port.UUID)
                continue
            }
            
            if _, exists := routerPort[routerPortUUID]; !exists {
                fmt.Printf("Warning: router port %s not found\n", routerPortUUID)
                continue
            }
            
            routerPortObj := externalmodel.RouterPort(*routerPort[routerPortUUID])

            RtoS:= &externalmodel.RtoSwitchPort{
                SwitchPort: &switchPort,
                RouterPort: &routerPortObj,
            }
            
            // 포인터로 저장 (동일한 객체 참조)
            portsPool[port.UUID] = RtoS
            portsPool[routerPortUUID] = RtoS

        } else if port.Type == "vif" {
            StoVM := &externalmodel.StoVMPort{				
                SwitchPort: &switchPort,
            }
            portsPool[StoVM.SwitchPort.UUID] = StoVM
        }
    }

    // 스위치를 먼저 처리
    for i:=range *LS{
        switchObj := (*LS)[i]
        err := o.AddExternSwitch(switchObj, portsPool)
        if err != nil {
            fmt.Printf("Error adding extern switch %s: %v\n", switchObj.UUID, err)
        }
    }

    // 라우터 처리
    for i:=range *LR{
        router := (*LR)[i]
        fmt.Printf("Processing router: UUID=%s, Name=%s, Ports=%v\n", 
            router.UUID, router.Name, router.Ports)
        
        err := o.UpdateDevices(router, portsPool)
        if err != nil {
            fmt.Printf("Error adding extern router %s: %v\n", router.UUID, err)
        }
    }

    fmt.Printf("Initialization complete: %d routers, %d switches\n", 
        len(o.ExternRouters), len(o.ExternSwitchs))
}

func (o* Operator)UpdateDevices (LR NBModel.LogicalRouter, ports map[string]externalmodel.NetInt)error {
    exR:= &externalmodel.ExternRouter{
        UUID:LR.UUID,
        InternalRouter: &LR,
        SubNetworks: make(map[string]externalmodel.NetInt),
    }

    for _, portUUID := range LR.Ports {
        if netInt, ok := ports[portUUID]; ok {
            switch port := netInt.(type) {
            case *externalmodel.RtoSwitchPort:  // 포인터로 type assertion
                port.ConnectedRouter = exR      // 원본 직접 수정
                fmt.Printf("Connected router-switch port %s to router %s\n", portUUID, LR.UUID)
                
                ip := port.RetriveAddress()
                if ip != "" {
                    exR.SubNetworks[ip] = port
                }
                if port.ConnectedSwitch == nil {
                    fmt.Printf("Warning: ConnectedSwitch is nil for port %s\n", portUUID)
                    continue
                }
                
                if port.ConnectedSwitch.InternalSwitch == nil {
                    fmt.Printf("Warning: InternalSwitch is nil for port %s\n", portUUID)
                    continue
                }
                
                ConnectedSwitch := port.ConnectedSwitch.InternalSwitch			
                for _, switchPortUUID := range ConnectedSwitch.Ports {
                    if switchPort, ok := ports[switchPortUUID]; ok {
                        Address := switchPort.RetriveAddress()
                        if Address != "" && Address != ip {
                            fmt.Printf("Adding switch port to ExternRouter: %s\n", Address)
                            exR.SubNetworks[Address] = switchPort
                        }
                    }
                }
            }
        }
    }

    o.ExternRouters[LR.UUID] = exR
    
    o.ExternRouters["10.5.15.4"] = exR
    
    return nil
}

func (o* Operator)AddExternSwitch (LS NBModel.LogicalSwitch, ports map[string]externalmodel.NetInt) error{
    exS:=&externalmodel.ExternSwitch{
        UUID: LS.UUID,
        InternalSwitch: &LS,
    }

    for _, portUUID := range LS.Ports {
        if netInt, ok := ports[portUUID]; ok {
            switch port := netInt.(type) {
            case *externalmodel.RtoSwitchPort:  
                port.ConnectedSwitch = exS     
                fmt.Printf("Connected RtoSwitchPort %s to switch %s\n", portUUID, LS.UUID)
                
            case *externalmodel.StoVMPort:      
                port.ConnectedSwitch = exS      
                fmt.Printf("Connected StoVMPort %s to switch %s\n", portUUID, LS.UUID)
            }
        }
    }

    o.ExternSwitchs[LS.UUID] = exS
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



    SP,err := o.AddSwitchAPort_Router(lsUUID, lrpuuid.String(), lspuuid.String())
    if err != nil {
        fmt.Printf("AddInterconnectR_S ERROR: Error in AddSwitchAPort_Router: %v\n", err)
        return err
    }

    natuuid,err:= util.UUIDGenerator()
    if err!=nil{
        fmt.Printf("AddInterconnectR_S ERROR: generating uuid error: %v\n", err)
        return err
    }
    routerPort, err := o.AddRouterPort(lrUUID, lrpuuid.String(), natuuid.String(), ip)
    if err != nil {
        fmt.Printf("AddInterconnectR_S ERROR: Error in AddRouterPort: %v\n", err)
        return err
    }




    InterPort:= &externalmodel.RtoSwitchPort{
		ConnectedRouter: o.ExternRouters[lrUUID],
		ConnectedSwitch: o.ExternSwitchs[lsUUID],
        RouterPort: routerPort,
		SwitchPort: SP,
        NatConnected: natuuid.String(),
	}


	externalmodel.AddNetInt(o.ExternRouters,ip , InterPort)

    return nil
}

func (o* Operator) InitialSetting()(error){
    EXTS_uuid,err:= o.AddSwitch()
    if err!=nil{
        return fmt.Errorf("failed to create external switch: %v", err)
    }
    fmt.Printf("Created EXTS_uuid: %s\n", EXTS_uuid)

    EXTR_uuid,err:=o.AddRouter(string(ROUTER))
    if err!=nil{
        return fmt.Errorf("failed to create external router: %v", err)
    }
    fmt.Printf("Created EXTR_uuid: %s\n", EXTR_uuid)

 


    // 포트 생성 및 연결
    lrpuuid,err:=util.UUIDGenerator()
    if err!=nil{
        return fmt.Errorf("lrpuuid generating error: %v", err)
    }
    lspuuid,err:=util.UUIDGenerator()
    if err!=nil{
        return fmt.Errorf("lspuuid generating error: %v", err)
    }



    SwitchPort, err := o.AddSwitchAPort_Router(EXTS_uuid, lrpuuid.String(), lspuuid.String())
    if err != nil {
        return fmt.Errorf("error in AddSwitchAPort_Router: %v", err)
    }
    natuuid, err := util.UUIDGenerator()
    if err != nil {
        return fmt.Errorf("generating uuid error: %v", err)
    }

    routerPort, err := o.AddRouterPort(EXTR_uuid, lrpuuid.String(), natuuid.String(), string(ROUTER))
    if err != nil {
        return fmt.Errorf("error in AddRouterPort: %v", err)
    }


        InterPort := &externalmodel.RtoSwitchPort{
        ConnectedRouter: o.ExternRouters[EXTR_uuid],
        ConnectedSwitch: o.ExternSwitchs[EXTS_uuid],
        SwitchPort: SwitchPort,
        RouterPort: routerPort,
        NatConnected: natuuid.String(),
    }

    o.ExternRouters[string(ROUTER)].SubNetworks[string(ROUTER)] = InterPort
    // 포인터로 생성


    // SubNetworks에 직접 추가

    err= o.ChassisInitializing(lrpuuid.String())
    if err!= nil{
        fmt.Printf("error adding chassis priority %v", err)
    }

    // Uplink 포트 생성
    br_EXTS_UUID,err := util.UUIDGenerator()
    if err!=nil{
        return fmt.Errorf("generating uuid error: br_exts_uuid")
    }

    newSP:=&externalmodel.SwitchPort{}
    operations:= []ovsdb.Operation{}
    
    ops,err:= newSP.Create(o.Client, br_EXTS_UUID.String(),"localnet", "unknown", map[string]string{
        "network_name": string(UPLINK),
    })
    if err != nil {
        return fmt.Errorf("creating switch port error %v", err)
    }
    operations = append(operations, ops...)

    request:= externalmodel.RequestControl{
        EXRList: o.ExternRouters,
        EXSList: o.ExternSwitchs,
        TargetUUID: EXTS_uuid,
        Client: o.Client,
    }

    ops,err = newSP.Connect(request)
    if err != nil {
        return fmt.Errorf("connecting switch port error %v", err)
    }
    operations = append(operations, ops...)

    result,err := o.Client.Transact(context.Background(),operations...)
    if err!=nil{
        return fmt.Errorf("transact error %v", err)
    }
    fmt.Println("Transaction result:", result)

    // 기본 경로 추가
    command := "ovn-nbctl" 
    args := []string{
        "lr-route-add",
        EXTR_uuid,
        "0.0.0.0/0",
        string(DEFAULT_GATEWAY),
    }
    
    cmd := exec.Command(command, args...) 
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout
    err = cmd.Run()
    if err != nil {
        return fmt.Errorf("error creating router command, %v", err)
    }

    return nil
}