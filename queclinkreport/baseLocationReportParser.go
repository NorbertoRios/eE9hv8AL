package queclinkreport

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
	"queclink-go/queclinkreport/fields"
)

//BaseLocationReportParser location message parser
type BaseLocationReportParser struct {
	ReportParser
}

//Parse default packet
func (parser *BaseLocationReportParser) Parse(packet []byte, config *Report) []report.IMessage {
	messages := []report.IMessage{}
	readFrom := parser.StartByte
	header := parser.ParseHeader(parser.ReportMask, packet, config.Header, &readFrom)
	header.SetValue(fields.MessageHeader, parser.MessageHeader)
	header.SetValue(fields.MessageType, parser.MessageType)
	header.SetValue(fields.ReportMask, parser.ReportMask)
	header.SetValue("RawData", packet)

	number := byte(0)
	if v, found := header.GetValue(fields.Number); found {
		number = v.(byte)
	}

	for i := 1; i <= int(number); i++ {
		body := parser.ParseBody(parser.ReportMask, packet, config.MultiPosition, &readFrom)
		messages = append(messages, body)
	}
	var tail = parser.ParseTail(parser.ReportMask, packet, config.Tail, &readFrom)

	for _, report := range messages {
		report.AppendRange(*header.GetData())
		report.AppendRange(*tail.GetData())
	}

	if len(messages) == 0 && (len(*header.GetData()) != 0 || len(*tail.GetData()) != 0) {
		header.AppendRange(*tail.GetData())
		messages = append(messages, header)
	}

	if _, found := header.GetValue(fields.Number); found && number == 0 {

		v, _ := utils.GetBytesValue(packet, len(packet)-6, 2)
		header.SetValue(fields.CountNumber, v)
		jMessage, _ := json.Marshal(messages)
		log.Println("Recovered binary: ", utils.InsertNth(utils.ByteToString(packet), 2, ' '))
		log.Println("Recovered json: ", string(jMessage))
	}

	return messages
}

//ParseHeader of message
func (parser *BaseLocationReportParser) ParseHeader(reportMask int32, packet []byte, confItems []Item, readFrom *int) IQueclinkMessage {
	return parser.ParsePart(reportMask, packet, confItems, readFrom)
}

//ParseBody of message
func (parser *BaseLocationReportParser) ParseBody(reportMask int32, packet []byte, confItems []Item, readFrom *int) IQueclinkMessage {
	return parser.ParsePart(reportMask, packet, confItems, readFrom)
}

//ParseTail of message
func (parser *BaseLocationReportParser) ParseTail(reportMask int32, packet []byte, confItems []Item, readFrom *int) IQueclinkMessage {
	return parser.ParsePart(reportMask, packet, confItems, readFrom)
}

//ParsePart parses part of packet and returns IQueclinkMessage
func (parser *BaseLocationReportParser) ParsePart(reportMask int32, packet []byte, confItems []Item, readFrom *int) IQueclinkMessage {

	message := NewQueclinkMessage()
	//in config.Where(reportItem => ByteIsSet(ReportMask, reportItem.MaskBit))
	for _, item := range confItems {
		if item.MaskBit != nil && !utils.BitIsSet(int64(reportMask), uint(*item.MaskBit)) {
			continue
		}
		value, _ := utils.GetBytesValue(packet, *readFrom, item.Size)
		parser.OnValueParsed(confItems, item.ItemName, value)
		message.SetValue(item.ItemName, value)
		*readFrom = *readFrom + item.Size
	}
	return message
}

//OnValueParsed post processing for field
func (parser *BaseLocationReportParser) OnValueParsed(confItems []Item, itemName string, value interface{}) {

}

