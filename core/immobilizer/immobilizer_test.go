package immobilizer

import (
	"testing"

	"queclink-go/base.device.service/config"
	"queclink-go/qconfig"
)

func TestGetFacadeOutputsState(t *testing.T) {
	config.Initialize("../../service/credentials.json", &qconfig.QConfiguration{})
	_, err := GetFacadeOutputsState("queclink_868034001645340")
	if err != nil {
		t.Error("To get outputs state")
	}
}

func TestGetImmobilizerCommand(t *testing.T) {
	config.Initialize("../../service/credentials.json", &qconfig.QConfiguration{})
	command, _ := GetImmobilizerCommand("queclink_868034001645340", "OUT0", "ARMED", "HIGH", 0)
	should := "AT+GTOUT=gv55,1,0,0,0,0,0,,,,0,,,,,,,FFFF$"
	if command != should {
		t.Errorf("Invalid value of output immobilizer command. Should:%v; Current:%v", should, command)
	}
}

func TestNewDeviceImmobilizer(t *testing.T) {
	config.Initialize("../../service/credentials.json", &qconfig.QConfiguration{})
	command, _ := GetImmobilizerCommand("queclink_000000000000001", "OUT0", "ARMED", "HIGH", 0)
	should := "AT+GTOUT=gv55,1,0,0,0,0,0,,,,0,,,,,,,FFFF$"
	if command != should {
		t.Errorf("Invalid value of output immobilizer command. Should:%v; Current:%v", should, command)
	}
}
