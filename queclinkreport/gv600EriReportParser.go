package queclinkreport

import (
	"math"

	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
	"queclink-go/queclinkreport/fields"
)

//GV600EriReportParser parser for eri reports
type GV600EriReportParser struct {
	BaseLocationReportParser
}

const (
	temperatureFlag = 1
)

//Parse eri report
func (parser *GV600EriReportParser) Parse(packet []byte, config *Report) []report.IMessage {
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

//PopulateTemperature converts temperature sensors
func (parser *GV600EriReportParser) PopulateTemperature(message report.IMessage) {
	qmessage, v := message.(*QueclinkMessage)
	if !v {
		return
	}

	if its, found := qmessage.GetValue(fields.OWireSensors); found && its != nil {
		if ts, valid := its.([]AcSensor); valid && len(ts) > 0 {
			qmessage.TemperatureSensors = &report.TemperatureSensors{}
			if len(ts) >= 1 {
				qmessage.TemperatureSensors.Sensor1 = parser.getTemperatureSensor(ts, 0)
			} else {
				qmessage.TemperatureSensors.Sensor1 = nil
			}
			if len(ts) >= 2 {
				qmessage.TemperatureSensors.Sensor2 = parser.getTemperatureSensor(ts, 1)
			} else {
				qmessage.TemperatureSensors.Sensor2 = nil
			}
			if len(ts) >= 3 {
				qmessage.TemperatureSensors.Sensor3 = parser.getTemperatureSensor(ts, 2)
			} else {
				qmessage.TemperatureSensors.Sensor3 = nil
			}
			if len(ts) >= 4 {
				qmessage.TemperatureSensors.Sensor4 = parser.getTemperatureSensor(ts, 3)
			} else {
				qmessage.TemperatureSensors.Sensor4 = nil
			}
		}
	}
}

func (parser *GV600EriReportParser) getTemperatureSensor(owireSensors []AcSensor, index int) *report.TemperatureSensor {
	return &report.TemperatureSensor{
		Id:               owireSensors[index].OWireDeviceID,
		TemperatureValue: owireSensors[index].OneWireDeviceData,
	}
}

//ParseSensors parses values of sensor data
func (parser *GV600EriReportParser) ParseSensors(message IQueclinkMessage, packet []byte, readFrom *int) {
	mask := int32(0)
	if v, found := message.GetValue(fields.ERIMask); found {
		mask = v.(int32)
	}
	if utils.BitIsSet(int64(mask), uint(temperatureFlag)) {
		parser.parse1WireSensorData(packet, message, readFrom)
	}

	parser.addValue(packet, 1, fields.Number, message, readFrom)
}

func (parser *GV600EriReportParser) parse1WireSensorData(packet []byte, message report.IMessage, readFrom *int) {
	iSn := parser.addValue(packet, 1, fields.OWireDeviceNumber, message, readFrom)
	sensorNumber := 0
	if bSn, valid := iSn.(byte); valid {
		sensorNumber = int(bSn)
	}

	var sensors = []AcSensor{}
	for i := 1; i <= sensorNumber; i++ {
		sensor := parser.parseAcSensor(packet, readFrom)
		sensors = append(sensors, sensor)
	}

	if len(sensors) > 0 {
		message.SetValue(fields.OWireSensors, sensors)
	}
}

func (parser *GV600EriReportParser) parseAcSensor(packet []byte, readFrom *int) AcSensor {
	sensor := AcSensor{}

	if value := parser.readValue(packet, 8, readFrom); value != nil {
		if sValue, valid := value.([]byte); valid {
			utils.Reverse(sValue)
			sensor.OWireDeviceID = utils.ByteToString(sValue)
		}
	}

	if value := parser.readValue(packet, 1, readFrom); value != nil {
		if bValue, valid := value.(byte); valid {
			sensor.OWireDeviceType = int32(bValue)
		}
	}

	if value := parser.readValue(packet, 1, readFrom); value != nil {
		if bValue, valid := value.(byte); valid {
			sensor.DeviceDataLength = int32(bValue)
		}
	}
	if value := parser.readValue(packet, int(sensor.DeviceDataLength), readFrom); value != nil {
		switch value.(type) {
		case byte:
			{
				sensor.OneWireDeviceData = float32(value.(byte)) * 0.0625
			}
		case int32:
			{
				tdata := value.(int32)
				if (tdata & 0x1F00) == 0x1F00 {
					tdata = int32(int16(tdata) + 1)
				}
				sensor.OneWireDeviceData = float32(tdata) * 0.0625
			}
		}
	}
	sensor.OneWireDeviceData = float32(math.Round(float64(sensor.OneWireDeviceData)*100.0) / 100)
	return sensor
}