//InvalidateFields makes fields human readable
func (parser *BaseLocationReportParser) InvalidateFields(config *Report, messages *[]report.IMessage) {
	if len(*messages) > 0 {
		parser.SetAck((*messages)[0])
	}
	for _, message := range *messages {
		parser.HumanizeUniqueID(message.(IQueclinkMessage), fields.UniqueID)
		parser.HumanizeSpeed(message.(IQueclinkMessage))
		parser.HumanizeDateTime(message.(IQueclinkMessage), fields.GPSUTCTime)
		parser.HumanizeDateTime(message.(IQueclinkMessage), fields.SendTime)
		parser.HumanizeCurrentMileage(message.(IQueclinkMessage))
		parser.HumanizeTotalMileage(message.(IQueclinkMessage))
		parser.HumanizeCurrentHourMeterCount(message.(IQueclinkMessage))
		parser.HumanizeTotalHourMeterCount(message.(IQueclinkMessage))
		parser.HumanizeTimestamp(message.(IQueclinkMessage))
		parser.HumanizeLngLat(message.(IQueclinkMessage), fields.Latitude)
		parser.HumanizeLngLat(message.(IQueclinkMessage), fields.Longitude)
		parser.HumanizeInputs(message.(IQueclinkMessage))
		parser.HumanizeValidity(config, message.(IQueclinkMessage))
		parser.HumanizeInt32ToFloat32(fields.Altitude, message.(IQueclinkMessage))
		parser.HumanizeInt32ToFloat32(fields.Heading, message.(IQueclinkMessage))
		parser.HumanizeSatellites(message.(IQueclinkMessage))
		parser.HumanizeBatteryPercentage(message.(IQueclinkMessage))
		parser.DecodeReason(message.(IQueclinkMessage))
		message.SetValue(fields.LocationMessage, config.Location)
	}
}

//HumanizeUniqueID converts byte set to string uniqueID
func (parser *BaseLocationReportParser) HumanizeUniqueID(message IQueclinkMessage, fn string) {
	var iID interface{}
	var v bool
	if iID, v = message.GetValue(fn); !v {
		return
	}

	bID, v := iID.([]byte)
	if !v {
		return
	}
	imei := ""

	for i := 0; i < 7; i++ {
		imei = fmt.Sprintf("%v%v", imei, fmt.Sprintf("%02d", bID[i]))
	}
	imei = fmt.Sprintf("%v%v", imei, bID[7])
	message.SetValue(fn, imei)
	message.SetValue(fields.DevID, fmt.Sprintf("%v_%v", "queclink", imei))
}

//HumanizeSpeed converts speed from byte array to float32
func (parser *BaseLocationReportParser) HumanizeSpeed(message IQueclinkMessage) {
	var iSpeed interface{}
	var v bool
	if iSpeed, v = message.GetValue(fields.Speed); !v {
		return
	}

	bSpeed, v := iSpeed.([]byte)
	if !v {
		return
	}

	intSpeed, _ := utils.GetBytesValue(bSpeed, 0, 2)

	speed := float32(intSpeed.(int32)) + float32(bSpeed[2])/10.0
	speed = speed / 3.6
	speed = float32(math.Round(float64(speed*100.0)) / 100.0)
	message.SetValue(fields.Speed, speed)
}

//HumanizeDateTime converts byte set to timestamp
func (parser *BaseLocationReportParser) HumanizeDateTime(message IQueclinkMessage, fn string) {
	var iTs interface{}
	var v bool
	if iTs, v = message.GetValue(fn); !v {
		return
	}
	bTs, v := iTs.([]byte)
	if !v || len(bTs) < 7 {
		return
	}
	y, yE := utils.GetBytesValue(bTs, 0, 2)
	m, mE := utils.GetBytesValue(bTs, 2, 1)
	d, dE := utils.GetBytesValue(bTs, 3, 1)
	h, hE := utils.GetBytesValue(bTs, 4, 1)
	mi, miE := utils.GetBytesValue(bTs, 5, 1)
	s, sE := utils.GetBytesValue(bTs, 6, 1)
	sts := "1989-01-01T00:00:00Z"
	if yE == nil && mE == nil && dE == nil && hE == nil && miE == nil && sE == nil {
		sts = fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", y.(int32), m.(byte), d.(byte), h.(byte), mi.(byte), s.(byte))
	}
	datetime, terr := time.Parse("2006-01-02T15:04:05Z", sts) //parse using yyyy-MM-ddTHH:mm:ssZ
	if terr != nil {
		datetime, _ = time.Parse("2006-01-02T15:04:05Z", "1989-01-01T00:00:00Z")
	}
	message.SetValue(fn, &utils.JSONTime{Time: datetime})
}

//HumanizeCurrentMileage makes field readable
func (parser *BaseLocationReportParser) HumanizeCurrentMileage(message IQueclinkMessage) {

	var iCm interface{}
	var v bool
	if iCm, v = message.GetValue(fields.CurrentMileage); !v {
		return
	}

	bCm, v := iCm.([]byte)
	if !v {
		return
	}

	intCm, _ := utils.GetBytesValue(bCm, 0, 2)

	cm := float32(intCm.(int32)) + float32(bCm[2])/10.0
	message.SetValue(fields.CurrentMileage, float32(cm))
}

