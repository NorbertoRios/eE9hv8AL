package core

import (
	"queclink-go/base.device.service/report"
	"queclink-go/queclinkreport/fields"
	"queclink-go/queclinkreport/gv55w/evt55w"
)

//QDeviceGV55W base queclink device
type QDeviceGV55W struct {
	QDevice
	virtualIgnition string
}

//UpdateCurrentState assign values from message to current state and fill up missing fields for message
func (device *QDeviceGV55W) UpdateCurrentState(message report.IMessage) {
	device.QDevice.UpdateCurrentState(message)
	if message.LocationMessage() {
		device.UpdateIgnitionState(message)
	}
}

//UpdateIgnitionState for ql55w device
func (device *QDeviceGV55W) UpdateIgnitionState(message report.IMessage) {
	if message.MessageType() == "+EVT" || message.MessageType() == "+BVT" {
		switch message.EventCode() {
		case evt55w.GTVGN:
			device.virtualIgnition = "On"
			break
		case evt55w.GTVGF:
			device.virtualIgnition = "Off"
			break
		case evt55w.GTIGN:
		case evt55w.GTIGF:
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
func (device *QDeviceGV55W) SetDefaults(deviceType int) {
	device.QDevice.SetDefaults(deviceType)
	device.Self = device
}

//Initialize device
func (device *QDeviceGV55W) Initialize(uniqueID string) error {
	device.QDevice.Initialize(uniqueID)
	strMsg := device.Activity.LastMessage
	if strMsg != "" {
		message, err := report.UnMarshalMessage(strMsg)
		if err == nil {
			if ivIgn, found := message.GetValue("VirtualIgnition"); found {
				vIgnState, valid := ivIgn.(string)
				if valid {
					device.virtualIgnition = vIgnState
				}
			}
		}
	}
	return nil
}
