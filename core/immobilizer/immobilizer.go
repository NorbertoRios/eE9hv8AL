package immobilizer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"queclink-go/base.device.service/config"
	"queclink-go/base.device.service/core"
	"queclink-go/base.device.service/core/immobilizer"
	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
	"queclink-go/qconfig"
)

const (
	//None status
	None = -1
	//Created status
	Created = 0
	//Sended status
	Sended = 1
	//RtoSended status
	RtoSended = 2
	//Confirmed status
	Confirmed = 3
)

//Handler manages device immobilizer state
type Handler struct {
	Device             core.IDevice
	LastUpdateDateTime time.Time
	SentCommand        core.IConfigurationItem
	CurrentStatus      int
}

//RemoveExpiredImmCommands from queue
func (immo *Handler) RemoveExpiredImmCommands() {
	if immo.SentCommand != nil {
		command, valid := immo.SentCommand.(*core.ConfigurationItemAPI)
		if valid {
			d := time.Duration(command.TTL) * time.Second
			if command.CreationTime.Add(d).Before(time.Now().UTC()) {
				immo.SentCommand = nil
			}
		}

		immoCommand, valid := immo.SentCommand.(*core.ConfigurationItemImmobilizerAPI)
		if valid {
			d := time.Duration(immoCommand.TTL) * time.Second
			if immoCommand.CreationTime.Add(d).Before(time.Now().UTC()) {
				immo.SentCommand = nil
			}
		}

	}
}

func (immo *Handler) generateImmobilizerOutput(port, state, trigger string) (portNum byte, portValue byte, found bool) {
	portNum, found = immobilizer.StrPortToInt(port)
	if !found {
		return portNum, 0, found
	}
	portValue, found = immobilizer.GetPortValue(state, trigger)
	return portNum, portValue, found
}

//IsImmobilizerOutputsEquals checks outputs are equals
func (immo *Handler) IsImmobilizerOutputsEquals(outputs []byte) bool {
	for i, v := range outputs {
		bs := utils.BitIsSet(int64(immo.Device.GetActivity().Relay), uint(i))
		if !(bs == (v == 1)) {
			return false
		}
	}
	return true
}

func (immo *Handler) generateImmobilizerConfigItem(cfg string, safetyOption bool) *core.ConfigurationItemImmobilizerAPI {
	if cfg != "" {
		result := &core.ConfigurationItemImmobilizerAPI{}
		result.Command = cfg
		result.MessageType = "8" //AT+GTOUT
		result.SafetyOption = safetyOption
		result.Type = "api_immobilizer_request"
		result.CreationTime = time.Now().UTC()
		return result
	}
	return nil
}

func (immo *Handler) generateImmStateCheckConfigItem() *core.ConfigurationItemAPI {
	cfg := &core.ConfigurationItemAPI{}
	cfg.Command = "AT+GTRTO=gv55,A,,,,,,FFFF$"
	cfg.MessageType = "16"
	return cfg
}

func (immo *Handler) immItemSended(callbackID, code string, result bool) {
	immo.RemoveExpiredImmCommands()
	if immo.SentCommand == nil {
		return
	}
	cID := ""
	var apiItem *core.ConfigurationItemAPI

	if c, v := immo.SentCommand.(*core.ConfigurationItemAPI); v {
		cID = c.CallbackID
		apiItem = c
	}

	if c, v := immo.SentCommand.(*core.ConfigurationItemImmobilizerAPI); v {
		cID = c.CallbackID
		apiItem = &c.ConfigurationItemAPI
	}
	if cID == callbackID {
		item := immo.generateImmStateCheckConfigItem()
		item.CallbackID = callbackID
		item.TTL = apiItem.TTL
		item.CreationTime = apiItem.CreationTime
		item.Type = "api_request"
		item.OnCommandSent(immo.immRtoSended)
		immo.CurrentStatus = Sended
		immo.SentCommand = item
		immo.Update()
	}
}

func (immo *Handler) immRtoSended(callbackID, code string, result bool) {
	immo.RemoveExpiredImmCommands()
	if immo.SentCommand == nil {
		return
	}
	cID := ""
	if c, v := immo.SentCommand.(*core.ConfigurationItemAPI); v {
		cID = c.CallbackID
	}

	if c, v := immo.SentCommand.(*core.ConfigurationItemImmobilizerAPI); v {
		cID = c.CallbackID
	}
	if cID == callbackID {
		if immo.CurrentStatus == Sended {
			immo.CurrentStatus = RtoSended
			immo.CheckStatusOfImmobilizerCommands()
		}
	}
}

