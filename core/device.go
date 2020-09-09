package core

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"queclink-go/core/immobilizer"
	"queclink-go/qconfig"
	"queclink-go/queclinkreport"
	"queclink-go/queclinkreport/fields"
	"queclink-go/queclinkreport/gv55/evt55"
	"queclink-go/queclinkreport/gv55/resp55"

	"queclink-go/base.device.service/comm"
	"queclink-go/base.device.service/config"
	"queclink-go/base.device.service/core"
	"queclink-go/base.device.service/core/models"
	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
)

//QDevice base queclink device
type QDevice struct {
	core.Device
	Statistic        *DeviceStatistic
	LastGpsMessage   report.IMessage
	DogMessage       report.IMessage
	Type             int
	RtoSentAt        time.Time
	handleRspMessage func(message report.IMessage) (bool, error)
	handleEvtMessage func(message report.IMessage) (bool, error)
}

//OnReceivePacket callback implementation
func (device *QDevice) OnReceivePacket(client comm.IChannel, packet []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("panic:Recovered in OnReceivePacket:", r)
		}
	}()

	log.Println("Received packet:", utils.InsertNth(utils.ByteToString(packet), 2, ' '), "from:", device.Identity)

	p := &queclinkreport.Parser{}
	messages, err := p.Parse(packet)

	if err != nil {
		return
	}

	for _, message := range messages {
		jMessage, jerr := json.Marshal(message)
		if jerr == nil {
			log.Println("Received packet:", string(jMessage), "from:", device.Identity)
		}

		ack, found := message.GetValue("Ack")
		err = device.ProcessMessage(message)
		if found && err == nil {
			device.Channel.SendBytes(ack.([]byte))
		}
	}
}

//ProcessMessage process incoming messages
func (device *QDevice) ProcessMessage(message report.IMessage) error {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in ProcessMessage:", r)
		}
	}()

	qmsg, valid := message.(*queclinkreport.QueclinkMessage)
	if !valid || qmsg == nil {
		return fmt.Errorf("[ProcessMessage] Invalid message for device:%v", device.Identity)
	}

	switch qmsg.MessageType() {
	case "+INF", "+BNF":
		{
			device.processOutputStatusMess(qmsg)
			break
		}
	case "+ACK":
		{
			device.processAcknowledgementMessage(qmsg)
			break
		}
	case "+HBD", "+BBD":
		{
			device.processHbdMessage(qmsg)
			break
		}
	default:
		{
			device.handleLocationMessage(qmsg)
			device.Immobilizer.Update()
			break
		}
	}
	if qmsg.MessageType() != "+BSP" && qmsg.MessageType() != "+BVT" {
		device.Configuration.SyncDeviceConfig()
	}

	return nil
}

func (device *QDevice) processOutputStatusMess(message report.IMessage) {
	switch message.MessageType() {
	case "+INF", "+BNF":
		{
			switch message.EventCode() {
			case 8:
				{
					device.updateOutputStatus(message)
					device.Immobilizer.CheckStatusOfImmobilizerCommands()
					break
				}
			}
			break
		}
	}
	device.SendSystemMessage(message, config.Config.GetBase().Rabbit.OrleansDebugRoutingKey)
}

func (device *QDevice) updateOutputStatus(message report.IMessage) bool {
	if r, found := message.GetValue("Relay"); found {
		if vr, valid := r.(byte); valid {
			device.Activity.Relay = vr
			return true
		}
	}
	return false
}

func (device *QDevice) handleLocationMessage(message report.IMessage) (bool, error) {
	switch message.MessageType() {
	case "+CRD", "+BRD":
		return device.handleCrdMessage(message)
	case "+RSP", "+BSP":
		return device.handleRspMessage(message)
	case "+EVT", "+BVT":
		return device.handleEvtMessage(message)
	}
	return true, nil
}

func (device *QDevice) processAcknowledgementMessage(message report.IMessage) (bool, error) {
	message.SetValue("IP", device.GetChannel().RemoteIP())
	message.SetValue("Port", device.GetChannel().RemotePort())
	device.Immobilizer.UpdateMessageSent(message)

	ack, ackFound := message.GetValue("Ack")
	if ackFound {
		go device.delayAck(3, ack)
	}

	return device.Device.ProcessAcknowledgementMessage(message)
}

