package queclinkreport

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"queclink-go/base.device.service/utils"

	"queclink-go/queclinkreport/fields"
	"queclink-go/queclinkreport/gv55/resp55"
)

func TestHumanizeUniqueId(t *testing.T) {
	parser := &BaseLocationReportParser{}
	var message IQueclinkMessage = NewQueclinkMessage()
	message.SetValue(fields.UniqueID, []byte{86, 38, 30, 03, 18, 02, 90, 9})
	parser.HumanizeUniqueID(message, fields.UniqueID)
	if message.UniqueID() != "863830031802909" {
		t.Error(fmt.Sprintf("[TestHumanizeUniqueId]Failed; Should:863830031802909; Got:%v", message.UniqueID()))
	}
	devID, found := message.GetValue(fields.DevID)

	if !found || devID != "queclink_863830031802909" {
		t.Error(fmt.Sprintf("[TestHumanizeUniqueId]Failed; Should:queclink_863830031802909; Got:%v", devID))
	}
}

func TestHumanizeSpeed(t *testing.T) {
	parser := &BaseLocationReportParser{}
	var message IQueclinkMessage = NewQueclinkMessage()
	message.SetValue(fields.Speed, []byte{0x00, 0x31, 0x08})
	parser.HumanizeSpeed(message)

	speed, found := message.GetValue(fields.Speed)

	if !found || math.Round(float64(speed.(float32))*100)/100 != 13.83 {
		t.Error(fmt.Sprintf("[TestHumanizeSpeed]Failed; Should:13.833333; Got:%v", math.Round(float64(speed.(float32)*100)/100)))
	}
}

func TestHumanizeDateTime(t *testing.T) {
	parser := &BaseLocationReportParser{}
	var message IQueclinkMessage = NewQueclinkMessage()
	message.SetValue(fields.GPSUTCTime, []byte{0x07, 0xE2, 0x0C, 0x12, 0x0D, 0x23, 0x33})
	parser.HumanizeDateTime(message, fields.GPSUTCTime)

	ts, _ := message.GetValue(fields.GPSUTCTime)

	timestamp, _ := ts.(*utils.JSONTime)

	value, _ := timestamp.MarshalJSON()

	should := "\"2018-12-18T13:35:51Z\""
	v := string(value)
	if v != should {
		t.Error("[TestHumanizeDateTime]Invalid GPSUTCTime value. Should:", should, "; Current value:", v)
	}
	jsonMessage, _ := json.Marshal(message)
	if !strings.Contains(string(jsonMessage), "\"GPSUTCTime\":\"2018-12-18T13:35:51Z\"") {
		t.Error("[TestHumanizeDateTime]Invalid serialize message result:", string(jsonMessage))
	}

}

func TestHumanizeCurrentMileage(t *testing.T) {
	parser := &BaseLocationReportParser{}
	var message IQueclinkMessage = NewQueclinkMessage()
	message.SetValue(fields.CurrentMileage, []byte{0x1B, 0x45, 0x04})
	parser.HumanizeCurrentMileage(message)

	cm, found := message.GetValue(fields.CurrentMileage)

	if !found || math.Round(float64(cm.(float32))*100)/100 != 6981.4 {
		t.Error(fmt.Sprintf("[TestHumanizeCurrentMileage]Failed; Should:6981.4; Got:%v", math.Round(float64(cm.(float32)*100)/100)))
	}
}

func TestHumanizeTotalMileage(t *testing.T) {
	parser := &BaseLocationReportParser{}
	var message IQueclinkMessage = NewQueclinkMessage()
	message.SetValue(fields.Odometer, []byte{0x00, 0x00, 0x5D, 0x51, 0x00})
	parser.HumanizeTotalMileage(message)

	cm, found := message.GetValue(fields.Odometer)

	if !found || math.Round(float64(cm.(int32))*10)/10 != 23889000.0 {
		t.Error(fmt.Sprintf("[TestHumanizeTotalMileage]Failed; Should:23889000.0; Got:%v", math.Round(float64(cm.(float32)*10)/10)))
	}
}

