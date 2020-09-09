package queclinkreport

import (
	"fmt"

	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
	"queclink-go/queclinkreport/fields"
)

//CrashReportParser parses crash reports
type CrashReportParser struct {
	BaseLocationReportParser
}

//Init crash parser parser
func (parser *CrashReportParser) Init(packet *[]byte) {
	parser.MessageHeader = string((*packet)[0:4])
	parser.MessageType = 0
	ireportMask, _ := utils.GetBytesValue(*packet, 4, 2)
	parser.ReportMask = ireportMask.(int32)
	parser.StartByte = 6
}

//Parse crash packet
func (parser *CrashReportParser) Parse(packet []byte, config *Report) []report.IMessage {
	parser.Init(&packet)
	messages := parser.BaseLocationReportParser.Parse(packet, config)
	parser.InvalidateFields(config, &messages)
	return messages
}

//InvalidateFields makes fields human readable
func (parser *CrashReportParser) InvalidateFields(config *Report, messages *[]report.IMessage) {
	parser.BaseLocationReportParser.InvalidateFields(config, messages)
	for _, message := range *messages {
		parser.HumanizeData(message.(IQueclinkMessage))
	}
}

//HumanizeData converts bytes to hex
func (parser *CrashReportParser) HumanizeData(message report.IMessage) {
	idata, found := message.GetValue(fields.Data)
	if !found {
		return
	}

	dataPair := idata.([]byte)
	hexData := fmt.Sprintf("%2X", dataPair)
	message.SetValue(fields.Data, hexData)

}
