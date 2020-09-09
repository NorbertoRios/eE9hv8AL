package core

import (
	"queclink-go/base.device.service/report"
	"queclink-go/core/immobilizerq600"
	"queclink-go/queclinkreport/fields"
	"queclink-go/queclinkreport/gv600/resp600"
)

//QDeviceGV600 base queclink device
type QDeviceGV600 struct {
	QDevice
	ts *report.TemperatureSensors
}

//Initialize device
func (device *QDeviceGV600) Initialize(uniqueID string) error {
	err := device.QDevice.Initialize(uniqueID)
	if err == nil {
		device.loadTemperatureSensors()
	}
	return err
}

//SetDefaults for device
func (device *QDeviceGV600) SetDefaults(deviceType int) {
	device.handleRspMessage = device.rspMessageHandler
	device.handleEvtMessage = device.evtMessageHandler
	device.Statistic = NewDeviceStatistic()
	device.Immobilizer = immobilizerq600.Initialize(device)
	device.Type = deviceType
	device.Self = device
}

//UpdateIgnitionState for ql600 device
func (device *QDeviceGV600) UpdateIgnitionState(message report.IMessage) {
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

//UpdateCurrentState assign values from message to current state and fill up missing fields for message
func (device *QDeviceGV600) UpdateCurrentState(message report.IMessage) {
	device.setCachedValues(message)
	device.validateSupply(message)
	message.RemoveKey(fields.OWireSensors)
	device.QDevice.UpdateCurrentState(message)
	if message.LocationMessage() {
		device.UpdateIgnitionState(message)
	}
}

func (device *QDeviceGV600) validateSupply(message report.IMessage) {
	if message.GetIntValue(fields.Supply, 0) == 0 {
		supply2 := message.GetIntValue(fields.Supply2, 0)
		if supply2 > 0 {
			message.SetValue(fields.Supply, supply2)
		}
	}
}

func (device *QDeviceGV600) setCachedValues(message report.IMessage) {
	if (message.MessageType() == "+RSP" || message.MessageType() == "+BSP") &&
		message.EventCode() == resp600.GTERI {
		ts := message.GetTemperatureSensors()
		if ts != nil {
			if ts.Sensor1 != nil {
				if ts.Sensor1.TemperatureValue < -80 || ts.Sensor1.TemperatureValue > 125 {
					if device.ts != nil && device.ts.Sensor1 != nil {
						ts.Sensor1 = device.ts.Sensor1
					} else {
						ts.Sensor1 = nil
					}
				}
			}
			if ts.Sensor2 != nil {
				if ts.Sensor2.TemperatureValue < -80 || ts.Sensor2.TemperatureValue > 125 {
					if device.ts != nil && device.ts.Sensor2 != nil {
						ts.Sensor2 = device.ts.Sensor2
					} else {
						ts.Sensor2 = nil
					}
				}
			}

			if ts.Sensor3 != nil {
				if ts.Sensor3.TemperatureValue < -80 || ts.Sensor3.TemperatureValue > 125 {
					if device.ts != nil && device.ts.Sensor3 != nil {
						ts.Sensor3 = device.ts.Sensor3
					} else {
						ts.Sensor3 = nil
					}
				}
			}

			if ts.Sensor4 != nil {
				if ts.Sensor4.TemperatureValue < -80 || ts.Sensor4.TemperatureValue > 125 {
					if device.ts != nil && device.ts.Sensor4 != nil {
						ts.Sensor4 = device.ts.Sensor4
					} else {
						ts.Sensor4 = nil
					}
				}
			}
			device.ts = ts
		}

	} else {
		if (message.MessageType() == "+RSP" || message.MessageType() == "+BSP") &&
			message.EventCode() == resp600.GTFRI {
			device.ts = nil
			message.SetTemperatureSensors(nil)
		}

		if device.ts != nil {
			message.SetTemperatureSensors(&report.TemperatureSensors{
				Sensor1: device.ts.Sensor1,
				Sensor2: device.ts.Sensor2,
				Sensor3: device.ts.Sensor3,
				Sensor4: device.ts.Sensor4,
			})
		}
	}
}

func (device *QDeviceGV600) loadTemperatureSensors() {
	if device.Activity == nil || device.Activity.LastMessage == "" {
		return
	}
	message, err := report.UnMarshalMessage(device.Activity.LastMessage)
	if err != nil {
		return
	}
	if message.TemperatureSensors != nil {
		device.ts = message.TemperatureSensors
	}
}
