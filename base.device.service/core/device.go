package core

import (
	"encoding/json"
	"log"
	"math"
	"strings"
	"time"

	"queclink-go/base.device.service/rabbit"
	"github.com/streadway/amqp"

	"queclink-go/base.device.service/config"
	"queclink-go/base.device.service/report"

	"queclink-go/base.device.service/comm"
	"queclink-go/base.device.service/core/immobilizer"
	"queclink-go/base.device.service/core/models"
	"queclink-go/base.device.service/utils"
)

//IDevice interface for device struct
type IDevice interface {
	Initialize(uniqueID string) error
	GetChannel() comm.IChannel
	SetChannel(channel comm.IChannel)
	GetConfiguration() IConfigurationManager
	SetConfiguration(configuration IConfigurationManager)
	GetIdentity() string
	SetIdentity(identity string)
	GetLastInitialization() time.Time
	GetActivity() *models.DeviceActivity
	GetVehicleSupport() *models.VehicleSupport
	SetOnDisconnect(func(c IDevice))
	OnChannelDisconnected(c comm.IChannel)
	OnReceivePacket(c comm.IChannel, packet []byte)
	SendSystemMessage(message report.IMessage, routingKey string)
	SendSystemMessageS(message string, routingKey string)

	ProcessMessage(message report.IMessage) error
	ProcessLocationMessage(message report.IMessage) (bool, error)
	ProcessAcknowledgementMessage(message report.IMessage) (bool, error)
	ProcessInfoMessage(message report.IMessage) (bool, error)
	ProcessDTCMessage(message report.IMessage) (bool, error)

	UpdateConfigurationQueue(message report.IMessage)
	SetDefaultValues(message report.IMessage)
	ValidateCoordinates(message report.IMessage)
	UpdateCurrentGpsTimeStamp(message report.IMessage)
	UpdateEventDependentStates(message report.IMessage)
	UpdateCurrentState(message report.IMessage)
	FillUpTimestamp(message report.IMessage)
	PublishMessage(message report.IMessage) error
	SaveActivity(message report.IMessage) error

	ValidateLocationMessage(message report.IMessage) bool

	SendString(command string)
	Send(command []byte)
	GetSelf() IDevice
	GetImmobilizer() immobilizer.IImmobilizer
}

//Device hardware device representation
type Device struct {
	Channel            comm.IChannel
	Configuration      IConfigurationManager
	Identity           string
	StartBit           int32
	LastInitialization time.Time
	Activity           *models.DeviceActivity
	VehicleSupport     *models.VehicleSupport
	onDisconnect       func(c IDevice)
	Immobilizer        immobilizer.IImmobilizer
	Self               IDevice
}

//ProcessMessage process incoming messages
func (device *Device) ProcessMessage(message report.IMessage) error {
	panic("Not implemented exception Device.ProcessMessage")
}

//OnReceivePacket from device
func (device *Device) OnReceivePacket(client comm.IChannel, packet []byte) {
	panic("Not implemented exception Device.OnReceivePacket")
}

//GetSelf return self inherited self instance
func (device *Device) GetSelf() IDevice {
	return device.Self
}

//GetImmobilizer returns instance of immobilizer handler
func (device *Device) GetImmobilizer() immobilizer.IImmobilizer {
	return device.Immobilizer
}

//ProcessAcknowledgementMessage handle ack message; Override on demand
func (device *Device) ProcessAcknowledgementMessage(message report.IMessage) (bool, error) {
	device.Configuration.UpdateMessageSent(message)
	device.Self.SendSystemMessage(message, config.Config.GetBase().Rabbit.OrleansDebugRoutingKey)
	return true, nil
}

//ProcessInfoMessage handle info message; Override on demand
func (device *Device) ProcessInfoMessage(message report.IMessage) (bool, error) {
	device.UpdateEventDependentStates(message)
	message.SetValue("PowerState", device.Activity.PowerState)
	device.Self.SendSystemMessage(message, config.Config.GetBase().Rabbit.OrleansDebugRoutingKey)
	return true, nil
}

//ProcessLocationMessage handle location message
func (device *Device) ProcessLocationMessage(message report.IMessage) (bool, error) {
	if !device.Self.ValidateLocationMessage(message) {
		return false, nil
	}
	device.Self.UpdateCurrentState(message)
	err := device.Self.PublishMessage(message)
	return true, err
}