//HumanizeTotalMileage makes field readable
func (parser *BaseLocationReportParser) HumanizeTotalMileage(message IQueclinkMessage) {
	var iCm interface{}
	var v bool
	if iCm, v = message.GetValue(fields.Odometer); !v {
		return
	}

	bCm, v := iCm.([]byte)
	if !v {
		return
	}

	intCm, _ := utils.GetBytesValue(bCm, 0, 4)

	cm := float32(intCm.(int32)) + float32(bCm[4])/10.0
	message.SetValue(fields.Odometer, int32(float32(cm)*1000))
}

//HumanizeCurrentHourMeterCount makes field readable in seconds
func (parser *BaseLocationReportParser) HumanizeCurrentHourMeterCount(message IQueclinkMessage) {
	var iCm interface{}
	var v bool
	if iCm, v = message.GetValue(fields.CurrentHourMeterCount); !v {
		return
	}

	bCm, v := iCm.([]byte)
	if !v || len(bCm) < 3 {
		return
	}

	duration := int32(bCm[0])*3600 + int32(bCm[1])*60 + int32(bCm[2])

	message.SetValue(fields.CurrentHourMeterCount, duration)
}

//HumanizeTotalHourMeterCount makes field readable in seconds
func (parser *BaseLocationReportParser) HumanizeTotalHourMeterCount(message IQueclinkMessage) {
	var iCm interface{}
	var v bool
	if iCm, v = message.GetValue(fields.TotalHourMeterCount); !v {
		return
	}

	bCm, v := iCm.([]byte)
	if !v || len(bCm) < 3 {
		return
	}
	ih, _ := utils.GetBytesValue(bCm, 0, 4)
	hh := ih.(int32)
	hhhh := int64(hh) * 3600
	duration := hhhh + int64(bCm[4])*60 + int64(bCm[5])

	message.SetValue(fields.TotalHourMeterCount, duration)
}

//HumanizeTimestamp makes field readable in seconds
func (parser *BaseLocationReportParser) HumanizeTimestamp(message IQueclinkMessage) {

	if message.MessageType() == "+BVT" || message.MessageType() == "+BSP" {
		if dt, found := message.GetValue(fields.GPSUTCTime); found {
			message.SetValue(fields.TimeStamp, dt)
		}
	} else {
		if dt, found := message.GetValue(fields.SendTime); found {
			message.SetValue(fields.TimeStamp, dt)
		} else {
			if dt, found := message.GetValue(fields.GPSUTCTime); found {
				message.SetValue(fields.TimeStamp, dt)
			}
		}
	}
}

//HumanizeLngLat makes field readable
func (parser *BaseLocationReportParser) HumanizeLngLat(message IQueclinkMessage, fn string) {
	var iL interface{}
	var v bool
	if iL, v = message.GetValue(fn); !v {
		return
	}

	var intL int32

	switch iL.(type) {
	case []byte:
		{
			bL := iL.([]byte)
			iL, _ = utils.GetBytesValue(bL, 0, len(bL))
			break
		}
	}

	intL, v = iL.(int32)
	if !v {
		return
	}

	fl64 := float64(intL)

	l := float32(fl64 / 1000000.0)
	message.SetValue(fn, l)
}

//HumanizeInputs generates gpio and IgnitionStatus
func (parser *BaseLocationReportParser) HumanizeInputs(message IQueclinkMessage) {
	ignitionStatus := byte(0)
	iS, found := message.GetValue(fields.IgnitionStatus)
	if found {
		ignitionStatus = iS.(byte)
	}

	if iDis, f := message.GetValue(fields.DigitalInputStatus); f {
		value := iDis.(byte)
		message.SetValue(fields.GPIO, byte(value>>1))

		if ((value & 1) == 1) || ignitionStatus == 32 {
			message.SetValue(fields.IgnitionState, byte(1))
		} else {
			message.SetValue(fields.IgnitionState, byte(0))
		}
	} else {
		message.SetValue(fields.GPIO, byte(0))
		if ignitionStatus == 32 {
			message.SetValue(fields.IgnitionState, byte(1))
		} else {
			message.SetValue(fields.IgnitionState, byte(0))
		}
	}
}