func (device *QDevice) delayAck(n time.Duration, ack interface{}) {
	c := time.Tick(n * time.Second)
	for range c {
		break
	}
	device.GetChannel().SendBytes(ack.([]byte))
}

func (device *QDevice) processHbdMessage(message report.IMessage) {
	device.Activity.Save()
	device.SendSystemMessage(message, config.Config.GetBase().Rabbit.OrleansDebugRoutingKey)
	if device.RtoSentAt.Add(4 * time.Hour).Before(time.Now().UTC()) {
		return
	}

	threshold := config.Config.(*qconfig.QConfiguration).RTLThreshold

	if device.RtoSentAt.Add(time.Duration(threshold) * time.Minute).After(time.Now().UTC()) {
		return
	}
	ist, found := message.GetValue("SendTime")
	if !found {
		return
	}
	sendTime, valid := ist.(*utils.JSONTime)
	if !valid {
		return
	}

	lastMessage, err := queclinkreport.UnMarshalMessage(device.Activity.LastMessage)
	if err != nil {
		return
	}

	reportConfig, err := report.ReportConfiguration.(*queclinkreport.ReportConfiguration).Find(device.Type, lastMessage.MessageType(), int(lastMessage.EventCode()))
	if err != nil || reportConfig == nil {
		log.Fatalf("[Device]processHbdMessage. Not found configuration for device type %v message header:%v message type:%v",
			device.Type, lastMessage.MessageType(), lastMessage.EventCode())
		return
	}

	if lastMessage.MessageType() == "+BSP" ||
		lastMessage.MessageType() == "+BVT" ||
		!reportConfig.Location {
		return
	}

	ilst, found := lastMessage.GetValue("SendTime")
	if !found {
		return
	}

	lastSendTime, valid := ilst.(*utils.JSONTime)

	if lastSendTime.AddDate(0, 0, 1).Before(time.Now().UTC()) {
		return
	}

	if lastSendTime.Add(time.Duration(threshold) * time.Minute).After(sendTime.Time) {
		return
	}
	log.Printf("[Device]SendRto for device:%v; last send time:%v; current send time:%v",
		device.Identity, lastSendTime.Time.String(), sendTime.Time.String())
	device.RtoSentAt = time.Now().UTC()

	deviceConfigs := report.ReportConfiguration.(*queclinkreport.ReportConfiguration)
	deviceConfig, err := deviceConfigs.FindDeviceType(device.Type)
	if err != nil {
		log.Fatalf("[Device]processHbdMessage. Not found configuration for device type %v", device.Type)
		return
	}

	item, err := deviceConfig.FindAckMessagesType("AT+GTRTO")
	if err != nil || item == nil {
		log.Fatalf("[Device]processHbdMessage. Not found ack configuration for device type %v", device.Type)
		return

	}
	device.Configuration.AddCommand(&core.ConfigurationItem{
		MessageType: fmt.Sprintf("%v", item.ID),
		Command:     "AT+GTRTO=gv55,1,,,,,,FFFF$",
	}, true)
}

