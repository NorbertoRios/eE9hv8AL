package queclinkreport

import (
	"fmt"

	"queclink-go/base.device.service/report"
	"queclink-go/queclinkreport/gv600/evt600"
	"queclink-go/queclinkreport/gv600/resp600"
)

//ParseRespGV600 parses packet for
func (parser *Parser) ParseRespGV600(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	reportConfig, err := report.ReportConfiguration.(*ReportConfiguration).Find(deviceType, messageHeader, messageType)
	if err != nil || reportConfig == nil {
		return nil, fmt.Errorf("Not found configuration for device type %v message header:%v message type:%v", deviceType, messageHeader, messageType)
	}
	var messages []report.IMessage
	switch messageType {
	case resp600.GTLBC:
		{
			messages = (&(LBCReportParser{})).Parse(packet, reportConfig)
		}
	case resp600.GTERI:
		{
			messages = (&(GV600EriReportParser{})).Parse(packet, reportConfig)
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

//ParseEventGV600 parses event report for gv55 devices
func (parser *Parser) ParseEventGV600(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	reportConfig, err := report.ReportConfiguration.(*ReportConfiguration).Find(deviceType, messageHeader, messageType)
	if err != nil || reportConfig == nil {
		return nil, fmt.Errorf("Not found configuration for device type %v message header:%v message type:%v", deviceType, messageHeader, messageType)
	}
	var messages []report.IMessage
	switch messageType {
	case evt600.GTUPD:
		{
			messages = (&(FirmwareUpdateReportParser{})).Parse(packet, reportConfig)
		}
	case evt600.GTDAT, evt600.GTUPC:
		{
			messages = (&(GV600DatReportParser{})).Parse(packet, reportConfig)
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
