package queclinkreport

import (
	"fmt"

	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
	"queclink-go/queclinkreport/fields"
)

//InfReportParser parses inf packets
type InfReportParser struct {
	BaseLocationReportParser
}

//Init crash parser parser
func (parser *InfReportParser) Init(packet *[]byte) {
	parser.MessageHeader = string((*packet)[0:4])
	parser.MessageType = int32((*packet)[4])
	ireportMask, _ := utils.GetBytesValue(*packet, 5, 2)
	parser.ReportMask = ireportMask.(int32)
	parser.StartByte = 7
}

//Parse INF packet
func (parser *InfReportParser) Parse(packet []byte, config *Report) []report.IMessage {
	parser.Init(&packet)
	messages := parser.BaseLocationReportParser.Parse(packet, config)
	parser.InvalidateFields(config, &messages)
	return messages
}

//InvalidateFields makes fields human readable
func (parser *InfReportParser) InvalidateFields(config *Report, messages *[]report.IMessage) {
	parser.BaseLocationReportParser.InvalidateFields(config, messages)
	for _, message := range *messages {
		parser.humanizeExternalPowerSupply(message)
		parser.humanizeGirType(message.(IQueclinkMessage))
		parser.humanizeBytesToHex(message.(IQueclinkMessage), "HardwareVersion")
		parser.humanizeBytesToHex(message.(IQueclinkMessage), "MCUVersion")
		parser.humanizeBytesToHex(message.(IQueclinkMessage), "FirmwareVersion")
		parser.humanizeBytesToHex(message.(IQueclinkMessage), "ProtocolVersion")
		parser.humanizeBytesToHex(message.(IQueclinkMessage), "ICCID")
		parser.HumanizeDateTime(message.(IQueclinkMessage), "LastFixUTCTime")
		parser.humanizeParameters(message.(IQueclinkMessage))
	}
}

func (parser *InfReportParser) humanizeExternalPowerSupply(message report.IMessage) {
	ivalue, found := message.GetValue(fields.Supply)
	if !found {
		return
	}
	value := ivalue.(int32)
	message.SetValue("Backup Battery Charge Mode", value&1)
	message.SetValue("LED On", (value>>4)&1)
	message.SetValue("Charging", (value>>5)&1)
	message.SetValue("Backup Battery On", (value>>6)&1)
	message.SetValue("Main Supply", (value>>7)&1)
}

func (parser *InfReportParser) humanizeGirType(message report.IMessage) {
	ivalue, found := message.GetValue("GIRTriggerType")
	if !found {
		return
	}

	girType := fmt.Sprintf("%v", ivalue)
	switch ivalue.(byte) {
	case 1:
		girType = "This cell information is for SOS requirement"
		break
	case 2:
		girType = "This cell information is for RTL requirement"
		break
	case 3:
		girType = "This cell information is for LBC requirement"
		break
	case 5:
		girType = "This cell information is for FRI requirement"
		break
	case 6:
		girType = "This cell information is for GIR requirement"
		break
	}
	message.SetValue("GIRTriggerType", girType)
}

func (parser *InfReportParser) humanizeBytesToHex(message report.IMessage, fn string) {
	ivalue, found := message.GetValue(fn)
	if !found {
		return
	}

	hexData := fmt.Sprintf("%2X", ivalue)
	message.SetValue(fn, hexData)
}

func (parser *InfReportParser) humanizeParameters(message report.IMessage) {
	ivalue, found := message.GetValue("PowerOWHAGPS")
	if !found {
		return
	}

	param := ivalue.(byte)
	message.SetValue("AGPS", param&1)
	message.SetValue("Outside working Hour", (param>>2)&1)
	message.SetValue("OWH Mode", (param>>3)&3)
	powerSavingEnable := (param >> 5) & 3
	switch powerSavingEnable {
	case 0:
		{
			message.SetValue("Power saving mode", "Disable power saving function")
			break
		}
	case 1:
		{
			message.SetValue("Power saving mode", "Deep saving mode")
			break
		}
	case 2:
		{
			message.SetValue("Power saving mode", "Low saving mode")
			break
		}
	}
}
