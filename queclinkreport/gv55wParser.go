package queclinkreport

import (
	"fmt"

	"queclink-go/base.device.service/report"
	"queclink-go/queclinkreport/gv55w/evt55w"
	"queclink-go/queclinkreport/gv55w/resp55w"
)

//ParseRespGV55w parses packet for
func (parser *Parser) ParseRespGV55w(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	reportConfig, err := report.ReportConfiguration.(*ReportConfiguration).Find(deviceType, messageHeader, messageType)
	if err != nil || reportConfig == nil {
		return nil, fmt.Errorf("Not found configuration for device type %v message header:%v message type:%v", deviceType, messageHeader, messageType)
	}
	var messages []report.IMessage
	switch messageType {
	case resp55w.GTLBC:
		{
			messages = (&(LBCReportParser{})).Parse(packet, reportConfig)
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

//ParseEventGV55w parses event report for gv55 devices
func (parser *Parser) ParseEventGV55w(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	reportConfig, err := report.ReportConfiguration.(*ReportConfiguration).Find(deviceType, messageHeader, messageType)
	if err != nil || reportConfig == nil {
		return nil, fmt.Errorf("Not found configuration for device type %v message header:%v message type:%v", deviceType, messageHeader, messageType)
	}
	var messages []report.IMessage
	switch messageType {
	case evt55w.GTUPD:
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