func TestHumanizeCurrentHourMeterCount(t *testing.T) {
	parser := &BaseLocationReportParser{}
	var message IQueclinkMessage = NewQueclinkMessage()
	message.SetValue(fields.CurrentHourMeterCount, []byte{10, 15, 23})
	parser.HumanizeCurrentHourMeterCount(message)

	cm, found := message.GetValue(fields.CurrentHourMeterCount)

	if !found || cm.(int32) != 36923 {
		t.Error(fmt.Sprintf("[TestHumanizeCurrentHourMeterCount]Failed; Should:36923; Got:%v", cm.(int32)))
	}
}

func TestHumanizeTotalHourMeterCount(t *testing.T) {
	parser := &BaseLocationReportParser{}
	var message IQueclinkMessage = NewQueclinkMessage()
	message.SetValue(fields.TotalHourMeterCount, []byte{0x10, 0x11, 0x12, 0x13, 0x15, 0x23})
	parser.HumanizeTotalHourMeterCount(message)

	cm, found := message.GetValue(fields.TotalHourMeterCount)

	if !found || cm.(int64) != 970395103295 {
		t.Error(fmt.Sprintf("[TestHumanizeCurrentHourMeterCount]Failed; Should:970395103295; Got:%v", cm.(int64)))
	}
}

func TestHumanizeTimeStamp(t *testing.T) {
	parser := &BaseLocationReportParser{}
	var message IQueclinkMessage = NewQueclinkMessage()
	message.SetValue(fields.GPSUTCTime, []byte{0x07, 0xE2, 0x0C, 0x12, 0x0D, 0x23, 0x33})
	message.SetValue(fields.SendTime, []byte{0x07, 0xE2, 0x0C, 0x12, 0x0D, 0x23, 0x35})
	parser.HumanizeDateTime(message, fields.GPSUTCTime)
	parser.HumanizeDateTime(message, fields.SendTime)
	parser.HumanizeTimestamp(message)

	iTs, _ := message.GetValue(fields.TimeStamp)

	timestamp, _ := iTs.(*utils.JSONTime)

	value, _ := timestamp.MarshalJSON()

	should := "\"2018-12-18T13:35:53Z\""
	v := string(value)
	if v != should {
		t.Error("[TestHumanizeDateTime]Invalid GPSUTCTime value. Should:", should, "; Current value:", v)
	}

	message.SetValue(fields.MessageHeader, "+BSP")
	parser.HumanizeTimestamp(message)

	iTs, _ = message.GetValue(fields.TimeStamp)
	timestamp, _ = iTs.(*utils.JSONTime)

	value, _ = timestamp.MarshalJSON()
	should = "\"2018-12-18T13:35:51Z\""
	v = string(value)
	if v != should {
		t.Error("[TestHumanizeDateTime]Invalid GPSUTCTime value. Should:", should, "; Current value:", v)
	}
}

func TestHumanizeLatLon(t *testing.T) {
	parser := &BaseLocationReportParser{}
	var message IQueclinkMessage = NewQueclinkMessage()
	message.SetValue(fields.Latitude, []byte{0x01, 0x7A, 0x65, 0xAD})
	message.SetValue(fields.Longitude, []byte{0xF9, 0x99, 0x4F, 0xB8})
	parser.HumanizeLngLat(message, fields.Latitude)
	parser.HumanizeLngLat(message, fields.Longitude)

	l, found := message.GetValue(fields.Latitude)

	rl := math.Round(float64(l.(float32)) * 1000000.0)

	if !found || rl/1000000.0 != 24.798637 {
		t.Error(fmt.Sprintf("[TestHumanizeLatLon]Failed; Should:24.798637; Got:%v", math.Round(float64(l.(float32)*1000000.0)/1000000.0)))
	}

	l, found = message.GetValue(fields.Longitude)

	rl = math.Round(float64(l.(float32)) * 1000000.0)

	if !found || rl/1000000.0 != -107.393097 {
		t.Error(fmt.Sprintf("[TestHumanizeLatLon]Failed; Should:-107.393097; Got:%v", math.Round((float64(l.(float32))*1000000.0)/1000000.0)))
	}
}

func TestHumanizeInputs(t *testing.T) {
	parser := &BaseLocationReportParser{}
	var message IQueclinkMessage = NewQueclinkMessage()
	message.SetValue(fields.IgnitionStatus, byte(16))
	message.SetValue(fields.DigitalInputStatus, byte(3))

	parser.HumanizeInputs(message)

	gpio, found := message.GetValue(fields.GPIO)

	if !found || gpio.(byte) != 1 {
		t.Error(fmt.Sprintf("[TestHumanizeInputs]Failed; GPIO Should:1; Got:%v", gpio))
	}

	ignition, found := message.GetValue(fields.IgnitionState)

	if !found || ignition.(byte) != 1 {
		t.Error(fmt.Sprintf("[TestHumanizeInputs]Failed; IgnitionStatus Should:1; Got:%v", ignition))
	}
}