//ProcessDTCMessage handle dtc message
func (device *Device) ProcessDTCMessage(message report.IMessage) (bool, error) {
	if device.Activity != nil && device.Activity.LastMessage != "" {
		iTs, found := message.GetValue("TimeStamp")
		if !found {
			message.SetValue("TimeStamp", &utils.JSONTime{Time: device.Activity.GPSTimeStamp})
			message.SetValue("GpsValidity", byte(0))
		} else {
			ts, _ := iTs.(*utils.JSONTime)
			if device.Activity.GPSTimeStamp.After(ts.Time) {
				message.SetValue("TimeStamp", &utils.JSONTime{Time: device.Activity.GPSTimeStamp})
				message.SetValue("GpsValidity", byte(0))
				message.SetValue("Odometer", device.Activity.Odometer)
			}
		}
	} else {
		message.SetValue("Odometer", device.Activity.Odometer)
	}
	message.SetValue("PowerState", device.Activity.PowerState)
	device.Self.UpdateCurrentState(message)
	err := device.Self.PublishMessage(message)
	return true, err
}

//UpdateVehicleInfoMsg update supported features for vehicle; Override on demand
func (device *Device) UpdateVehicleInfoMsg(message report.IMessage) {
	device.VehicleSupport.UpdateState(device.Identity, message)
	device.VehicleSupport.Sync(device.Identity, message.GetStringValue("VIN", ""))
}

//SendSystemMessage sends message via rabbit
func (device *Device) SendSystemMessage(message report.IMessage, routingKey string) {
	jMessage, _ := json.Marshal(message)
	headers := make(amqp.Table, 0)
	rabbit.RabbitConnection.Publish(string(jMessage), config.Config.GetBase().Rabbit.SystemExchange, routingKey, headers, 10)
}

//SendSystemMessageS sends info message to via rabbit
func (device *Device) SendSystemMessageS(message string, routingKey string) {
	dto := report.NewMessage()
	dto.SetValue("Message", message)
	dto.SetValue("DevId", device.Identity)
	dto.SetValue("TimeStamp", &utils.JSONTime{Time: time.Now().UTC()})
	dto.SetValue("LocationMessage", false)
	device.Self.SendSystemMessage(dto, routingKey)
}

//GetChannel returns current device channel(TCP/UDP)
func (device *Device) GetChannel() comm.IChannel {
	return device.Channel
}

//SetChannel sets data channell for device
func (device *Device) SetChannel(channel comm.IChannel) {
	device.Channel = channel
}

//GetConfiguration returns instance of configuration manager for device
func (device *Device) GetConfiguration() IConfigurationManager {
	return device.Configuration
}

//SetConfiguration assign configuration manager to device instance
func (device *Device) SetConfiguration(configuration IConfigurationManager) {
	device.Configuration = configuration
}

//GetIdentity returns identity for device
func (device *Device) GetIdentity() string {
	return device.Identity
}

//SetIdentity sets identity for device
func (device *Device) SetIdentity(identity string) {
	device.Identity = identity
}

//GetLastInitialization returns last refresh date/time
func (device *Device) GetLastInitialization() time.Time {
	return device.LastInitialization
}

//GetActivity returns device activity[last message timestamp, software, last message etc]
func (device *Device) GetActivity() *models.DeviceActivity {
	return device.Activity
}

//GetVehicleSupport return supported ECU feature(vin, fuel, DTC codes etc)
func (device *Device) GetVehicleSupport() *models.VehicleSupport {
	return device.VehicleSupport
}

//SetOnDisconnect assign callback on device disconnect
func (device *Device) SetOnDisconnect(onDisconnect func(c IDevice)) {
	device.onDisconnect = onDisconnect
}

//UpdateConfigurationQueue handle incoming message by configuration manager
func (device *Device) UpdateConfigurationQueue(message report.IMessage) {
	device.Configuration.UpdateMessageSent(message)
	device.Configuration.SyncDeviceConfig()
}

//FillUpTimestamp checks for Timestamp field exists
func (device *Device) FillUpTimestamp(message report.IMessage) {
	t, f := message.GetValue("Timestamp")
	if !f {
		message.SetValue("Timestamp", &utils.JSONTime{Time: device.Activity.TimeStamp})
	}
	if time.Now().UTC().Sub(t.(*utils.JSONTime).Time).Hours() > 8760 {
		message.SetValue("Timestamp", &utils.JSONTime{Time: device.Activity.TimeStamp})
	}
}

