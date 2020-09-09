package core

import (
	"testing"
	"time"

	"queclink-go/base.device.service/core/models"
	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
)

func TestValidTimeStamp(t *testing.T) {
	device := &Device{}
	device.Activity = &models.DeviceActivity{}
	device.Activity.GPSTimeStamp = time.Now().UTC().Add(-5 * time.Minute)
	message := report.NewMessage()
	message.SetValue("GpsValidity", byte(1))
	nextTime := time.Now().UTC()
	message.SetValue("TimeStamp", &utils.JSONTime{Time: nextTime})
	device.UpdateCurrentGpsTimeStamp(message)
	mt, _ := message.GetValue("TimeStamp")
	mTimeStamp := mt.(*utils.JSONTime)
	if mTimeStamp.Time != device.Activity.GPSTimeStamp {
		t.Error("[TestValidTimeStamp]Invalid timestamp value. Should:", mTimeStamp.Time.String(), "; Current:", device.Activity.GPSTimeStamp.String())
	}
}

func TestTimeStampFromFeature(t *testing.T) {
	device := &Device{}
	device.Activity = &models.DeviceActivity{}
	lastKnownTimeStimp := time.Now().UTC().Add(-5 * time.Minute)
	device.Activity.GPSTimeStamp = lastKnownTimeStimp

	message := report.NewMessage()
	message.SetValue("GpsValidity", byte(1))
	nextTime := time.Now().UTC().Add(11 * time.Minute)
	message.SetValue("TimeStamp", &utils.JSONTime{Time: nextTime})

	device.UpdateCurrentGpsTimeStamp(message)
	mt, _ := message.GetValue("TimeStamp")
	mTimeStamp := mt.(*utils.JSONTime)
	if mTimeStamp.Time != lastKnownTimeStimp {
		t.Error("[TestValidTimeStamp]Invalid timestamp value. Should:", mTimeStamp.Time.String(), "; Current:", lastKnownTimeStimp.String())
	}
}

func TestTimeStampBC(t *testing.T) {
	device := &Device{}
	device.Activity = &models.DeviceActivity{}
	lastKnownTimeStimp := time.Now().UTC().Add(-5 * time.Minute)
	device.Activity.GPSTimeStamp = lastKnownTimeStimp

	message := report.NewMessage()
	message.SetValue("GpsValidity", byte(1))
	nextTime := time.Now().UTC().AddDate(-6, 0, 0)
	message.SetValue("TimeStamp", &utils.JSONTime{Time: nextTime})

	device.UpdateCurrentGpsTimeStamp(message)

	if device.Activity.GPSTimeStamp != lastKnownTimeStimp {
		t.Error("[TestValidTimeStamp]Invalid timestamp value. Should:", device.Activity.GPSTimeStamp.String(), "; Current:", lastKnownTimeStimp.String())
	}
}

