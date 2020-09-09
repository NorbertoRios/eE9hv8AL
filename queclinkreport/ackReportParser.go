package queclinkreport

import (
	"queclink-go/base.device.service/report"
)

//AckReportParser for ack messages
type AckReportParser struct {
	BaseLocationReportParser
}

//Init crash parser parser
func (parser *AckReportParser) Init(packet *[]byte) {
	parser.MessageHeader = string((*packet)[0:4])
	parser.MessageType = int32((*packet)[4])
	parser.ReportMask = int32((*packet)[5])
	parser.StartByte = 6
}

//Parse Ack packet
func (parser *AckReportParser) Parse(packet []byte, config *Report) []report.IMessage {
	parser.Init(&packet)
	messages := parser.BaseLocationReportParser.Parse(packet, config)
	parser.InvalidateFields(config, &messages)
	return messages
}
