package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	externalmodel "github.com/kwonkwonn/ovn-go-cms/ovs/externalModel"
	"github.com/kwonkwonn/ovn-go-cms/ovs/operation"
	"github.com/kwonkwonn/ovn-go-cms/ovs/util"
)

func (h *Handler) CreateNewVm(w http.ResponseWriter, r *http.Request) {
	h.Operator.Lock()
	defer h.Operator.Unlock()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("CreateNewVm: read body error: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	request := &NewInstanceRequeset{}
	if err = json.Unmarshal(body, request); err != nil {
		log.Printf("CreateNewVm: unmarshal error: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var RtoSInterfaceIP = request.RequestSubnet + "1"
	RtoSInterface := externalmodel.GetNetInt(h.Operator.ExternRouters, RtoSInterfaceIP)

	var swUUID string
	newvifIP := externalmodel.FindRemainIP(h.Operator.ExternRouters, request.RequestSubnet, externalmodel.VIF)

	mac, err := util.MacGenerator()
	if err != nil {
		log.Printf("CreateNewVm: mac generating error: %v", err)
		writeError(w, http.StatusInternalServerError, "mac generating error")
		return
	}
	InstUUID, err := util.UUIDGenerator()
	if err != nil {
		log.Printf("CreateNewVm: uuid generating error: %v", err)
		writeError(w, http.StatusInternalServerError, "uuid generating error")
		return
	}

	if len(RtoSInterface) == 0 {
		log.Println("CreateNewVm: request for new subnet, creating new router port")
		routerUUID := h.Operator.ExternRouters[string(operation.ROUTER)].UUID
		if routerUUID == "" {
			writeError(w, http.StatusInternalServerError, "router not found")
			return
		}

		swUUID, err = h.Operator.AddSwitch()
		if err != nil {
			log.Printf("CreateNewVm: add switch error: %v", err)
			writeError(w, http.StatusInternalServerError, "add switch error")
			return
		}
		if err = h.Operator.AddInterconnectR_S(swUUID, routerUUID, RtoSInterfaceIP); err != nil {
			log.Printf("CreateNewVm: add interconnect error: %v", err)
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		swUUID = RtoSInterface[0].(*externalmodel.RtoSwitchPort).ConnectedSwitch.UUID
	}

	if err = h.Operator.SwitchesPortConnect([]string{swUUID}, newvifIP, InstUUID.String(), mac); err != nil {
		log.Printf("CreateNewVm: switch port connect error: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	result := &NewInstanceResult{
		MacAddress: mac,
		IP:         newvifIP,
		IfaceID:    InstUUID.String(),
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) DeleteInstance(w http.ResponseWriter, r *http.Request) {
	ip := r.PathValue("ip")
	if ip == "" {
		writeError(w, http.StatusBadRequest, "missing ip path parameter")
		return
	}
	h.deleteInstance(w, ip)
}

func (h *Handler) DelNet(w http.ResponseWriter, r *http.Request) {
	h.Operator.Lock()
	defer h.Operator.Unlock()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("DelNet: read body error: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	request := &DelInstanceRequest{}
	if err = json.Unmarshal(body, request); err != nil {
		log.Printf("DelNet: unmarshal error: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.deleteInstance(w, request.IP)
}

func (h *Handler) deleteInstance(w http.ResponseWriter, ip string) {
	h.Operator.Lock()
	defer h.Operator.Unlock()

	NetSignifier, err := util.GetNetWorkSignifier(ip)
	if err != nil {
		log.Printf("deleteInstance: network signifier parse error: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	NetInt := externalmodel.GetNetInt(h.Operator.ExternRouters, ip)
	if len(NetInt) == 0 {
		writeError(w, http.StatusNotFound, "no such switch exist")
		return
	}

	if _, ok := NetInt[0].(*externalmodel.RtoSwitchPort); ok {
		writeError(w, http.StatusConflict, "cannot delete switch port, connected to router")
		return
	}

	Port, ok := NetInt[0].(*externalmodel.StoVMPort)
	if !ok {
		writeError(w, http.StatusNotFound, "no such switch port exist")
		return
	}

	SwitchPort := Port.ConnectedSwitch
	if SwitchPort == nil {
		writeError(w, http.StatusInternalServerError, "switch port not connected")
		return
	}

	if err = h.Operator.DelSwitchPort(ip); err != nil {
		log.Printf("deleteInstance: del switch port error: %v", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	result := &DelInstanceResult{Detail: "delete switch port operation"}

	delete(h.Operator.ExternRouters[string(operation.ROUTER)].SubNetworks, ip)
	nets := externalmodel.GetAllVIF(h.Operator.ExternRouters, NetSignifier)
	if len(nets) == 0 {
		log.Println("deleteInstance: deleting connected switch")
		if err = h.Operator.DelSwitch(SwitchPort.UUID); err != nil {
			log.Printf("deleteInstance: del switch error: %v", err)
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err = h.Operator.DelRouterPort(NetSignifier + "1"); err != nil {
			log.Printf("deleteInstance: del router port error: %v", err)
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("del router port error: %v", err))
			return
		}
		result.Detail = "delete switch and router port success"
		delete(h.Operator.ExternRouters[string(operation.ROUTER)].SubNetworks, NetSignifier+"1")
	}

	result.Detail = fmt.Sprintf("%s, switch port deleted", result.Detail)
	log.Println("deleteInstance: delete vm success")
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	h.Operator.Lock()
	defer h.Operator.Unlock()

	h.Operator.DeleteAll()
	writeJSON(w, http.StatusOK, map[string]string{"detail": "work done"})
}


