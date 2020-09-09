package queclinkreport

import (
	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
	"queclink-go/queclinkreport/fields"
)

//GV600DatReportParser parser for event reports GTSTC and GTBTC
type GV600DatReportParser struct {
	BaseLocationReportParser
}

//Parse crash packet
func (parser *GV600DatReportParser) Parse(packet []byte, config *Report) []report.IMessage {
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
	return messages
}

//ParseHeader of message
func (parser *GV600DatReportParser) ParseHeader(reportMask int32, packet []byte, confItems []Item, readFrom *int) IQueclinkMessage {
	message := NewQueclinkMessage()
	for _, item := range confItems {
		if item.MaskBit != nil && !utils.BitIsSet(int64(reportMask), uint(*item.MaskBit)) {
			continue
		}
		size := item.Size
		var value interface{}

		if item.SizeIn != "" {
			if v, f := message.GetValue(item.SizeIn); f {
				switch v.(type) {
				case int32:
					size = int(v.(int32))
				case byte:
					size = int(v.(byte))
				}
			}

		}
		value, _ = utils.GetBytesValue(packet, *readFrom, size)

		parser.OnValueParsed(confItems, item.ItemName, value)
		message.SetValue(item.ItemName, value)
		*readFrom = *readFrom + size
	}
	return message
}
