package immobilizerq600

import (
	"strings"
	"time"

	"queclink-go/base.device.service/core"

	"queclink-go/core/immobilizer"
	"queclink-go/queclinkreport/devicetypes"
)

//Q600Handler manages device immobilizer state
type Q600Handler struct {
	immobilizer.Handler
}

//SendAPIImmobilizerCommand updates device immobilizer state
func (immo *Q600Handler) SendAPIImmobilizerCommand(callbackID, port, state string, ttl int, trigger string, safetyOption bool) string {
	immo.CurrentStatus = immobilizer.Created
	t := strings.ToUpper(trigger)
	if trigger == "" {
		t = "HIGH"
	}
	cmd := &immobilizer.CommandData{
		CallbackID:   callbackID,
		CreatedAt:    time.Now().UTC(),
		TTL:          ttl,
		Port:         strings.ToUpper(port),
		State:        strings.ToUpper(state),
		Trigger:      t,
		SafetyOption: safetyOption,
		Status:       immobilizer.Created,
	}

	cfg, output := immobilizer.GetImmobilizerCommand(immo.Device.GetIdentity(), cmd.Port, cmd.State, cmd.Trigger, devicetypes.GV600W)
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

//Initialize creates new instance of immobiliser handler
func Initialize(device core.IDevice) *Q600Handler {
	h := &Q600Handler{}
	h.Device = device
	return h
}
