package queclinkreport

import (
	"queclink-go/base.device.service/report"
	"queclink-go/queclinkreport/fields"
)

//GV350LBCReportParser location by call message parser
type GV350LBCReportParser struct {
	GV350LocationReportParser
	LBCReportParser
}

//Parse LBC packet
func (parser *GV350LBCReportParser) Parse(packet []byte, config *Report) []report.IMessage {
	parser.LBCReportParser.Init(&packet)
	messages := []report.IMessage{}
	readFrom := parser.LBCReportParser.StartByte
	header := parser.LBCReportParser.ParseHeader(parser.LBCReportParser.ReportMask, packet, config.Header, &readFrom)
	header.SetValue(fields.MessageHeader, parser.LBCReportParser.MessageHeader)
	header.SetValue(fields.MessageType, parser.LBCReportParser.MessageType)
	header.SetValue(fields.ReportMask, parser.LBCReportParser.ReportMask)
	header.SetValue("RawData", packet)

	number := byte(0)
	if v, found := header.GetValue(fields.Number); found {
		number = v.(byte)
	}

	for i := 1; i <= int(number); i++ {
		body := parser.LBCReportParser.ParseBody(parser.LBCReportParser.ReportMask, packet, config.MultiPosition, &readFrom)
		messages = append(messages, body)
	}
	var tail = parser.LBCReportParser.ParseTail(parser.LBCReportParser.ReportMask, packet, config.Tail, &readFrom)

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

//InvalidateFields makes fields human readable
func (parser *GV350LBCReportParser) InvalidateFields(config *Report, messages *[]report.IMessage) {
	parser.GV350LocationReportParser.InvalidateFields(config, messages)
	for _, message := range *messages {
		parser.HumanizePhone(message.(IQueclinkMessage))
	}
}