//SendAPIImmobilizerCommand updates device immobilizer state
func (immo *Handler) SendAPIImmobilizerCommand(callbackID, port, state string, ttl int, trigger string, safetyOption bool) string {
	immo.CurrentStatus = Created
	t := strings.ToUpper(trigger)
	if trigger == "" {
		t = "HIGH"
	}
	cmd := &CommandData{
		CallbackID:   callbackID,
		CreatedAt:    time.Now().UTC(),
		TTL:          ttl,
		Port:         strings.ToUpper(port),
		State:        strings.ToUpper(state),
		Trigger:      t,
		SafetyOption: safetyOption,
		Status:       Created,
	}

	cfg, output := GetImmobilizerCommand(immo.Device.GetIdentity(), cmd.Port, cmd.State, cmd.Trigger, 0)
	immo.RemoveExpiredImmCommands()
	cmd.Command = cfg

	//Check empty command queue and current state equal to needed
	if immo.SentCommand == nil && immo.IsImmobilizerOutputsEquals(output) {
		core.SendConfigurationAPIResponse(callbackID, "Actual", true)
		return cfg
	}
	immo.UpdateSendingItem(cmd)
	return cfg
}

//UpdateSendingItem state
func (immo *Handler) UpdateSendingItem(cmd *CommandData) {
	item := immo.generateImmobilizerConfigItem(cmd.Command, cmd.SafetyOption)
	if item == nil {
		return
	}

	item.CallbackID = cmd.CallbackID
	item.TTL = cmd.TTL
	item.OnCommandSent(immo.immItemSended)
	log.Println("[Immobilizer]Enqueue imm command:", cmd.Command, " for device:", immo.Device.GetIdentity())
	immo.SentCommand = item
	immo.Device.SendString("AT+GTRTO=gv55,A,,,,,,FFFF$")
}

func (immo *Handler) checkSendImmAbility(item *core.ConfigurationItemImmobilizerAPI) bool {
	if item == nil {
		return true
	}

	if immo.Device.GetActivity().Ignition == "Off" &&
		immo.Device.GetActivity().MessageTime.Add(600*time.Second).After(time.Now().UTC()) {
		log.Println("[Immobilizer]checkSendImmAbility. Result:true; Ignition is:", immo.Device.GetActivity().Ignition)
		return true
	}
	if item.SafetyOption {
		if immo.Device.GetActivity().Ignition == "On" {
			log.Println("[Immobilizer]checkSendImmAbility. Result:false; Ignition is:", immo.Device.GetActivity().Ignition)
			return false
		}
	}
	if immo.Device.GetActivity().MessageTime.Add(600 * time.Second).Before(time.Now().UTC()) {
		return true
	}
	return true
}

//Update immobilizer state
func (immo *Handler) Update() {

	log.Println("[Immobilizer]Update. Device:", immo.Device.GetIdentity())

	if immo.SentCommand == nil {
		return
	}

	if immo.SentCommand.GetSentAt().Add(10 * time.Second).After(time.Now().UTC()) {
		return
	}

	log.Println("[Immobilizer]Update. Device:", immo.Device.GetIdentity(), ";sending", immo.SentCommand.GetCommand())
	c, v := immo.SentCommand.(*core.ConfigurationItemImmobilizerAPI)
	if (v && immo.checkSendImmAbility(c)) || !v {
		immo.SentCommand.SetSentAt(time.Now().UTC())
		immo.Device.SendString(immo.SentCommand.GetCommand())
		immo.LastUpdateDateTime = time.Now().UTC()
		log.Println("[Immobilizer]Update. Device:", immo.Device.GetIdentity(), ";sent", immo.SentCommand.GetCommand())
	}
}

//CheckStatusOfImmobilizerCommands checks immo status and send response to facade
func (immo *Handler) CheckStatusOfImmobilizerCommands() {
	log.Println("[Immobilizer]CheckStatusOfImmobilizerCommands. Device:", immo.Device.GetIdentity(), "; command is nil:", immo.SentCommand == nil, "current state:", immo.CurrentStatus)
	if immo.SentCommand != nil {
		if immo.CurrentStatus != RtoSended {
			immo.Update()
		} else if immo.CurrentStatus == RtoSended {
			currentStates, _ := GetFacadeOutputsState(immo.Device.GetIdentity())
			outputs := GetOutputArray(currentStates)
			if immo.IsImmobilizerOutputsEquals(outputs) {
				if outputs != nil {
					for _, state := range currentStates.Items {
						core.SendConfigurationAPIResponse(state.CallbackID, "Done", true)
					}
				}
				immo.SentCommand = nil
				immo.CurrentStatus = None
			}
		}

	}
	immo.RemoveExpiredImmCommands()
}