func (device *QDevice) handleCrdMessage(message report.IMessage) (bool, error) {
	if device.LastGpsMessage == nil {
		device.SendSystemMessage(message, config.Config.GetBase().Rabbit.OrleansDebugRoutingKey)
		return true, nil
	}

	if err := device.replaceMessageValidity(device.LastGpsMessage, message); err != nil {
		return true, nil
	}

	if v, f := device.LastGpsMessage.GetValue(fields.DigitalInputStatus); f {
		message.SetValue(fields.DigitalInputStatus, v)
	} else {
		return true, nil
	}
	if v, f := device.LastGpsMessage.GetValue(fields.IgnitionState); f {
		message.SetValue(fields.IgnitionState, v)
	} else {
		return true, nil
	}

	if v, f := device.LastGpsMessage.GetValue(fields.Relay); f {
		message.SetValue(fields.Relay, v)
	} else {
		return true, nil
	}

	if v, f := device.LastGpsMessage.GetValue(fields.Relay); f {
		message.SetValue(fields.Relay, v)
	} else {
		return true, nil
	}

	if v, f := device.LastGpsMessage.GetValue(fields.Altitude); f {
		message.SetValue(fields.Altitude, v)
	} else {
		return true, nil
	}
	if v, f := device.LastGpsMessage.GetValue(fields.Heading); f {
		message.SetValue(fields.Heading, v)
	} else {
		return true, nil
	}

	if v, f := device.LastGpsMessage.GetValue(fields.BatteryPercentage); f {
		message.SetValue(fields.BatteryPercentage, v)
	} else {
		return true, nil
	}

	if v, f := device.LastGpsMessage.GetValue(fields.CurrentMileage); f {
		message.SetValue(fields.CurrentMileage, v)
	} else {
		return true, nil
	}

	if v, f := device.LastGpsMessage.GetValue(fields.Odometer); f {
		message.SetValue(fields.Odometer, v)
	} else {
		return true, nil
	}

	if v, f := device.LastGpsMessage.GetValue(fields.Supply); f {
		message.SetValue(fields.Supply, v)
	} else {
		return true, nil
	}
	if v, f := device.LastGpsMessage.GetValue(fields.IgnitionStatus); f {
		message.SetValue(fields.IgnitionStatus, v)
	} else {
		return true, nil
	}

	if v, f := device.LastGpsMessage.GetValue(fields.GPSUTCTime); f {
		message.SetValue(fields.GPSUTCTime, v)
	} else {
		return true, nil
	}

	if v, f := device.LastGpsMessage.GetValue(fields.Speed); f {
		message.SetValue(fields.Speed, v)
	} else {
		return true, nil
	}

	if v, f := device.LastGpsMessage.GetValue(fields.TimeStamp); f {
		message.SetValue(fields.TimeStamp, v)
	} else {
		return true, nil
	}
	device.Self.ProcessLocationMessage(message)
	return true, nil
}

func (device *QDevice) rspMessageHandler(message report.IMessage) (bool, error) {

	if message.EventCode() == resp55.GTDOG {
		device.DogMessage = message
		device.handleRspMessage = device.rspIgnoreMessageHandler
		device.handleEvtMessage = device.evtIgnoreMessageHandler
	}
	device.syncMessageValidity(message)
	device.Self.ProcessLocationMessage(message)
	return true, nil
}

func (device *QDevice) rspIgnoreMessageHandler(message report.IMessage) (bool, error) {

	if message.EventCode() == resp55.GTDOG {
		device.DogMessage = message
		device.handleRspMessage = device.rspIgnoreMessageHandler
		device.handleEvtMessage = device.evtIgnoreMessageHandler
	}
	if device.DogMessage == nil {
		device.handleRspMessage = device.rspMessageHandler
		device.handleEvtMessage = device.evtMessageHandler
		return true, nil
	}

	idst, found := device.DogMessage.GetValue(fields.SendTime)
	if !found {
		idst = &utils.JSONTime{Time: utils.MinTimeStamp()}
	}

	dst, valid := idst.(*utils.JSONTime)
	if !valid {
		dst = &utils.JSONTime{Time: utils.MinTimeStamp()}
	}

	imst, found := message.GetValue(fields.SendTime)
	if !found {
		idst = &utils.JSONTime{Time: utils.MinTimeStamp()}
	}

	mst, valid := imst.(*utils.JSONTime)
	if !valid {
		mst = &utils.JSONTime{Time: utils.MinTimeStamp()}
	}

	if dst.Add(5 * time.Minute).Before(mst.Time) {
		device.handleRspMessage = device.rspMessageHandler
		device.handleEvtMessage = device.evtMessageHandler
	}
	if message.EventCode() == resp55.GTEPS {
		message.SetValue("LocationMessage", false)
	}
	device.rspMessageHandler(message)
	return true, nil
}

func (device *QDevice) evtMessageHandler(message report.IMessage) (bool, error) {
	if message.LocationMessage() {
		device.updatePowerState(message)
	}
	device.syncMessageValidity(message)
	device.Self.ProcessLocationMessage(message)
	return true, nil
}

