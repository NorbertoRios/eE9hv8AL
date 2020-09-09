package core

import (
	"queclink-go/base.device.service/report"
	"queclink-go/queclinkreport/fields"
	"queclink-go/queclinkreport/gv300w/evt300w"
	"queclink-go/queclinkreport/gv300w/resp300w"
)

//QDeviceGV300W base queclink device
type QDeviceGV300W struct {
	QDevice
	fuelLevel       float32
	virtualIgnition string
}

//ProcessLocationMessage handle location message
func (device *QDeviceGV300W) ProcessLocationMessage(message report.IMessage) (bool, error) {

	if (message.MessageType() == "+RSP" || message.MessageType() == "+BSP") && message.EventCode() == resp300w.GTERI {
		ifv, found := message.GetValue(fields.FuelLevel)
		if found {
			fv, valid := ifv.(float32)
			if valid {
				device.fuelLevel = fv
			}
		}
	} else {
		if device.fuelLevel > 0 {
			message.SetValue(fields.FuelLevel, device.fuelLevel)
		}
	}

	return device.QDevice.ProcessLocationMessage(message)
}

//UpdateCurrentState assign values from message to current state and fill up missing fields for message
func (device *QDeviceGV300W) UpdateCurrentState(message report.IMessage) {
	device.QDevice.UpdateCurrentState(message)
	if message.LocationMessage() {
		device.UpdateIgnitionState(message)
	}
}

//UpdateIgnitionState for ql300 device
func (device *QDeviceGV300W) UpdateIgnitionState(message report.IMessage) {
	if message.MessageType() == "+EVT" || message.MessageType() == "+BVT" {
		switch message.EventCode() {
		case evt300w.GTVGN:
			device.virtualIgnition = "On"
			break
		case evt300w.GTVGF:
			device.virtualIgnition = "Off"
			break
		case evt300w.GTIGN:
		case evt300w.GTIGF:
			device.virtualIgnition = ""
			break
		}
	}

	iv, found := message.GetValue(fields.DigitalInputStatus)
	if !found {
		return
	}
	diStatus, valid := iv.(byte)
	if !valid {
		return
	}

	if (diStatus & 1) == 1 {
		device.virtualIgnition = ""
	}

	if device.virtualIgnition != "" {
		if device.virtualIgnition == "On" {
			message.SetValue(fields.IgnitionState, byte(1))
			diStatus = diStatus | 1
			message.SetValue(fields.DigitalInputStatus, byte(diStatus|1))
		}

		if device.virtualIgnition == "Off" {
			message.SetValue(fields.IgnitionState, byte(0))
			diStatus = diStatus & 0xFE
			message.SetValue(fields.DigitalInputStatus, diStatus)

		}
	}
	if (diStatus & 1) == 1 {
		device.Activity.Ignition = "On"
	} else {
		device.Activity.Ignition = "Off"
	}
	message.SetValue("VirtualIgnition", device.virtualIgnition)
}

//SetDefaults for device
func (device *QDeviceGV300W) SetDefaults(deviceType int) {
	device.QDevice.SetDefaults(deviceType)
	device.Self = device
}

//UpdateEventDependentStates ...
func (device *QDeviceGV300W) UpdateEventDependentStates(message report.IMessage) {
	if device.Activity.PowerState == "Backup battery" && message.GetIntValue("Supply", 0) > 8000 {
		device.Activity.PowerState = "Powered"
	}
	message.SetValue("PowerState", device.Activity.PowerState)
	if device.fuelLevel == 0.0 {
		message.RemoveKey(fields.FuelLevel)
	}
	device.Activity.FuelLevel = device.fuelLevel
}