//UpdateMessageSent confirms message is sent
func (immo *Handler) UpdateMessageSent(message report.IMessage) {
	if immo.SentCommand == nil {
		return
	}
	messageType := fmt.Sprint(message.EventCode())
	if immo.SentCommand.GetMessageType() == messageType {
		immo.Device.SendSystemMessageS(fmt.Sprintf("Api request \"%v\" is done", immo.SentCommand.GetCommand()), config.Config.GetBase().Rabbit.OrleansDebugRoutingKey)
		if c, v := immo.SentCommand.(*core.ConfigurationItemAPI); v {
			c.SendResponse("Done", true)
		}
		if c, v := immo.SentCommand.(*core.ConfigurationItemImmobilizerAPI); v {
			c.SendResponse("Done", true)
		}

	} else if immo.CurrentStatus != RtoSended {
		immo.Update()
	}
}

//GetImmobilizerCommand returns immobilizer command generated from parameters
func GetImmobilizerCommand(identity string, port string, state string, trigger string, deviceType int) (string, []byte) {

	out0 := 0
	out1 := 1
	out2 := 2
	out3 := 3
	outputs := GenerateImmobilizerOutputs(identity, port, state, trigger)
	if deviceType == 252 || deviceType == 237 { //GV600
		return fmt.Sprintf("AT+GTOUT=gv55,%v,0,0,%v,0,0,%v,0,0,%v,0,0,F,0,0,0,0,FFFF$", outputs[out0], outputs[out1], outputs[out2], outputs[out3]), outputs
	}
	return fmt.Sprintf("AT+GTOUT=gv55,%v,0,0,%v,0,0,,,,0,,,,,,,FFFF$", outputs[out0], outputs[out1]), outputs
}

//GenerateImmobilizerOutputs generates output states
func GenerateImmobilizerOutputs(identity string, port string, state string, trigger string) []byte {
	currentStates, _ := GetFacadeOutputsState(identity)
	outputs := GetOutputArray(currentStates)
	if port != "" {
		value, res := immobilizer.GetPortValue(state, trigger)
		if !res {
			return outputs
		}
		p, res := immobilizer.StrPortToInt(port)
		if !res {
			return outputs
		}
		outputs[p] = value
	}
	return outputs
}

//GetOutputArray converts StatesResponse to slice
func GetOutputArray(currentStates *StatesResponse) []byte {
	outputs := []byte{0, 0, 0, 0}

	if len(currentStates.Items) > 0 {
		for _, out := range currentStates.Items {
			value, res := immobilizer.GetPortValue(out.State, out.Trigger)
			if !res {
				continue
			}
			p, res := immobilizer.StrPortToInt(out.Port)
			if !res {
				continue
			}
			outputs[p] = value
		}
	}
	return outputs
}

//GetFacadeOutputsState returns response from facade with outputs
func GetFacadeOutputsState(identity string) (*StatesResponse, error) {
	outStatuses := &StatesResponse{}
	resp, err := http.Get(fmt.Sprintf("%v/device/analytics/service_immobilizer_states?identity=%v", config.Config.(*qconfig.QConfiguration).DeviceFacadeHost, identity))
	if err != nil {
		return outStatuses, fmt.Errorf("[Immobilizer]GetFacadeOutputsState. Error response. Error:%v; ", err.Error())
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 404 {
		return outStatuses, nil
	}

	if resp.StatusCode != 200 {
		return outStatuses, fmt.Errorf("[Immobilizer]GetFacadeOutputsState. Invalid response. Code:%v; Content:%v", resp.StatusCode, contents)
	}

	if err != nil {
		return outStatuses, fmt.Errorf("[Immobilizer]GetFacadeOutputsState. Content read error:%v ", err.Error())
	}

	log.Println("[Immobilizer]GetFacadeOutputsState. response:", string(contents))
	err = json.Unmarshal(contents, outStatuses)
	if err != nil {
		return outStatuses, err
	}
	return outStatuses, nil
}

//Initialize creates new instance of immobiliser handler
func Initialize(device core.IDevice) *Handler {
	return &Handler{
		Device: device,
	}
}

//CommandData struct to describe immobilizer state
type CommandData struct {
	CallbackID    string
	Port          string
	State         string
	Trigger       string
	SafetyOption  bool
	Status        byte
	CreatedAt     time.Time
	TTL           int
	Command       string
	onCommandSent func(callbackId string, code string, result bool)
}

//StatesResponse set of outputs state
type StatesResponse struct {
	Items []StateResponse
}

//StateResponse utput state
type StateResponse struct {
	Port       string `json:"port"`
	State      string `json:"state"`
	Trigger    string `json:"trigger"`
	CallbackID string `json:"CallbackId"`
}