func TestHumanizeInputsIOff(t *testing.T) {
	parser := &BaseLocationReportParser{}
	var message IQueclinkMessage = NewQueclinkMessage()
	message.SetValue(fields.IgnitionStatus, byte(16))
	message.SetValue(fields.DigitalInputStatus, byte(2))

	parser.HumanizeInputs(message)

	gpio, found := message.GetValue(fields.GPIO)

	if !found || gpio.(byte) != 1 {
		t.Error(fmt.Sprintf("[TestHumanizeInputsIOff]Failed; GPIO Should:1; Got:%v", gpio))
	}

	ignition, found := message.GetValue(fields.IgnitionState)

	if !found || ignition.(byte) != 0 {
		t.Error(fmt.Sprintf("[TestHumanizeInputsIOff]Failed; IgnitionStatus Should:1; Got:%v", ignition))
	}
}

func TestHumanizeInputsIOn(t *testing.T) {
	parser := &BaseLocationReportParser{}
	var message IQueclinkMessage = NewQueclinkMessage()
	message.SetValue(fields.IgnitionStatus, byte(32))
	message.SetValue(fields.DigitalInputStatus, byte(2))

	parser.HumanizeInputs(message)

	gpio, found := message.GetValue(fields.GPIO)

	if !found || gpio.(byte) != 1 {
		t.Error(fmt.Sprintf("[TestHumanizeInputsIOff]Failed; GPIO Should:1; Got:%v", gpio))
	}

	ignition, found := message.GetValue(fields.IgnitionState)

	if !found || ignition.(byte) != 1 {
		t.Error(fmt.Sprintf("[TestHumanizeInputsIOff]Failed; IgnitionStatus Should:1; Got:%v", ignition))
	}
}

func TestValidity(t *testing.T) {
	config := &ReportConfiguration{}
	LoadReportConfiguration("ReportConfiguration.xml", config)
	parser := &BaseLocationReportParser{}
	var message IQueclinkMessage = NewQueclinkMessage()

	message.SetValue(fields.GPSAccuracy, byte(2))
	message.SetValue(fields.MessageType, int32(7))
	message.SetValue(fields.TimeStamp, &utils.JSONTime{Time: time.Now().UTC()})
	message.SetValue(fields.MessageHeader, "+RSP")
	message.SetValue(fields.Satellites, byte(7))
	message.SetValue(fields.Latitude, float64(24.123456))
	message.SetValue(fields.Longitude, float64(-107.123456))

	reportConfig, _ := config.Find(16, "+RSP", resp55.GTFRI)
	parser.HumanizeValidity(reportConfig, message)

	ivalidity, found := message.GetValue(fields.GpsValidity)

	if !found || ivalidity.(byte) != 1 {
		t.Error(fmt.Sprintf("[TestValidity]Failed; Validity Should:1; Got:%v", ivalidity.(byte)))
	}
}

func TestLastKnownLocationValidity(t *testing.T) {
	config := &ReportConfiguration{}
	LoadReportConfiguration("ReportConfiguration.xml", config)
	parser := &BaseLocationReportParser{}
	var message IQueclinkMessage = NewQueclinkMessage()

	message.SetValue(fields.GPSAccuracy, byte(0))
	message.SetValue(fields.MessageType, int32(21))
	message.SetValue(fields.TimeStamp, &utils.JSONTime{Time: time.Now().UTC()})
	message.SetValue(fields.MessageHeader, "+EVT")
	message.SetValue(fields.Satellites, byte(7))
	message.SetValue(fields.Latitude, float32(0))
	message.SetValue(fields.Longitude, float32(0))

	reportConfig, _ := config.Find(16, "+EVT", 21)
	parser.HumanizeValidity(reportConfig, message)

	ivalidity, found := message.GetValue(fields.GpsValidity)

	if !found || ivalidity.(byte) != 0 {
		t.Error(fmt.Sprintf("[TestValidity]Failed; Validity Should:0; Got:%v", ivalidity.(byte)))
	}
}

