package queclinkreport

import "queclink-go/base.device.service/report"

//LocationReportParser parses inf packets
type LocationReportParser struct {
	BaseLocationReportParser
}

//Parse location reports packet
func (parser *LocationReportParser) Parse(packet []byte, config *Report) []report.IMessage {
	parser.Init(&packet)
	messages := parser.BaseLocationReportParser.Parse(packet, config)
	parser.InvalidateFields(config, &messages)
	return messages
}
