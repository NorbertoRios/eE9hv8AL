package queclinkreport

import (
	"fmt"

	"queclink-go/queclinkreport/gv300w/evt300w"

	"queclink-go/base.device.service/report"
	"queclink-go/queclinkreport/gv300w/resp300w"
)

//ParseRespGV300W parses packet for
func (parser *Parser) ParseRespGV300W(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	reportConfig, err := report.ReportConfiguration.(*ReportConfiguration).Find(deviceType, messageHeader, messageType)
	if err != nil || reportConfig == nil {
		return nil, fmt.Errorf("Not found configuration for device type %v message header:%v message type:%v", deviceType, messageHeader, messageType)
	}
	var messages []report.IMessage
	switch messageType {
	case resp300w.GTLBC:
		{
			messages = (&(LBCReportParser{})).Parse(packet, reportConfig)
		}
	case resp300w.GTERI:
		{
			messages = (&(GV300WEriParser{})).Parse(packet, reportConfig)
		}
	default:
		{
			messages = (&(LocationReportParser{})).Parse(packet, reportConfig)
		}

	}
	if messages == nil || len(messages) == 0 {
		return nil, fmt.Errorf("Parsing result is empty for message type:%v", messageType)
	}
	return messages, nil
}

//ParseEventGV300W parses event report for gv55 devices
func (parser *Parser) ParseEventGV300W(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	reportConfig, err := report.ReportConfiguration.(*ReportConfiguration).Find(deviceType, messageHeader, messageType)
	if err != nil || reportConfig == nil {
		return nil, fmt.Errorf("Not found configuration for device type %v message header:%v message type:%v", deviceType, messageHeader, messageType)
	}
	var messages []report.IMessage
	switch messageType {
	case evt300w.GTUPD:
		{
			messages = (&(FirmwareUpdateReportParser{})).Parse(packet, reportConfig)
		}
	default:
		{
			messages = (&(LocationReportParser{})).Parse(packet, reportConfig)
		}

	}
	if messages == nil || len(messages) == 0 {
		return nil, fmt.Errorf("Parsing result is empty for message type:%v", messageType)
	}
	return messages, nil
}
