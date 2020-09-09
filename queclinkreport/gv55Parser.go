package queclinkreport

import (
	"fmt"

	"queclink-go/base.device.service/report"
	"queclink-go/queclinkreport/gv55/evt55"
	"queclink-go/queclinkreport/gv55/resp55"
)

//ParseRespGV55 parses packet for
func (parser *Parser) ParseRespGV55(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	reportConfig, err := report.ReportConfiguration.(*ReportConfiguration).Find(deviceType, messageHeader, messageType)
	if err != nil || reportConfig == nil {
		return nil, fmt.Errorf("Not found configuration for device type %v message header:%v message type:%v", deviceType, messageHeader, messageType)
	}
	var messages []report.IMessage
	switch messageType {
	case resp55.GTLBC:
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

//ParseEventGV55 parses event report for gv55 devices
func (parser *Parser) ParseEventGV55(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	reportConfig, err := report.ReportConfiguration.(*ReportConfiguration).Find(deviceType, messageHeader, messageType)
	if err != nil || reportConfig == nil {
		return nil, fmt.Errorf("Not found configuration for device type %v message header:%v message type:%v", deviceType, messageHeader, messageType)
	}
	var messages []report.IMessage
	switch messageType {
	case evt55.GTUPD:
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

//ParseCrashReport parses packets for device
func (parser *Parser) ParseCrashReport(deviceType int, messageHeader string, packet []byte) ([]report.IMessage, error) {
	reportConfig, err := report.ReportConfiguration.(*ReportConfiguration).Find(deviceType, messageHeader, 0)
	if err != nil || reportConfig == nil {
		return nil, fmt.Errorf("Not found configuration for device type %v message header:%v", deviceType, messageHeader)
	}
	return (&(CrashReportParser{})).Parse(packet, reportConfig), nil
}

//ParseAckReport parses ack reports
func (parser *Parser) ParseAckReport(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	reportConfig, err := report.ReportConfiguration.(*ReportConfiguration).Find(deviceType, messageHeader, messageType)
	if err != nil || reportConfig == nil {
		return nil, fmt.Errorf("Not found configuration for device type %v message header:%v message type:%v", deviceType, messageHeader, messageType)
	}

	return (&(AckReportParser{})).Parse(packet, reportConfig), nil
}

//ParseHbdReport parses ack reports
func (parser *Parser) ParseHbdReport(deviceType int, messageHeader string, packet []byte) ([]report.IMessage, error) {
	reportConfig, err := report.ReportConfiguration.(*ReportConfiguration).Find(deviceType, messageHeader, 0)
	if err != nil || reportConfig == nil {
		return nil, fmt.Errorf("Not found configuration for device type %v message header:%v message type:%v", deviceType, messageHeader, 0)
	}
	return (&(HBDReportParser{})).Parse(packet, reportConfig), nil
}

//ParseInfReport parses ack reports
func (parser *Parser) ParseInfReport(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	reportConfig, err := report.ReportConfiguration.(*ReportConfiguration).Find(deviceType, messageHeader, messageType)
	if err != nil || reportConfig == nil {
		return nil, fmt.Errorf("Not found configuration for device type %v message header:%v message type:%v", deviceType, messageHeader, messageType)
	}

	return (&(InfReportParser{})).Parse(packet, reportConfig), nil
}