//HumanizeValidity validator for validity
func (parser *BaseLocationReportParser) HumanizeValidity(config *Report, message IQueclinkMessage) {
	if _, found := message.GetValue(fields.GPSAccuracy); found {
		if parser.IsValid(config, message) {
			message.SetValue(fields.GpsValidity, byte(1))
		} else {
			message.SetValue(fields.GpsValidity, byte(0))
		}
	}
}

//IsValid checks message is valid
func (parser *BaseLocationReportParser) IsValid(config *Report, message IQueclinkMessage) bool {
	lastKnownPosition := false
	if messageType, found := config.GetType(message.EventCode()); found {
		lastKnownPosition = messageType.LastKnownPosition
	}

	if lastKnownPosition {
		return parser.GetPassiveValidity(message)
	}
	return parser.GetValidity(message)
}

//GetPassiveValidity returns alternative GPS validity
func (parser *BaseLocationReportParser) GetPassiveValidity(message IQueclinkMessage) bool {

	gpsValidity := false
	if a, f := message.GetValue(fields.GPSAccuracy); f {
		gpsValidity = a.(byte) != 0
	}

	tsValidity := false
	if t, f := message.GetValue(fields.TimeStamp); f {
		tsValidity = t.(*utils.JSONTime).Before(time.Now().UTC().AddDate(0, 0, 1))
	}
	if gpsValidity && tsValidity {
		return true
	}

	if isat, f := message.GetValue(fields.Satellites); f {
		sat := isat.(byte)
		lat, _ := message.GetValue(fields.Latitude)
		lon, _ := message.GetValue(fields.Longitude)
		return sat != 0 && lat.(float32) != 0 && lon.(float32) != 0 && tsValidity
	}
	return false
}

//GetValidity returns GPS validity
func (parser *BaseLocationReportParser) GetValidity(message IQueclinkMessage) bool {
	gpsValidity := false
	if a, f := message.GetValue(fields.GPSAccuracy); f {
		gpsValidity = a.(byte) != 0
	}

	tsValidity := false
	if t, f := message.GetValue(fields.TimeStamp); f {
		tsValidity = t.(*utils.JSONTime).Before(time.Now().UTC().AddDate(0, 0, 1))
	}
	return gpsValidity && tsValidity
}

//HumanizeInt32ToFloat32 converts int32 to float32 value
func (parser *BaseLocationReportParser) HumanizeInt32ToFloat32(fieldName string, message IQueclinkMessage) {
	if iv, f := message.GetValue(fieldName); f {
		if v, valid := iv.(int32); valid {
			message.SetValue(fieldName, float32(v))
		}
	}
}

//HumanizeSatellites converts satellites
func (parser *BaseLocationReportParser) HumanizeSatellites(message IQueclinkMessage) {
	if iv, f := message.GetValue(fields.Satellites); f {
		if v, valid := iv.(byte); valid {
			message.SetValue(fields.Satellites, int32(v))
		}
	}
}

//HumanizeBatteryPercentage converts BatteryPercentage
func (parser *BaseLocationReportParser) HumanizeBatteryPercentage(message IQueclinkMessage) {
	if iv, f := message.GetValue(fields.BatteryPercentage); f {
		if v, valid := iv.(byte); valid {
			message.SetValue(fields.BatteryPercentage, float32(v))
		}
	}
}

//DecodeReason for message
func (parser *BaseLocationReportParser) DecodeReason(message IQueclinkMessage) {
	message.SetValue("Reason", GetReason(message))

}

func (parser *BaseLocationReportParser) readMaskedValue(packet []byte, mask int32, position byte, size int, readFrom *int) interface{} {
	if utils.BitIsSet(int64(mask), uint(position)) {
		return parser.readValue(packet, size, readFrom)
	}
	return nil
}

func (parser *BaseLocationReportParser) readValue(packet []byte, size int, readFrom *int) interface{} {
	value, _ := utils.GetBytesValue(packet, *readFrom, size)
	*readFrom = *readFrom + size
	return value
}

func (parser *BaseLocationReportParser) addValue(packet []byte, size int, fieldname string, message report.IMessage, readFrom *int) interface{} {
	value, _ := utils.GetBytesValue(packet, *readFrom, size)
	message.SetValue(fieldname, value)
	*readFrom = *readFrom + size
	return value
}

func (parser *BaseLocationReportParser) addMaskedValue(packet []byte, mask int32, position byte, size int, fieldname string, message report.IMessage, readFrom *int) interface{} {
	if utils.BitIsSet(int64(mask), uint(position)) {
		return parser.addValue(packet, size, fieldname, message, readFrom)
	}
	return nil
}
