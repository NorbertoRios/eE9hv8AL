package config

import (
	"testing"
)

func TestDeviceConfigurationLoad(t *testing.T) {
	err := Initialize("..", "/credentials.json")
	if err != nil {
		t.Error("Error load configuration:", err.Error())
	}
	if Config.GetBase().DeviceFacadeHost == "" {
		t.Error("Empty configuration value DeviceFacadeHost")
	}
}

func TestInvalidConfigurationLoad(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {

		}
	}()
	err := Initialize("..", "/credentials.json")
	if err != nil {
		t.Error("Error load configuration:", err.Error())
	}
}

type ExtendedConfiguration struct {
	Configuration
	SomeSetting string
}

func TestAccessViaParent(t *testing.T) {
	err := Initialize("..", "/credentials.json")
	if err != nil {
		t.Error("Error load configuration:", err.Error())
	}
	if Config.GetBase().DeviceFacadeHost == "" {
		t.Error("Empty configuration value DeviceFacadeHost")
	}
}
