package queclinkreport

import (
	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
	"queclink-go/queclinkreport/fields"
)

//EriReportParser parser for eri reports
type EriReportParser struct {
	BaseLocationReportParser
}

const (
	digitFuelSensorDataFlag = 0
)

//Parse eri report
func (parser *EriReportParser) Parse(packet []byte, config *Report) []report.IMessage {
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
	}

	if len(messages) == 0 && (len(*header.GetData()) != 0 || len(*tail.GetData()) != 0) {
		header.AppendRange(*tail.GetData())
		messages = append(messages, header)
	}
	parser.InvalidateFields(config, &messages)
	return messages
}

//ParseSensors parses values of sensor data
func (parser *EriReportParser) ParseSensors(message IQueclinkMessage, packet []byte, readFrom *int) {
	mask := int32(0)
	if v, found := message.GetValue(fields.ERIMask); found {
		mask = v.(int32)
	}
	parser.addMaskedValue(packet, mask, digitFuelSensorDataFlag, 2, fields.FuelLevel, message, readFrom)

	uartDeviceType := byte(0)
	if v, found := message.GetValue(fields.UARTDeviceType); found {
		uartDeviceType = v.(byte)
	}
	switch uartDeviceType {
	case 1: //DiginFuelSensor
		{
			parser.parseFuelData(packet, message, mask, readFrom)
			break
		}
	case 2: //Ac100_1Wire_Bus
		{
			parser.parseAc100SensorData(packet, message, mask, readFrom)
			break
		}
	}
	//main uart device
	parser.parseMainDeviceUart(packet, message, mask, readFrom)
	parser.parsePulseFrequency(packet, message, mask, readFrom)
	parser.addValue(packet, 1, fields.Number, message, readFrom)
}

func (parser *EriReportParser) parseMainDeviceUart(packet []byte, message report.IMessage, mask int32, readFrom *int) {
	imainUartDevice := parser.addValue(packet, 1, fields.MainUARTDevice, message, readFrom)
	mainUartDevice := 0
	if v, valid := imainUartDevice.(byte); valid {
		mainUartDevice = int(v)
	}
	switch mainUartDevice {
	case 1: //Ac200_1Wire_Bus
		{
			parser.parseAc200SensorData(packet, message, mask, readFrom)
			break
		}
	}
}

func (parser *EriReportParser) parseAc200SensorData(packet []byte, message report.IMessage, mask int32, readFrom *int) {
	iSn := parser.addMaskedValue(packet, mask, 2, 2, fields.AC200DeviceNumber, message, readFrom)
	sensorNumber := 0
	if bSn, valid := iSn.(int32); valid {
		sensorNumber = int(bSn)
	}

	var sensors = []AcSensor{}
	for i := 1; i <= sensorNumber; i++ {
		sensor := parser.parseAcSensor(packet, mask, readFrom)
		sensors = append(sensors, sensor)
	}

	if len(sensors) > 0 {
		message.SetValue(fields.AC200Sensors, sensors)
	}
}

func (parser *EriReportParser) parseAcSensor(packet []byte, mask int32, readFrom *int) AcSensor {
	sensor := AcSensor{}

	if value := parser.readValue(packet, 8, readFrom); value != nil {
		if sValue, valid := value.([]byte); valid {
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
				sensor.OneWireDeviceData = float32(value.(int32)) * 0.0625
			}
		}
	}
	return sensor
}

func (parser *EriReportParser) parsePulseFrequency(packet []byte, message report.IMessage, mask int32, readFrom *int) {
	if utils.BitIsSet(int64(mask), 6) {
		parser.addValue(packet, 4, "PulseFrequency", message, readFrom)
	}
}

func (parser *EriReportParser) parseAc100SensorData(packet []byte, message report.IMessage, mask int32, readFrom *int) {
	iSn := parser.addMaskedValue(packet, mask, 1, 1, fields.AC100DeviceNumber, message, readFrom)
	sensorNumber := 0
	if bSn, valid := iSn.(byte); valid {
		sensorNumber = int(bSn)
	}
	sensors := []AcSensor{}
	for i := 1; i <= sensorNumber; i++ {
		sensor := parser.parseAcSensor(packet, mask, readFrom)
		sensors = append(sensors, sensor)
	}
	if len(sensors) > 0 {
		message.SetValue(fields.AC100Sensors, sensors)
	}
}

func (parser *EriReportParser) parseFuelData(packet []byte, message report.IMessage, mask int32, readFrom *int) {
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
		if value := parser.readMaskedValue(packet, mask, 3, 2, readFrom); value != nil {
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
