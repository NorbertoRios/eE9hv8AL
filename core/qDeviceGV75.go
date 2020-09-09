package core

import (
	"queclink-go/base.device.service/report"
	"queclink-go/queclinkreport/fields"
	"queclink-go/queclinkreport/gv75/resp75"
)

//QDeviceGV75 base queclink device
type QDeviceGV75 struct {
	QDevice
	fuelLevel int32
}

//ProcessLocationMessage handle location message
func (device *QDeviceGV75) ProcessLocationMessage(message report.IMessage) (bool, error) {

	if (message.MessageType() == "+RSP" || message.MessageType() == "+BSP") && message.EventCode() == resp75.GTERI {
		ifv, found := message.GetValue(fields.FuelLevel)
		if found {
			fv, valid := ifv.(int32)
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
func (device *QDeviceGV75) UpdateCurrentState(message report.IMessage) {
	device.QDevice.UpdateCurrentState(message)
	if message.LocationMessage() {
		device.UpdateIgnitionState(message)
	}
}

//UpdateIgnitionState for ql300 device
func (device *QDeviceGV75) UpdateIgnitionState(message report.IMessage) {
	iv, found := message.GetValue(fields.DigitalInputStatus)
	if !found {
		return
	}
	diStatus, valid := iv.(byte)
	if !valid {
		return
	}

	if (diStatus & 1) == 1 {
		device.Activity.Ignition = "On"
	} else {
		device.Activity.Ignition = "Off"
	}
}

//SetDefaults for device
func (device *QDeviceGV75) SetDefaults(deviceType int) {
	device.QDevice.SetDefaults(deviceType)
	device.Self = device
}
