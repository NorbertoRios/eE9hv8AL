package queclinkreport

import (
	"fmt"

	"queclink-go/base.device.service/report"
	"queclink-go/queclinkreport/gv75/evt75"
	"queclink-go/queclinkreport/gv75/resp75"
)

//ParseRespGV75 parses packet for
func (parser *Parser) ParseRespGV75(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	reportConfig, err := report.ReportConfiguration.(*ReportConfiguration).Find(deviceType, messageHeader, messageType)
	if err != nil || reportConfig == nil {
		return nil, fmt.Errorf("Not found configuration for device type %v message header:%v message type:%v", deviceType, messageHeader, messageType)
	}
	var messages []report.IMessage
	switch messageType {
	case resp75.GTLBC:
		{
			messages = (&(LBCReportParser{})).Parse(packet, reportConfig)
		}
	case resp75.GTERI:
		{
			messages = (&(EriReportParser{})).Parse(packet, reportConfig)
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

//ParseEventGV75 parses event report for gv55 devices
func (parser *Parser) ParseEventGV75(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	reportConfig, err := report.ReportConfiguration.(*ReportConfiguration).Find(deviceType, messageHeader, messageType)
	if err != nil || reportConfig == nil {
		return nil, fmt.Errorf("Not found configuration for device type %v message header:%v message type:%v", deviceType, messageHeader, messageType)
	}
	var messages []report.IMessage
	switch messageType {
	case evt75.GTUPD:
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
