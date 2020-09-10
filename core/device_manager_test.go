package core

import (
	"testing"

	"queclink-go/base.device.service/comm"
	"queclink-go/base.device.service/config"
	"queclink-go/base.device.service/core"
	"queclink-go/base.device.service/report"
)

func TestDeviceManagerEmbedding(t *testing.T) {
	config.Initialize("..", "/credentials.json")
	instance := &ExtendedDM{}
	instance.DeviceManager.InitializeDeviceCallback = instance.InitializeDeviceCallback
	instance.DeviceManager.InitializeUDPDeviceCallback = instance.InitializeUDPDeviceCallback

	core.InitializeDeviceManager(instance)

	d, err := core.InstanceDM.InitializeUDPDevice(&comm.UDPChannel{}, &report.Message{})
	if err != nil {
		t.Error("[TestDeviceManagerEmbedding]Error in device manager inheritance:", err.Error())
	}
	if d.GetIdentity() != "genx_123456789012" {
		t.Error("[TestDeviceManagerEmbedding]Error in device manager inheritance: invalid device identity")
	}
}

type ExtendedDM struct {
	core.DeviceManager
}

func (manager *ExtendedDM) InitializeUDPDevice(c comm.IChannel, message report.IMessage) (core.IDevice, error) {
	return manager.InitializeDevice(c, message)
}

func (manager *ExtendedDM) InitializeDevice(c comm.IChannel, message report.IMessage) (core.IDevice, error) {
	return &core.Device{Identity: "genx_123456789012"}, nil
}
