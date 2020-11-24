package models

import (
	"testing"

	"queclink-go/base.device.service/config"
)

func TestDeviceConfigurationLoad(t *testing.T) {
	err := config.Initialize("..", "/credentials.json")
	err = InitializeConnections(config.Config.GetBase().MysqDeviceMasterConnectionString)
	if err != nil {
		t.Error("Error establish connection:", err.Error())
	}
}

func TestInvalidConfigurationLoad(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {

		}
	}()
	err := config.Initialize("..", "/credentials.json")
	config.Config.GetBase().MysqDeviceMasterConnectionString = ""
	err = InitializeConnections(config.Config.GetBase().MysqDeviceMasterConnectionString)
	if err != nil {
		t.Error("Error load configuration:", err.Error())
	}
}
