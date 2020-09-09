package queclinkreport

import (
	"fmt"

	"queclink-go/base.device.service/report"
	"queclink-go/queclinkreport/gv350/evt350"
	"queclink-go/queclinkreport/gv350/resp350"
)

//ParseRespGV350 parses packet for
func (parser *Parser) ParseRespGV350(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	reportConfig, err := report.ReportConfiguration.(*ReportConfiguration).Find(deviceType, messageHeader, messageType)
	if err != nil || reportConfig == nil {
		return nil, fmt.Errorf("Not found configuration for device type %v message header:%v message type:%v", deviceType, messageHeader, messageType)
	}
	var messages []report.IMessage
	switch messageType {
	case resp350.GTLBC:
		{
			messages = (&(GV350LBCReportParser{})).Parse(packet, reportConfig)
		}
	case resp350.GTERI:
		{
			messages = (&(GV350EriReportParser{})).Parse(packet, reportConfig)
		}
	default:
		{
			messages = (&(GV350LocationReportParser{})).Parse(packet, reportConfig)
		}

	}
	if messages == nil || len(messages) == 0 {
		return nil, fmt.Errorf("Parsing result is empty for message type:%v", messageType)
	}
	return messages, nil
}

//ParseEventGV350 parses event report for gv55 devices
func (parser *Parser) ParseEventGV350(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	reportConfig, err := report.ReportConfiguration.(*ReportConfiguration).Find(deviceType, messageHeader, messageType)
	if err != nil || reportConfig == nil {
		return nil, fmt.Errorf("Not found configuration for device type %v message header:%v message type:%v", deviceType, messageHeader, messageType)
	}
	var messages []report.IMessage
	switch messageType {
	case evt350.GTUPD:
		{
			messages = (&(FirmwareUpdateReportParser{})).Parse(packet, reportConfig)
		}
	case evt350.GTDAT, evt350.GTUPC:
		{
			messages = (&(GV600DatReportParser{})).Parse(packet, reportConfig)
		}
	default:
		{
			messages = (&(GV350LocationReportParser{})).Parse(packet, reportConfig)
		}

	}
	if messages == nil || len(messages) == 0 {
		return nil, fmt.Errorf("Parsing result is empty for message type:%v", messageType)
	}
	return messages, nil
}
