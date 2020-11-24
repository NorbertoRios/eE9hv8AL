package models

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"queclink-go/base.device.service/config"

	"queclink-go/base.device.service/report"
)

//VehicleSupport feature based on ECU
type VehicleSupport struct {
	ID               int       `gorm:"column:vehId"`
	Identity         string    `gorm:"column:identity"`
	Vin              string    `gorm:"column:vehVin"`
	UpdateTime       time.Time `gorm:"column:vehUpdateTime"`
	ObdDescription   string    `gorm:"column:vehObdDescription"`
	MpgType          int       `gorm:"column:vehMpgType"`
	MpgCoefficient   string    `gorm:"column:vehMpgCoefficient"`
	VinSupport       bool      `gorm:"column:vehVinSupport"`
	FuelLevelSupport bool      `gorm:"column:vehFuelLevelSupport"`
	FuelTypeSupport  bool      `gorm:"column:vehFuelTypeSupport"`
	MafSupport       bool      `gorm:"column:vehMafSupport"`
	EfrSupport       bool      `gorm:"column:vehEfrSupport"`
	DtcSupport       bool      `gorm:"column:vehDtcSupport"`

	Uploaded bool `gorm:"-" sql:"-"`
}

//TableName for DeviceActivity model
func (VehicleSupport) TableName() string {
	return "ats.tblXirgoVehicleInformation"
}

//Save device state cache to database
func (vehicle *VehicleSupport) Save() {

	v := &VehicleSupport{
		Identity: vehicle.Identity,
		Vin:      vehicle.Vin}
	if rawdb.Model(&vehicle).Where(v).Updates(map[string]interface{}{
		"vehVin":              vehicle.Vin,
		"vehUpdateTime":       vehicle.UpdateTime,
		"vehObdDescription":   vehicle.ObdDescription,
		"vehMpgType":          vehicle.MpgType,
		"vehMpgCoefficient":   vehicle.MpgCoefficient,
		"vehVinSupport":       vehicle.VinSupport,
		"vehFuelLevelSupport": vehicle.FuelLevelSupport,
		"vehFuelTypeSupport":  vehicle.FuelTypeSupport,
		"vehMafSupport":       vehicle.MafSupport,
		"vehEfrSupport":       vehicle.EfrSupport,
		"vehDtcSupport":       vehicle.DtcSupport}).RowsAffected == 0 {
		rawdb.Save(&vehicle)
	}
	vehicle.Uploaded = true
}

//Delete vehicle information by ID
func (vehicle *VehicleSupport) Delete() error {
	v := &VehicleSupport{
		Identity: vehicle.Identity,
		Vin:      vehicle.Vin}
	return rawdb.Model(&vehicle).Where(v).Delete(v).Error
}

//EnableDTCSupport update dtc support for vehicle
func (vehicle *VehicleSupport) EnableDTCSupport() {
	if vehicle.DtcSupport {
		return
	}
	vehicle.DtcSupport = true
	if vehicle.Uploaded {
		log.Println("[VehicleSupport] Updated Vehicle DTC Support.  Identity:", vehicle.Identity)
		vehicle.Update(true)
	}
}

//Update current state for vehicle information
func (vehicle *VehicleSupport) Update(onlyDtc bool) {
	vehicle.Save()
	if !onlyDtc {
		vehicle.SendVinDecodeRequest()
	}
}