func (device *QDevice) evtIgnoreMessageHandler(message report.IMessage) (bool, error) {

	if device.DogMessage == nil {
		device.handleRspMessage = device.rspMessageHandler
		device.handleEvtMessage = device.evtMessageHandler
		return true, nil
	}

	idst, found := device.DogMessage.GetValue(fields.SendTime)
	if !found {
		idst = &utils.JSONTime{Time: utils.MinTimeStamp()}
	}

	dst, valid := idst.(*utils.JSONTime)
	if !valid {
		dst = &utils.JSONTime{Time: utils.MinTimeStamp()}
	}

	imst, found := message.GetValue(fields.SendTime)
	if !found {
		idst = &utils.JSONTime{Time: utils.MinTimeStamp()}
	}

	mst, valid := imst.(*utils.JSONTime)
	if !valid {
		mst = &utils.JSONTime{Time: utils.MinTimeStamp()}
	}

	if dst.Add(5 * time.Minute).Before(mst.Time) {
		device.handleRspMessage = device.rspMessageHandler
		device.handleEvtMessage = device.evtMessageHandler
	}

	switch message.EventCode() {
	case evt55.GTPFA, evt55.GTPNA,
		evt55.GTMPN, evt55.GTMPF,
		evt55.GTIGN, evt55.GTIGF:
		{
			message.SetValue("LocationMessage", false)
		}
	}

	device.evtMessageHandler(message)
	return true, nil
}

//SetDefaults for device
func (device *QDevice) SetDefaults(deviceType int) {
	device.handleRspMessage = device.rspMessageHandler
	device.handleEvtMessage = device.evtMessageHandler
	device.Statistic = NewDeviceStatistic()
	device.Immobilizer = immobilizer.Initialize(device)
	device.Type = deviceType
	device.Self = device
}

//UpdateCurrentState assign values from message to current state and fill up missing fields for message
func (device *QDevice) UpdateCurrentState(message report.IMessage) {
	device.updateDefPowerState(message)
	device.updateOutputStatus(message)
	if message.LocationMessage() {
		device.UpdateIgnitionState(message)
	}
	device.Device.UpdateCurrentState(message)
}

//UpdateIgnitionState for ql device
func (device *QDevice) UpdateIgnitionState(message report.IMessage) {
	iv, found := message.GetValue(fields.DigitalInputStatus)
	if !found {
		return
	}
	diStatus, valid := iv.(byte)
	if !valid {
		return
	}

	iis, found := message.GetValue(fields.IgnitionStatus)
	if !found {
		return
	}
	iStatus, valid := iis.(byte)
	if !valid {
		return
	}

	if (diStatus&1) == 1 || iStatus == 32 {
		device.Activity.Ignition = "On"
	} else {
		device.Activity.Ignition = "Off"
	}
}

func (device *QDevice) updateDefPowerState(message report.IMessage) {
	if message.GetIntValue("Supply", 0) > 0 {
		device.Activity.PowerState = "Powered"
	} else {
		device.Activity.PowerState = "Backup battery"
	}
}

func (device *QDevice) updatePowerState(message report.IMessage) {
	switch message.EventCode() {
	case evt55.GTPNA, evt55.GTMPN:
		device.Activity.PowerState = "Powered"
		break
	case evt55.GTPFA:
		device.Activity.PowerState = "Power off"
		break
	case evt55.GTMPF:
		device.Activity.PowerState = "Backup battery"
		break
	case evt55.GTBPL:
		device.Activity.PowerState = "Backup battery"
		break
	default:
		device.updateDefPowerState(message)
		break
	}
}