//SetDefaultValues sets values for fields Supply, GPIO, IgnitionState, Odometer, PowerState, ReceivedTime
func (device *Device) SetDefaultValues(message report.IMessage) {
	if _, found := message.GetValue("Supply"); !found {
		message.SetValue("Supply", int32(0))
	}

	if _, found := message.GetValue("GPIO"); !found {
		message.SetValue("GPIO", byte(0))
	}

	if _, found := message.GetValue("IgnitionState"); !found {
		if device.Activity.Ignition == "On" {
			message.SetValue("IgnitionState", byte(1))
		} else {
			message.SetValue("IgnitionState", byte(0))
		}
	}

	if _, found := message.GetValue("Odometer"); !found {
		message.SetValue("Odometer", device.Activity.Odometer)
	}

	message.SetValue("PowerState", device.Activity.PowerState)
	message.SetValue("ReceivedTime", &utils.JSONTime{Time: time.Now().UTC()})
}

//UpdateCurrentGpsTimeStamp gps timestamp validation
func (device *Device) UpdateCurrentGpsTimeStamp(message report.IMessage) {
	if v, found := message.GetValue("GpsValidity"); found {
		validity := v.(byte)
		t, _ := message.GetValue("TimeStamp")
		timestamp := t.(*utils.JSONTime)
		if validity == 1 &&
			timestamp.After(time.Now().UTC().AddDate(-5, 0, 0)) &&
			timestamp.Before(time.Now().UTC().Add(10*time.Minute)) &&
			timestamp.After(device.Activity.GPSTimeStamp) {

			device.Activity.GPSTimeStamp = timestamp.Time

		} else if timestamp.After(time.Now().UTC().Add(10 * time.Minute)) {
			message.SetValue("TimeStamp", &utils.JSONTime{Time: device.Activity.GPSTimeStamp})
		}
	}
}

//UpdateCurrentState assign values from message to current state and fill up missing fields for message
func (device *Device) UpdateCurrentState(message report.IMessage) {
	device.Self.SetDefaultValues(message)
	device.Self.UpdateCurrentGpsTimeStamp(message)
	device.Self.ValidateCoordinates(message)
	device.Self.UpdateEventDependentStates(message)
	device.Activity.TimeStamp = time.Now().UTC()

	if device.Activity.DTC != nil && len(device.Activity.DTC.Codes) > 0 {
		message.SetValue("DTCCode", device.Activity.DTC.Codes)
	}

	if v, f := message.GetValue("Odometer"); f {
		device.Activity.Odometer = v.(int32)
	}

	if v, f := message.GetValue("Satellites"); f {
		device.Activity.SatFix = v.(int32)
	}
	if v, f := message.GetValue("BatteryPercentage"); f {
		device.Activity.BatteryLevel = v.(float32)
	}
	if v, f := message.GetValue("Relay"); f {
		device.Activity.Relay = v.(byte)
	}
}

//UpdateEventDependentStates make map events on states like PowerState, Fuel etc
func (device *Device) UpdateEventDependentStates(message report.IMessage) {
	if device.Activity.PowerState == "Backup battery" && message.GetIntValue("Supply", 0) > 8000 {
		device.Activity.PowerState = "Powered"
	}

	message.SetValue("PowerState", device.Activity.PowerState)

	if v, f := message.GetValue("FuelLevel"); f && device.Activity.Ignition == "On" {
		fv, _ := v.(float32)
		if fv != 0.0 {
			device.Activity.FuelLevel = fv
		}
	}

	if device.Activity.FuelLevel == 0 || (device.VehicleSupport.Uploaded && !device.VehicleSupport.FuelLevelSupport) {
		message.RemoveKey("FuelLevel")
		device.Activity.FuelLevel = 0
	}
	if device.Activity.FuelLevel != 0 {
		message.SetValue("FuelLevel", device.Activity.FuelLevel)
	}
}

