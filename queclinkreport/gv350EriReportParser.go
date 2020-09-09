package queclinkreport

import (
	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
	"queclink-go/queclinkreport/fields"
)

//GV350EriReportParser parser for eri reports
type GV350EriReportParser struct {
	GV600EriReportParser
}

const (
	fuelPercantageFlag = 3
	fuelVolumeFlag     = 4
)

//Parse eri report
func (parser *GV350EriReportParser) Parse(packet []byte, config *Report) []report.IMessage {
	parser.Init(&packet)
	messages := []report.IMessage{}
	readFrom := parser.StartByte
	header := parser.ParseHeader(parser.ReportMask, packet, config.Header, &readFrom)
	header.SetValue(fields.MessageHeader, parser.MessageHeader)
	header.SetValue(fields.MessageType, parser.MessageType)
	header.SetValue(fields.ReportMask, parser.ReportMask)
	header.SetValue("RawData", packet)

	parser.ParseSensors(header, packet, &readFrom)

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
		parser.PopulateTemperature(report)
	}

	if len(messages) == 0 && (len(*header.GetData()) != 0 || len(*tail.GetData()) != 0) {
		header.AppendRange(*tail.GetData())
		messages = append(messages, header)
	}
	parser.InvalidateFields(config, &messages)
	return messages
}

//ParseSensors parses values of sensor data
func (parser *GV350EriReportParser) ParseSensors(message IQueclinkMessage, packet []byte, readFrom *int) {
	mask := int32(0)
	if v, found := message.GetValue(fields.ERIMask); found {
		mask = v.(int32)
	}
	if utils.BitIsSet(int64(mask), uint(temperatureFlag)) {
		parser.parse1WireSensorData(packet, message, readFrom)
	}

	if utils.BitIsSet(int64(mask), uint(fuelPercantageFlag)) || utils.BitIsSet(int64(mask), uint(fuelVolumeFlag)) {
		parser.parseFuelData(packet, message, mask, readFrom)
	}

	parser.addValue(packet, 1, fields.Number, message, readFrom)
}

func (parser *GV350EriReportParser) parseFuelData(packet []byte, message report.IMessage, mask int32, readFrom *int) {
	sensorNumber := 0
	if utils.BitIsSet(int64(mask), uint(3)) || utils.BitIsSet(int64(mask), uint(4)) {
		iSN := parser.addValue(packet, 1, fields.SensorNumber, message, readFrom)
		bSN, _ := iSN.(byte)
		sensorNumber = int(bSN)
	}

	var sensors = []FuelSensor{}

	for i := 1; i <= sensorNumber; i++ {
		sensor := FuelSensor{}
		if value := parser.readValue(packet, 1, readFrom); value != nil {
			if bValue, valid := value.(byte); valid {
				sensor.SensorType = int32(bValue)
			}
		}

		*readFrom = *readFrom + 1 //strange description in documentation... reserved with length 0)
		if value := parser.readMaskedValue(packet, mask, 3, 1, readFrom); value != nil {
			if bValue, valid := value.(int32); valid {
				sensor.Percentage = bValue
			}
		}

		if value := parser.readMaskedValue(packet, mask, 4, 2, readFrom); value != nil {
			if bValue, valid := value.(int32); valid {
				sensor.Volume = bValue
			}
		}
		sensors = append(sensors, sensor)
	}
	if len(sensors) > 0 {
		message.SetValue(fields.FuelSensors, sensors)
	}
}

//InvalidateFields makes fields human readable
func (parser *GV350EriReportParser) InvalidateFields(config *Report, messages *[]report.IMessage) {
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

//HumanizeSatellites converts satellites
func (parser *GV350EriReportParser) HumanizeSatellites(message IQueclinkMessage) {
	if iv, f := message.GetValue(fields.Satellites); f {
		if v, valid := iv.(byte); valid {
			message.SetValue(fields.Satellites, int32(v&0x0F))
			message.SetValue("ExternalGPSAntennaStatus", int32((v>>6)&0x03))
		}
	}
}
