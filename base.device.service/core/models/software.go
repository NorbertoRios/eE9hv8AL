package models

import (
	"encoding/json"
	"log"
	"time"

	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
)

//Software struct
type Software struct {
	FirmwareVersion    string
	SoftwareVersion    string
	OBDFirmwareVersion string
	UpdatedAt          *utils.JSONTime
}

//Marshal software struct
func (data *Software) Marshal() string {
	jMessage, jerr := json.Marshal(data)
	if jerr != nil {
		log.Println("Marshal software data error:", jerr)
		return ""
	}
	return string(jMessage)
}

//UpdateDeviceSoftware handle software version
func (data *Software) UpdateDeviceSoftware(message report.IMessage) bool {
	if _, f := message.GetValue("MainFirmwareVersion"); !f {
		return false
	}
	if _, f := message.GetValue("OBDFirmwareVersion"); !f {
		return false
	}

	software := &Software{
		FirmwareVersion:    message.GetStringValue("MainFirmwareVersion", ""),
		OBDFirmwareVersion: message.GetStringValue("OBDFirmwareVersion", ""),
	}
	return !data.compare(software)
}

func (data *Software) compare(software *Software) bool {
	if data.FirmwareVersion != software.FirmwareVersion ||
		data.OBDFirmwareVersion != software.OBDFirmwareVersion ||
		(data.UpdatedAt != nil && time.Now().UTC().Sub(data.UpdatedAt.Time).Seconds() >= 86400) {
		software.UpdatedAt = &utils.JSONTime{Time: time.Now().UTC()}
		data.updateSoftware(software)
		return false
	}
	return true
}

func (data *Software) updateSoftware(software *Software) {
	data.FirmwareVersion = software.FirmwareVersion
	data.OBDFirmwareVersion = software.FirmwareVersion
	data.UpdatedAt = software.UpdatedAt
}

//UnMarshalSoftware given string to Software struct
func UnMarshalSoftware(str string) (*Software, error) {
	software := &Software{}
	err := json.Unmarshal([]byte(str), software)
	if err != nil {
		return &Software{}, err
	}
	return software, nil
}