func (device *QDevice) replaceMessageValidity(from report.IMessage, to report.IMessage) error {

	if v, f := from.GetValue(fields.GPSAccuracy); f {
		to.SetValue(fields.GPSAccuracy, v)
	} else {
		return fmt.Errorf("Not found field:%v", fields.GPSAccuracy)
	}
	if v, f := from.GetValue(fields.GpsValidity); f {
		to.SetValue(fields.GpsValidity, v)
	} else {
		return fmt.Errorf("Not found field:%v", fields.GpsValidity)
	}
	if v, f := from.GetValue(fields.Satellites); f {
		to.SetValue(fields.Satellites, v)
	} else {
		return fmt.Errorf("Not found field:%v", fields.Satellites)
	}
	if v, f := from.GetValue(fields.Latitude); f {
		to.SetValue(fields.Latitude, v)
	} else {
		return fmt.Errorf("Not found field:%v", fields.Latitude)
	}
	if v, f := from.GetValue(fields.Longitude); f {
		to.SetValue(fields.Longitude, v)
	} else {
		return fmt.Errorf("Not found field:%v", fields.Longitude)
	}
	return nil
}

func (device *QDevice) syncMessageValidity(message report.IMessage) {
	switch message.MessageType() {
	case "+RSP", "+BSP":
		{
			reportConfig, err := report.ReportConfiguration.(*queclinkreport.ReportConfiguration).Find(device.Type,
				message.MessageType(), int(message.EventCode()))
			if err != nil {
				return
			}

			if !reportConfig.Location && device.LastGpsMessage != nil {
				device.replaceMessageValidity(device.LastGpsMessage, message)
			} else {
				device.LastGpsMessage = message
			}
		}
	case "+EVT", "+BVT":
		reportConfig, err := report.ReportConfiguration.(*queclinkreport.ReportConfiguration).Find(device.Type,
			message.MessageType(), int(message.EventCode()))
		if err != nil {
			return
		}

		if !reportConfig.Location && device.LastGpsMessage != nil {
			if v, f := message.GetValue(fields.GpsValidity); f && v == 1 && message.EventCode() == evt55.GTSTT {
				device.LastGpsMessage = message
			} else {
				device.replaceMessageValidity(device.LastGpsMessage, message)
			}
		} else {
			device.LastGpsMessage = message
		}
		return
	}
}

func (device *QDevice) updateSoftware(message report.IMessage) {
	ipv, found := message.GetValue(fields.ProtocolVersion)
	if !found {
		return
	}
	pv, valid := ipv.(int32)
	if !valid {
		return
	}

	spv := fmt.Sprintf("%v", pv)
	if device.Activity.Software.SoftwareVersion != spv {
		device.Activity.Software.SoftwareVersion = spv
		device.Activity.Software.UpdatedAt = &utils.JSONTime{Time: time.Now().UTC()}
	}

	ifv, found := message.GetValue(fields.FirmwareVersion)
	if !found {
		return
	}
	fv, valid := ifv.(int32)
	if !valid {
		return
	}

	sfv := fmt.Sprintf("%v", fv)
	if device.Activity.Software.FirmwareVersion != sfv {
		device.Activity.Software.FirmwareVersion = sfv
		device.Activity.Software.UpdatedAt = &utils.JSONTime{Time: time.Now().UTC()}
	}
}

//SaveActivity save activity to database
func (device *QDevice) SaveActivity(message report.IMessage) error {
	device.updateSoftware(message)
	return device.Device.SaveActivity(message)
}

//SetDefaultValues sets values for fields Supply, GPIO, IgnitionState, Odometer, PowerState, ReceivedTime
func (device *QDevice) SetDefaultValues(message report.IMessage) {
	device.Device.SetDefaultValues(message)
	message.SetValue("IP", device.GetChannel().RemoteIP())
	message.SetValue("Port", device.GetChannel().RemotePort())
}

//Initialize device
func (device *QDevice) Initialize(uniqueID string) error {
	identity := strings.Replace(uniqueID, "imei:", "", -1)
	device.Identity = "queclink_" + identity
	if time.Now().UTC().Sub(device.LastInitialization).Seconds() < 30 {
		return nil
	}
	device.LastInitialization = time.Now().UTC()
	device.Activity, _ = models.FindDeviceActivityInfo(device.Identity)
	device.Activity.Identity = device.Identity
	device.VehicleSupport, _ = models.FindVehicleSupport(device.Identity, "")
	device.Configuration = InitializeConfigurationManager(device.Type, device)
	device.Configuration.SetSelf(device.Configuration)
	device.Configuration.Load()
	return nil
}
