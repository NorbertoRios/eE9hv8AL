package queclinkreport

import (
	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
	"queclink-go/queclinkreport/fields"
)

//GV300WEriParser eri report parser for gv300w devices
type GV300WEriParser struct {
	BaseLocationReportParser
}

//Parse eri report
func (parser *GV300WEriParser) Parse(packet []byte, config *Report) []report.IMessage {
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
func (parser *GV300WEriParser) ParseSensors(message IQueclinkMessage, packet []byte, readFrom *int) {
	uartDeviceType := byte(0)
	if v, found := message.GetValue(fields.UARTDeviceType); found {
		uartDeviceType = v.(byte)
	}
	switch uartDeviceType {
	case 1: //DiginFuelSensor
		{
			parser.addFuelValue(packet, readFrom, message)
			break
		}
	default:
		*readFrom = *readFrom + 2
		break
	}
	parser.addValue(packet, 1, fields.Number, message, readFrom)
}

func (parser *GV300WEriParser) addValue(packet []byte, size int, fieldname string, message report.IMessage, readFrom *int) interface{} {
	value, _ := utils.GetBytesValue(packet, *readFrom, size)
	message.SetValue(fieldname, value)
	*readFrom = *readFrom + size
	return value
}

func (parser *GV300WEriParser) addFuelValue(packet []byte, readFrom *int, message report.IMessage) {
	value, _ := utils.GetBytesValue(packet, *readFrom, 2)
	fValue := float32(float32(value.(int32)) / 10)
	message.SetValue(fields.FuelLevel, fValue)
	*readFrom = *readFrom + 2
}
