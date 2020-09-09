package core

import (
	"queclink-go/base.device.service/report"
	"queclink-go/queclinkreport/fields"
)

//QDeviceGV55 base queclink device
type QDeviceGV55 struct {
	QDevice
}

//UpdateCurrentState assign values from message to current state and fill up missing fields for message
func (device *QDeviceGV55) UpdateCurrentState(message report.IMessage) {
	device.QDevice.UpdateCurrentState(message)
	if message.LocationMessage() {
		device.UpdateIgnitionState(message)
	}
}

//UpdateIgnitionState for ql300 device
func (device *QDeviceGV55) UpdateIgnitionState(message report.IMessage) {
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
func (device *QDeviceGV55) SetDefaults(deviceType int) {
	device.QDevice.SetDefaults(deviceType)
	device.Self = device
}
