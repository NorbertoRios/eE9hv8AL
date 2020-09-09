package core

import (
	"queclink-go/core/immobilizer"
	"queclink-go/queclinkreport/devicetypes"
)

//Queclink600DM overrided version of core.DeviceManager
type Queclink600DM struct {
	QueclinkDM
}

//GetImmobilizerCommand return immobilizer command for network usage
func (manager *Queclink600DM) GetImmobilizerCommand(identity string, port string, state string, trigger string) string {
	command, _ := immobilizer.GetImmobilizerCommand(identity, port, state, trigger, devicetypes.GV600W)
	return command
}

//GetSmsImmobilizerCommand return immobilizer command for cell usage
func (manager *Queclink600DM) GetSmsImmobilizerCommand(identity string, port string, state string, trigger string) string {
	command, _ := immobilizer.GetImmobilizerCommand(identity, port, state, trigger, devicetypes.GV600W)
	return command
}