//UpdateState chenge state using message
//Used message fields:
//@VINSupport
//@FuelLevelSupport
//@MPGType
//@FuelLevelSupport
//@FuelTypeSupport
//@MAFSupport
//@EFRSupport
//@VIN
func (vehicle *VehicleSupport) UpdateState(identity string, message report.IMessage) {
	vehicle.Identity = identity
	vehicle.Vin = message.GetStringValue("VIN", "")
	vehicle.UpdateTime = time.Now().UTC()
	vehicle.ObdDescription = message.GetStringValue("OBDProtocolDescription", "")
	mpgType := message.GetStringValue("MPGType", "")
	impgType, err := strconv.ParseInt(mpgType, 10, 32)
	if err != nil {
		impgType = 0
	}
	vehicle.MpgType = int(impgType)
	vehicle.MpgCoefficient = message.GetStringValue("MPGCoefficient", "")
	if message.GetStringValue("VINSupport", "") == "1" {
		vehicle.VinSupport = true
	} else {
		vehicle.VinSupport = false
	}

	if message.GetStringValue("FuelLevelSupport", "") == "1" {
		vehicle.FuelLevelSupport = true
	} else {
		vehicle.FuelLevelSupport = false
	}

	if message.GetStringValue("FuelTypeSupport", "") == "1" {
		vehicle.FuelTypeSupport = true
	} else {
		vehicle.FuelTypeSupport = false
	}

	if message.GetStringValue("MAFSupport", "") == "1" {
		vehicle.MafSupport = true
	} else {
		vehicle.MafSupport = false
	}
	if message.GetStringValue("EFRSupport", "") == "1" {
		vehicle.EfrSupport = true
	} else {
		vehicle.EfrSupport = false
	}
}

//Sync vehicle information
func (vehicle *VehicleSupport) Sync(identity, vin string) {
	fromBase, error := FindVehicleSupport(identity, vin)
	vehicle.Identity = identity
	vehicle.Vin = vin
	if error != nil {
		vehicle.Save()
		return
	}

	vehicle.ID = fromBase.ID
	var syncRes = vehicle.Synchronize(fromBase)

	switch syncRes {
	case 1: //DifferentId
	case 2: //DifferentVin
	case 4: //NotEquivalent
		{
			vehicle.Save()
			break
		}
	}
}

//Synchronize vehicle information
func (vehicle *VehicleSupport) Synchronize(vs *VehicleSupport) int {
	if vehicle.Identity != vs.Identity {
		return 1 //DifferentId
	}
	if vehicle.Vin != vs.Vin {
		return 2 //DifferentVin
	}
	if vehicle.IsEqual(vs) {
		return 3 //Equivalent
	}
	return 4 //NotEquivalent
}

//IsEqual checks for equals to input vs
func (vehicle *VehicleSupport) IsEqual(vs *VehicleSupport) bool {
	return vehicle.ObdDescription == vs.ObdDescription &&
		vehicle.MpgType == vs.MpgType &&
		vehicle.MpgCoefficient == vs.MpgCoefficient &&
		vehicle.VinSupport == vs.VinSupport &&
		vehicle.FuelLevelSupport == vs.FuelLevelSupport &&
		vehicle.FuelTypeSupport == vs.FuelTypeSupport &&
		vehicle.MafSupport == vs.MafSupport &&
		vehicle.EfrSupport == vs.EfrSupport
}

//SendVinDecodeRequest Sends request to facade service to decode VIN code
func (vehicle *VehicleSupport) SendVinDecodeRequest() (*http.Response, error) {
	log.Println("[VehicleSupport] Try send vin decode request.  VIN:", vehicle.Vin)
	if vehicle.Vin == "" {
		return nil, nil
	}
	log.Println("[VehicleSupport] Sending vin decode request:", vehicle.Vin, " to:", config.Config.GetBase().DeviceFacadeHost)
	resp, err := http.PostForm(fmt.Sprintf("%v/vin/decode", config.Config.GetBase().DeviceFacadeHost), url.Values{"VIN": {vehicle.Vin}})
	if err != nil || resp.StatusCode != 200 {
		log.Fatal("[VehicleSupport] Invalid response for:", vehicle.Vin, " to:", config.Config.GetBase().DeviceFacadeHost)
	}
	return resp, err
}

//FindVehicleSupport lookup vehicle information by identity and vin
func FindVehicleSupport(identity, vin string) (*VehicleSupport, error) {

	v := &VehicleSupport{}
	err := rawdb.Model(v).Select("vehId, identity, vehVin, vehUpdateTime, vehObdDescription, vehMpgType, vehMpgCoefficient," +
		"vehVinSupport, vehFuelLevelSupport, vehFuelTypeSupport,  vehMafSupport, vehEfrSupport, vehDtcSupport").
		Where(&VehicleSupport{Identity: identity, Vin: vin}).Scan(v).Error
	v.Uploaded = false
	return v, err
}
