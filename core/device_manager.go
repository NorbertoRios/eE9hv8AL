package core

import (
	"fmt"

	"queclink-go/core/immobilizer"
	"queclink-go/queclinkreport/devicetypes"
	"queclink-go/queclinkreport/fields"

	"queclink-go/base.device.service/comm"
	"queclink-go/base.device.service/core"
	"queclink-go/base.device.service/report"
)

//QueclinkDM overrided version of core.DeviceManager
type QueclinkDM struct {
	core.DeviceManager
}

//InitializeDevice creates new device
func (manager *QueclinkDM) InitializeDevice(c comm.IChannel, message report.IMessage) (core.IDevice, error) {

	deviceType := int(0)
	idt, found := message.GetValue(fields.DeviceType)
	if !found {
		return nil, fmt.Errorf("DeviceType not found")
	}

	if _, valid := idt.(byte); !valid {
		return nil, fmt.Errorf("Invalid DeviceType data type")
	}

	var device core.IDevice
	deviceType = int(idt.(byte))
	switch deviceType {
	case devicetypes.GV320:
		d := &QDeviceGV300{}
		d.SetDefaults(deviceType)
		device = d
		break
	case devicetypes.GV55,
		devicetypes.GV55N:
		d := &QDeviceGV55{}
		d.SetDefaults(deviceType)
		device = d
		break
	case devicetypes.GV55Lite,
		devicetypes.GV55NLite:
		d := &QDevice{}
		d.SetDefaults(deviceType)
		device = d
		break
	case devicetypes.GV75,
		devicetypes.GV75W:
		d := &QDeviceGV75{}
		d.SetDefaults(deviceType)
		device = d
		break
	case devicetypes.GV55W:
		d := &QDeviceGV55W{}
		d.SetDefaults(deviceType)
		device = d
		break
	case devicetypes.GV600W, devicetypes.GV600Fake, devicetypes.GV600MG:
		d := &QDeviceGV600{}
		d.SetDefaults(deviceType)
		device = d
		break
	case devicetypes.GV350MA:
		d := &QDeviceGV350{}
		d.SetDefaults(deviceType)
		device = d
		break
	case devicetypes.GV300W:
		d := &QDeviceGV300W{}
		d.SetDefaults(deviceType)
		device = d
		break
	}

	device.SetChannel(c)
	err := device.Initialize(message.UniqueID())
	manager.AddRegisteredDevice(device)
	return device, err
}

//GetImmobilizerCommand return immobilizer command for network usage
func (manager *QueclinkDM) GetImmobilizerCommand(identity string, port string, state string, trigger string) string {
	command, _ := immobilizer.GetImmobilizerCommand(identity, port, state, trigger, devicetypes.GV55)
	return command
}

//GetSmsImmobilizerCommand return immobilizer command for cell usage
func (manager *QueclinkDM) GetSmsImmobilizerCommand(identity string, port string, state string, trigger string) string {
	command, _ := immobilizer.GetImmobilizerCommand(identity, port, state, trigger, devicetypes.GV55)
	return command
}

//InitializeUDPDevice creates new device
func (manager *QueclinkDM) InitializeUDPDevice(c comm.IChannel, message report.IMessage) (core.IDevice, error) {
	device, err := manager.InitializeDevice(c, message)
	if err != nil {
		return nil, err
	}

	device.GetChannel().(*comm.UDPChannel).OnProcessMessage(device.ProcessMessage)
	manager.DeviceWorkers.AddDevice(device)
	return device, nil
}
