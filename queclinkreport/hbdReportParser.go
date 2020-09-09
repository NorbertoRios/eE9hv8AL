package queclinkreport

import "queclink-go/base.device.service/report"

//HBDReportParser parser for heartbeat reports
type HBDReportParser struct {
	BaseLocationReportParser
}

//Init crash parser parser
func (parser *HBDReportParser) Init(packet *[]byte) {
	parser.MessageHeader = string((*packet)[0:4])
	parser.MessageType = 0
	ireportMask := (*packet)[4]
	parser.ReportMask = int32(ireportMask)
	parser.StartByte = 5
}

//Parse HBD packet
func (parser *HBDReportParser) Parse(packet []byte, config *Report) []report.IMessage {
	parser.Init(&packet)
	messages := parser.BaseLocationReportParser.Parse(packet, config)
	parser.InvalidateFields(config, &messages)
	return messages
}
