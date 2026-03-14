package operation

import (
	"fmt"
	"os"
	"sync"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"

	"github.com/ovn-kubernetes/libovsdb/client"
)

type isOCcupied string

const (
	Occupied   isOCcupied = "occupied"
	UnOccupied isOCcupied = "unoccupied"
) // subnet 검사를 위해 사용, xxx.xxx.xxx.0에 저장 됨

type KNOWN_DEVICES string

const (
	UPLINK KNOWN_DEVICES = "UPLINK"
)

var (
	DEFAULT_GATEWAY KNOWN_DEVICES
	ROUTER          KNOWN_DEVICES
)

func init() {
	DEFAULT_GATEWAY = KNOWN_DEVICES(getEnvOrDefault("OVN_DEFAULT_GATEWAY", "10.5.15.1"))
	ROUTER = KNOWN_DEVICES(getEnvOrDefault("OVN_ROUTER_IP", "10.5.15.4"))
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

type Operator struct {
	mu            sync.Mutex
	Client        client.Client
	ExternRouters externalmodel.EXRList
	ExternSwitchs externalmodel.EXSList
}

func (o *Operator) Lock()   { o.mu.Lock() }
func (o *Operator) Unlock() { o.mu.Unlock() }

func (o *Operator) IPMapToDev(IP string) externalmodel.NetInt {
	list := externalmodel.GetNetInt(o.ExternRouters, IP)
	if len(list) > 0 {
		return list[0]
	}
	return nil
}

func (o *Operator) SwitchesPortConnect(uuids []string, IP string, VMUUID string, VMMac string) error {
	for _, uuid := range uuids {
		if _, err := o.AddSwitchAPort(uuid, IP, VMUUID, VMMac); err != nil {
			return err
		}
	}
	return nil
}

func (o *Operator) findDevByUUID(uuid string) (any, error) {
	dev, ok := o.ExternRouters[uuid]
	if !ok {
		dev, ok := o.ExternSwitchs[uuid]
		if !ok {
			return nil, fmt.Errorf("no such device for uuid")
		}
		return dev, nil
	}
	return dev, nil
}
