package models

import (
	"testing"
	"time"

	"queclink-go/base.device.service/config"
)

func TestVehicleInformation(t *testing.T) {
	err := config.Initialize("..", "/credentials.json")
	InitializeConnections(config.Config.GetBase().MysqDeviceMasterConnectionString)
	v := &VehicleSupport{
		Identity:         "xirgo_123456789012345",
		Vin:              "WASDF3234343GHR5",
		DtcSupport:       true,
		EfrSupport:       true,
		FuelLevelSupport: true,
		FuelTypeSupport:  true,
		MafSupport:       true,
		VinSupport:       true,
		MpgCoefficient:   "1.00",
		MpgType:          15,
		ObdDescription:   "0080",
		UpdateTime:       time.Now().UTC(),
	}
	err = v.Delete()
	if err != nil {
		t.Error("Unable to delete existing vehicle information")
	}
	v.Save()
	v1, err := FindVehicleSupport("xirgo_123456789012345", "WASDF3234343GHR5")

	if err != nil {
		t.Error("Unable to find vehicle information")
	}
	if v1.Identity != v.Identity || v1.DtcSupport != v.DtcSupport || v1.Vin != v.Vin || v1.EfrSupport != v.EfrSupport ||
		v1.FuelLevelSupport != v.FuelLevelSupport || v1.FuelTypeSupport != v.FuelTypeSupport || v1.MafSupport != v.MafSupport ||
		v1.VinSupport != v.VinSupport || v1.MpgCoefficient != v.MpgCoefficient || v1.MpgType != v.MpgType || v1.ObdDescription != v.ObdDescription {

		t.Error("Invalid vehicle information value")
	}
}

func TestFindNonExistingVehicleInformation(t *testing.T) {
	err := config.Initialize("../..", "/credentials.json")
	InitializeConnections(config.Config.GetBase().MysqDeviceMasterConnectionString)
	v := &VehicleSupport{
		Identity: "xirgo_123456789012345",
		Vin:      "WASDF3234343GHR5",
	}
	err = v.Delete()
	if err != nil {
		t.Error("Unable to delete existing vehicle information")
	}
	v1, err := FindVehicleSupport("xirgo_123456789012345", "WASDF3234343GHR5")

	if v1 == nil {
		t.Error("Vehicle information is null")
	}
}

func TestSendVinDecodeRequest(t *testing.T) {
	err := config.Initialize("..", "/credentials.json")
	InitializeConnections(config.Config.GetBase().MysqDeviceMasterConnectionString)
	v := &VehicleSupport{
		Identity: "xirgo_123456789012345",
		Vin:      "WMZYV5C32J3E04723",
	}
	resp, err := v.SendVinDecodeRequest()

	if err != nil {
		t.Error("Unable to delete existing vehicle information")
	}

	if resp.StatusCode != 200 {
		t.Error("Invalid VIN decode request")
	}
}

func TestVehicleInformationIsEqual(t *testing.T) {
	err := config.Initialize("..", "/credentials.json")
	InitializeConnections(config.Config.GetBase().MysqDeviceMasterConnectionString)
	v := &VehicleSupport{
		Identity:         "xirgo_123456789012345",
		Vin:              "WASDF3234343GHR5",
		DtcSupport:       true,
		EfrSupport:       true,
		FuelLevelSupport: true,
		FuelTypeSupport:  true,
		MafSupport:       true,
		VinSupport:       true,
		MpgCoefficient:   "1.00",
		MpgType:          15,
		ObdDescription:   "0080",
		UpdateTime:       time.Now().UTC(),
	}
	err = v.Delete()
	if err != nil {
		t.Error("Unable to delete existing vehicle information")
	}
	v.Save()
	v1, _ := FindVehicleSupport("xirgo_123456789012345", "WASDF3234343GHR5")
	if !v.IsEqual(v1) {
		t.Error("Invalid Equal result")
	}
}

func TestVehicleSyncNewRecord(t *testing.T) {
	err := config.Initialize("..", "/credentials.json")
	InitializeConnections(config.Config.GetBase().MysqDeviceMasterConnectionString)
	v := &VehicleSupport{
		Identity:         "xirgo_123456789012345",
		Vin:              "WASDF3234343GHR5",
		DtcSupport:       true,
		EfrSupport:       true,
		FuelLevelSupport: true,
		FuelTypeSupport:  true,
		MafSupport:       true,
		VinSupport:       true,
		MpgCoefficient:   "1.00",
		MpgType:          15,
		ObdDescription:   "0080",
		UpdateTime:       time.Now().UTC(),
	}
	err = v.Delete()
	if err != nil {
		t.Error("Unable to delete existing vehicle information")
	}
	v.Sync(v.Identity, v.Vin)
	v1, _ := FindVehicleSupport("xirgo_123456789012345", "WASDF3234343GHR5")
	if !v.IsEqual(v1) {
		t.Error("Invalid Equal result")
	}
}

func TestVehicleSyncNewVin(t *testing.T) {
	err := config.Initialize("..", "/credentials.json")
	InitializeConnections(config.Config.GetBase().MysqDeviceMasterConnectionString)
	v := &VehicleSupport{
		Identity:         "xirgo_123456789012345",
		Vin:              "WASDF3234343GHR5",
		DtcSupport:       true,
		EfrSupport:       true,
		FuelLevelSupport: true,
		FuelTypeSupport:  true,
		MafSupport:       true,
		VinSupport:       true,
		MpgCoefficient:   "1.00",
		MpgType:          15,
		ObdDescription:   "0080",
		UpdateTime:       time.Now().UTC(),
	}
	err = v.Delete()
	if err != nil {
		t.Error("Unable to delete existing vehicle information")
	}
	v.Sync(v.Identity, v.Vin)
	v1, _ := FindVehicleSupport("xirgo_123456789012345", "WASDF3234343GHR5")
	v1.Vin = "WASDF3234343GHR5"
	v.Sync(v.Identity, v.Vin)
	v2, _ := FindVehicleSupport("xirgo_123456789012345", "WASDF3234343GHR5")
	if !v2.IsEqual(v1) {
		t.Error("Invalid Equal result")
	}
}