/*
func TestDeviceInitialization(t *testing.T) {
	config.Initialize("../config/credentials.json")
	models.InitializeConnections()
	rabbit.InitializeRabbitConnection()
	device := &Device{}
	device.Initialize("862522030140750")
	if device.Activity == nil {
		t.Error("[TestDeviceInitialization]Activity is nil")
	}
	if device.VehicleSupport == nil {
		t.Error("[TestDeviceInitialization]Activity is nil")
	}
	if device.Configuration == nil {
		t.Error("[TestDeviceInitialization]Configuration is nil")
	}
}

func TestLocationMessageProcessing(t *testing.T) {
	config.Initialize("../config/credentials.json", &config.Configuration{})
	report.LoadReportConfiguration("ReportConfiguration.xml", &report.ReportConfiguration{})

	report.LoadReportConfiguration("../report/ReportConfiguration.xml")

	rabbit.InitializeRabbitConnection()
	models.InitializeConnections()
	device := &Device{}
	device.Initialize("355922063699116")
	packet := "$$355922063699116,6001,2018/11/01,14:01:24,57.31913,-111.50500,298,10.5,0.6,0.0,1247,167,18,1.2,1495,4.3,14.6,12,1,0,91.7,-156,4,-1031,158##"
	message, _ := report.Parse(packet)

	err := device.ProcessMessage(message)
	if err != nil {
		t.Error("[TestLocationMessageProcessing]Handle message error:", err.Error())
	}
}

func TestInfoMessageProcessing(t *testing.T) {
	config.Initialize("../config/credentials.json")
	report.LoadReportConfiguration("../report/ReportConfiguration.xml")
	rabbit.InitializeRabbitConnection()
	models.InitializeConnections()
	device := &Device{}
	device.Initialize("355922061549941")
	packet := "$$355922061549941,7002,0080,2C4RDGBG7GR217706,60386,0,12.5,0.0,16,1.00,1,1,1,0,0,002##"
	message, _ := report.Parse(packet)

	err := device.ProcessMessage(message)
	if err != nil {
		t.Error("[TestInfoMessageProcessing]Handle message error:", err.Error())
	}
}

func TestAckMessageProcessing(t *testing.T) {
	config.Initialize("../config/credentials.json")
	report.LoadReportConfiguration("../report/ReportConfiguration.xml")
	rabbit.InitializeRabbitConnection()
	models.InitializeConnections()

	cfg := &models.DeviceConfig{
		Identity:  "xirgo_355922061549942",
		Command:   "+XT:1010,4040,216.187.77.151,,,www.kyivstar.net,,6,0,0,255,240,1\r\n+XT:355922061549942,3001,1,1,0,0,0",
		SentAt:    utils.NullTime{Time: time.Now().UTC(), Valid: false},
		DevID:     1,
		CreatedAt: time.Now().UTC(),
	}

	cfg.DeleteAll()
	cfg.Save()

	_, found := models.FindDeviceConfigByIdentity("xirgo_355922061549942")
	if !found {
		t.Error("[TestAckMessageProcessing]Config doesn't exists:")
	}

	device := &Device{}
	device.Initialize("355922061549942")
	packet := "$$355922061549942,3001,1.0,1,0,0,0##"
	message, _ := report.Parse(packet)

	err := device.ProcessMessage(message)
	if err != nil {
		t.Error("[TestInfoMessageProcessing]Handle message error:", err.Error())
	}
	_, found = models.FindDeviceConfigByIdentity(device.Identity)
	if found {
		t.Error("[TestAckMessageProcessing]Handle Ack error. Config exists:")
	}
}

func TestDtcMessageProcessing(t *testing.T) {
	config.Initialize("../config/credentials.json")
	report.LoadReportConfiguration("../report/ReportConfiguration.xml")
	rabbit.InitializeRabbitConnection()
	models.InitializeConnections()
	device := &Device{}
	device.Initialize("355922061549941")
	packet := "$$355922061549941,6045,2018/10/15,01:54:38,56.75673,-111.43954,338,17,1.1,1,0,0,2,P0457,P0458,094##"
	message, _ := report.Parse(packet)

	if err := device.ProcessMessage(message); err != nil {
		t.Error("[TestDtcMessageProcessing]Handle message error:", err.Error())
	}

	if device.Activity.DTC.Count() != 2 {
		t.Error("[TestDtcMessageProcessing]Invalid DTC message processing")
	}

	packet = "$$355922063699116,6001,2018/11/01,14:01:24,57.31913,-111.50500,298,10.5,0.6,0.0,1247,167,18,1.2,1495,4.3,14.6,12,1,0,91.7,-156,4,-1031,158##"
	message, _ = report.Parse(packet)
	device.ProcessMessage(message)

	v, f := message.GetValue("DTCCode")
	if !f {
		t.Error("[TestDtcMessageProcessing]DTC Codes doesn't exists")
	}
	if len(v.([]string)) != 2 {
		t.Error("[TestDtcMessageProcessing]Invalid output DTC codes count")
	}
}

func TestPeriodicMessage(t *testing.T) {
	config.Initialize("../config/credentials.json")
	report.LoadReportConfiguration("../report/ReportConfiguration.xml")
	rabbit.InitializeRabbitConnection()
	models.InitializeConnections()
	device := &Device{}
	device.Initialize("355922064286244")
	packet := "$$355922064286244,4001,2018/11/14,08:47:38,51.03266,-114.02741,1047,37.9,0.0,0.0,982,41,14,1.1,42611,0.0,14.1,21,1,0,63.9,-371,-1059,-512,009##"
	message, _ := report.Parse(packet)

	err := device.ProcessMessage(message)
	if err != nil {
		t.Error("[TestLocationMessageProcessing]Handle message error:", err.Error())
	}
}

func TestOBDIIIgnitionOnMessage(t *testing.T) {
	config.Initialize("../config/credentials.json")
	report.LoadReportConfiguration("../report/ReportConfiguration.xml")
	rabbit.InitializeRabbitConnection()
	models.InitializeConnections()
	device := &Device{}
	device.Initialize("355922063875096")
	packet := "$$355922063875096,6011,2018/11/14,10:30:07,44.21442,-79.58988,304,1.2,1.1,0.0,965,85,12,1.4,8488,0.0,14.4,7,1,0.0,0,-969,-8,313,002##"
	message, _ := report.Parse(packet)

	err := device.ProcessMessage(message)
	if err != nil {
		t.Error("[TestOBDIIIgnitionOnMessage]Handle message error:", err.Error())
	}
}

func TestEmptyTimestampMessage(t *testing.T) {
	config.Initialize("../config/credentials.json")
	report.LoadReportConfiguration("../report/ReportConfiguration.xml")
	rabbit.InitializeRabbitConnection()
	models.InitializeConnections()
	device := &Device{}
	device.Initialize("355922063875096")
	packet := "$$355922061857732,4001,2000/00/00,00:00:00,0.00000,0.00000,0,8.0,0.0,0.5,630,0,0,0.0,4,13.7,14.2,14,0,0,76.7,47,348,922,010##"
	message, _ := report.Parse(packet)

	if _, f := message.GetValue("TimeStamp"); !f {
		t.Error("[TestEmptyTimestampMessage]TimeStamp not found")
	}

	v, _ := message.GetValue("TimeStamp")
	if v == nil {
		t.Error("[TestEmptyTimestampMessage]Invalid timestamp")
	}
	if _, valid := v.(*utils.JSONTime); !valid {
		t.Error("[TestEmptyTimestampMessage]Invalid timestamp")
	}
	err := device.ProcessMessage(message)
	if err != nil {
		t.Error("[TestEmptyTimestampMessage]Handle message error:", err.Error())
	}
}

func TestInvalidTimestampMessage(t *testing.T) {
	config.Initialize("../config/credentials.json")
	report.LoadReportConfiguration("../report/ReportConfiguration.xml")
	rabbit.InitializeRabbitConnection()
	models.InitializeConnections()
	device := &Device{}
	device.Initialize("355922063875096")
	packet := "$$355922063940353,4006,20118/11/14,13:34:10,57.39091,-111.06342,378,9.9,0.0,0.6,946,15,12,1.2,17539,0.0,14.5,22,1,1FTEW1E55JFB35598,AAz1-1168AA2,0b01,Unknown,0,90.0,-8,-78,-1055,002#"
	message, _ := report.Parse(packet)

	if _, f := message.GetValue("TimeStamp"); !f {
		t.Error("[TestEmptyTimestampMessage]TimeStamp not found")
	}

	v, _ := message.GetValue("TimeStamp")
	if v == nil {
		t.Error("[TestEmptyTimestampMessage]Invalid timestamp")
	}
	if _, valid := v.(*utils.JSONTime); !valid {
		t.Error("[TestEmptyTimestampMessage]Invalid timestamp")
	}
	err := device.ProcessMessage(message)
	if err != nil {
		t.Error("[TestEmptyTimestampMessage]Handle message error:", err.Error())
	}
}

//$$355922064340504,6011,2018/11/14,14:42:07,43.74029,-79.32384,129,1.3,0.3,0.0,709,195,15,1.0,5438,0.0,14.0,23,1,83.9,0,55,-8,-980,057##

func TestInvalidIgnitionMessage(t *testing.T) {
	config.Initialize("../config/credentials.json")
	report.LoadReportConfiguration("../report/ReportConfiguration.xml")
	rabbit.InitializeRabbitConnection()
	models.InitializeConnections()
	device := &Device{}
	device.Initialize("355922063875096")
	packet := "$$355922064340504,6011,2018/11/14,14:42:07,43.74029,-79.32384,129,1.3,0.3,0.0,709,195,15,1.0,5438,0.0,14.0,23,1,83.9,0,55,-8,-980,057##"
	message, _ := report.Parse(packet)

	err := device.ProcessMessage(message)
	if err != nil {
		t.Error("[TestEmptyTimestampMessage]Handle message error:", err.Error())
	}

	if ignitonState, valid := message.GetValue("IgnitionState"); !valid || ignitonState.(byte) != 1 {
		t.Error("[TestEmptyTimestampMessage]Invalid ignition state")
	}

}

func TestInvalidIgnitionPeriodicMessage(t *testing.T) {
	config.Initialize("../config/credentials.json")
	report.LoadReportConfiguration("../report/ReportConfiguration.xml")
	rabbit.InitializeRabbitConnection()
	models.InitializeConnections()
	device := &Device{}
	device.Initialize("355922064340504")
	packet := "$$355922064340504,4001,2018/11/14,14:45:07,43.74195,-79.31014,175,0.0,0.0,0.0,696,175,13,1.2,5439,11.2,14.7,23,1,0,83.9,0,63,-980,062##"
	message, _ := report.Parse(packet)

	err := device.ProcessMessage(message)
	if err != nil {
		t.Error("[TestEmptyTimestampMessage]Handle message error:", err.Error())
	}

	if ignitonState, valid := message.GetValue("IgnitionState"); !valid || ignitonState.(byte) != 1 {
		t.Error("[TestEmptyTimestampMessage]Invalid ignition state")
	}

}
*/
