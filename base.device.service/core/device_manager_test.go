package core

import (
	"testing"

	"queclink-go/base.device.service/comm"
	"queclink-go/base.device.service/config"
	"queclink-go/base.device.service/report"
)

func TestDeviceManagerEmbedding(t *testing.T) {
	config.Initialize("../config/credentials.json", &config.Configuration{})
	instance := &ExtendedDM{}
	instance.DeviceManager.InitializeDevice = instance.InitializeDevice
	instance.DeviceManager.InitializeUDPDevice = instance.InitializeUDPDevice
	InitializeDeviceManager(instance)
	d, err := InstanceDM.InitializeUDPDevice(&comm.UDPChannel{}, &report.Message{})
	if err != nil {
		t.Error("[TestDeviceManagerEmbedding]Error in device manager inheritance:", err.Error())
	}
	if d.GetIdentity() != "genx_123456789012" {
		t.Error("[TestDeviceManagerEmbedding]Error in device manager inheritance: invalid device identity")
	}
}

type ExtendedDM struct {
	DeviceManager
}

func (manager *ExtendedDM) InitializeUDPDevice(c comm.IChannel, message report.IMessage) (IDevice, error) {
	return manager.InitializeDevice(c, message)
}

func (manager *ExtendedDM) InitializeDevice(c comm.IChannel, message report.IMessage) (IDevice, error) {
	return &Device{Identity: "genx_123456789012"}, nil
}