func TestPacketWithAdditionalBytes(t *testing.T) {
	packet := []byte{0x2B, 0x42, 0x53, 0x50, 0x07, 0x00, 0x3E, 0x1F, 0xBF, 0x00, 0x5B, 0x50, 0x01, 0x00, 0x02, 0x07, 0x56, 0x30, 0x02, 0x03, 0x01, 0x58, 0x1D, 0x00, 0x64, 0x35, 0xEB, 0x01, 0x00, 0x22, 0x0B, 0x00, 0x00, 0x00, 0xEE, 0x01, 0x00, 0x00, 0x0D, 0x08, 0x01, 0x0E, 0x06, 0xEE, 0xFA, 0x15, 0xF6, 0x6E, 0x01, 0x21, 0x73, 0xB4, 0x07, 0xE3, 0x04, 0x0C, 0x14, 0x01, 0x02, 0x03, 0x34, 0x00, 0x20, 0x52, 0x08, 0x07, 0xD0, 0xB3, 0xE7, 0x00, 0x00, 0x0E, 0x00, 0x00, 0x00, 0x06, 0x0B, 0x09, 0x07, 0xE3, 0x04, 0x0C, 0x14, 0x01, 0x02, 0x5D, 0xBE, 0x75, 0x69, 0x0D, 0x0A}
	config := &ReportConfiguration{}
	LoadReportConfiguration("ReportConfiguration.xml", config)
	parser := Parser{}
	messages, _ := parser.Parse(packet)
	ack, found := messages[0].GetValue("Ack")

	sack := string(ack.([]byte))
	if !found || sack != "+SACK:5DBE$" {
		t.Error(fmt.Sprintf("[TestPacketWithAdditionalBytes]Failed; ACK Should:+SACK:5DBE$; Got:%v", sack))
	}

}

func TestPacketNilError(t *testing.T) {
	packet := []byte{0x2B, 0x42, 0x56, 0x54, 0x0D, 0x03, 0xFE, 0x1F, 0xFF, 0x00, 0xBD, 0x41, 0x05, 0x02, 0x05, 0x0B, 0x67, 0x76, 0x35, 0x35, 0x00, 0x00, 0x00, 0x00, 0x0C, 0x00, 0x00, 0x01, 0x00, 0x21, 0x0A, 0x00, 0x00, 0xB6, 0xCB, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xCD, 0xFB, 0x43, 0x3E, 0xE9, 0x02, 0x9B, 0xC4, 0x61, 0x07, 0xE3, 0x05, 0x02, 0x14, 0x37, 0x0E, 0x03, 0x02, 0x07, 0x20, 0xEA, 0xB0, 0x00, 0x22, 0x4A, 0x8D, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07, 0xE3, 0x05, 0x03, 0x09, 0x2B, 0x11, 0x00, 0xF2, 0x14, 0x36, 0x0D, 0x0A}
	config := &ReportConfiguration{}
	LoadReportConfiguration("ReportConfiguration.xml", config)
	parser := Parser{}
	messages, _ := parser.Parse(packet)
	ack, found := messages[0].GetValue("Ack")

	sack := string(ack.([]byte))
	if !found || sack != "+SACK:00F2$" {
		t.Error(fmt.Sprintf("[TestPacketWithAdditionalBytes]Failed; ACK Should:+SACK:5DBE$; Got:%v", sack))
	}

}

/*
func TestGV600LTEFRIParsing(t *testing.T) {
	packet := []byte{0x2B, 0x42, 0x53, 0x50, 0x07, 0x00, 0x9E, 0xBF, 0x00, 0xED, 0x01, 0x05, 0x05, 0x0C, 0x56, 0x40, 0x19, 0x03, 0x4F, 0x13, 0x06, 0x07, 0x0E, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0x08, 0x10, 0x01, 0x01, 0xF8, 0xAE, 0xC6, 0xF5, 0x02, 0xEF, 0x5D, 0x7A, 0x07, 0xE3, 0x07, 0x1A, 0x13, 0x11, 0x1F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x94, 0x5A, 0x0D, 0x0A}
	config := &ReportConfiguration{}
	LoadReportConfiguration("ReportConfiguration.xml", config)
	parser := Parser{}
	messages, _ := parser.Parse(packet)
	if len(messages) == 0 {
		t.Error("[TestGV600LTEFRIParsing]Messages count is 0")
	}
}
*/
