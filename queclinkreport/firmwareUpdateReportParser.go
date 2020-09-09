package queclinkreport

import (
	"queclink-go/base.device.service/report"
	"queclink-go/queclinkreport/fields"
)

//FirmwareUpdateReportParser parser for event reports GTSTC and GTBTC
type FirmwareUpdateReportParser struct {
	BaseLocationReportParser
}

//Parse crash packet
func (parser *FirmwareUpdateReportParser) Parse(packet []byte, config *Report) []report.IMessage {
	parser.Init(&packet)
	messages := parser.BaseLocationReportParser.Parse(packet, config)
	parser.InvalidateFields(config, &messages)
	return messages
}

//InvalidateFields makes fields human readable
func (parser *FirmwareUpdateReportParser) InvalidateFields(config *Report, messages *[]report.IMessage) {
	parser.BaseLocationReportParser.InvalidateFields(config, messages)
	for _, message := range *messages {
		parser.HumanizeDescription(message.(IQueclinkMessage))
	}
}

//HumanizeDescription translate code to text
func (parser *FirmwareUpdateReportParser) HumanizeDescription(message IQueclinkMessage) {
	icode, found := message.GetValue(fields.Code)
	if !found {
		return
	}
	code, valid := icode.(int32)
	if !valid {
		return
	}

	if reason := parser.getReason(code); reason != "" {
		message.SetValue(fields.Description, reason)
	}
}

func (parser *FirmwareUpdateReportParser) getReason(code int32) string {
	switch code {
	case 100:
		return "The update command is confirmed by the device"
	case 101:
		return "The update command is refused by the device"
	case 102:
		return "The update process is canceled by the backend server"
	case 103:
		return "The update process is refused because the battery is low"
	case 200:
		return "The device starts to download the package"
	case 201:
		return "The device finishes downloading the package successfully"
	case 202:
		return "The device fails to download the package"
	case 300:
		return "The device starts to update the firmware"
	case 301:
		return "The device finishes updating the firmware successfully"
	case 302:
		return "The device fails to update the firmware"
	case 303:
		return "The update process does not start because the battery is low"
	}
	return ""
}