//ValidateCoordinates checks coordinate for valid value
func (device *Device) ValidateCoordinates(message report.IMessage) {
	lt, f := message.GetValue("Latitude")
	if !f {
		return
	}
	lat := lt.(float32)
	ln, f1 := message.GetValue("Longitude")
	if !f1 {
		return
	}
	lon := ln.(float32)
	v, _ := message.GetValue("GpsValidity")
	vld := v.(byte)
	if math.Abs(float64(lat)) > 1e-10 && math.Abs(float64(lon)) > 1e-10 && vld == 1 {
		device.Activity.Latitude = lat
		device.Activity.Longitude = lon
	} else {
		if vld != 1 {
			if math.IsNaN(float64(lat)) || lat < -90.0 ||
				lat > 90.0 || math.Abs(float64(lat)) < 1e-10 {
				message.SetValue("Latitude", device.Activity.Latitude)
				message.SetValue("Longitude", device.Activity.Longitude)
			}
			if math.IsNaN(float64(lon)) || lon < -180.0 ||
				lon > 180.0 || math.Abs(float64(lon)) < 1e-10 {
				message.SetValue("Latitude", device.Activity.Latitude)
				message.SetValue("Longitude", device.Activity.Longitude)
			}
		}
	}
	if al, found := message.GetValue("Altitude"); found {
		alt := al.(float32)
		if math.IsNaN(float64(alt)) || math.Abs(float64(alt)) < 1e-10 {
			message.SetValue("Altitude", float32(0.0))
		}
	}
	if h, found := message.GetValue("Heading"); found &&
		(h.(float32) > 360 || h.(float32) < 0) {
		message.SetValue("Heading", float32(0))
	}
}

//PublishMessage uploads message to database, rabbit etc
func (device *Device) PublishMessage(message report.IMessage) error {

	if config.Config.GetBase().Storage != "" {
		message.SetValue("storage", config.Config.GetBase().Storage)
	}

	message.SetValue("PrevSourceId", device.Activity.LastMessageID)

	if herr := device.saveMessageHistory(message); herr != nil {
		return herr
	}
	if aerr := device.Self.SaveActivity(message); aerr != nil {
		return aerr
	}
	device.Self.SendSystemMessage(message, models.GetMessageHistoryTableName(device.Identity))
	return nil
}

//SaveActivity save activity to database
func (device *Device) SaveActivity(message report.IMessage) error {
	timestamp, _ := message.GetValue("TimeStamp")
	device.Activity.MessageTime = timestamp.(*utils.JSONTime).Time
	device.Activity.LastUpdateTime = time.Now().UTC()
	device.Activity.LastMessageID = message.SourceID()
	jMessage, _ := json.Marshal(message)
	device.Activity.LastMessage = string(jMessage)
	return device.Activity.Save()
}

//ValidateLocationMessage checks message is location and contains coordinates
func (device *Device) ValidateLocationMessage(message report.IMessage) bool {
	if !message.LocationMessage() {
		return false
	}
	if _, f := message.GetValue("Latitude"); !f {
		return false
	}
	if _, f1 := message.GetValue("Longitude"); !f1 {
		return false
	}
	return true
}

//SendString command to device using channel
func (device *Device) SendString(command string) {
	if len(command) == 0 {
		return
	}

	if device.Channel != nil {
		device.Channel.Send(command)
		device.Self.SendSystemMessageS(command, config.Config.GetBase().Rabbit.OrleansDebugRoutingKey)
		log.Println("[Device]SendCommand: Send message to device ", device.Identity, " Message:", command)
	} else {
		log.Println("[Device]SendCommand: Unable send message to device ", device.Identity, " Message:", command)
	}
}

//Send command to device using channel
func (device *Device) Send(command []byte) {
	if len(command) == 0 {
		return
	}

	if device.Channel != nil {
		device.Channel.SendBytes(command)
		device.Self.SendSystemMessageS(string(command), config.Config.GetBase().Rabbit.OrleansDebugRoutingKey)
		log.Println("[Device]SendBytes: Send message to device ", device.Identity, " Message:", string(command))
	} else {
		log.Println("[Device]SendBytes: Unable send message to device ", device.Identity, " Message:", string(command))
	}
}

//Initialize device
func (device *Device) Initialize(uniqueID string) error {
	identity := strings.Replace(uniqueID, "imei:", "", -1)
	device.Identity = "xirgo_" + identity
	if time.Now().UTC().Sub(device.LastInitialization).Seconds() < 30 {
		return nil
	}
	device.LastInitialization = time.Now().UTC()
	device.Activity, _ = models.FindDeviceActivityInfo(device.Identity)
	device.Activity.Identity = device.Identity
	device.VehicleSupport, _ = models.FindVehicleSupport(device.Identity, "")
	device.Configuration = InitializeConfigurationManager(device)
	device.Configuration.Load()
	return nil
}

//OnChannelDisconnected indicates client disconnected
func (device *Device) OnChannelDisconnected(c comm.IChannel) {
	device.onDisconnect(device)
}

func (device *Device) saveMessageHistory(message report.IMessage) error {
	h := &models.MessageHistory{
		DevID: device.Identity,
	}
	_, err := h.Save(message)
	return err
}
