package queclinkreport

import (
	"queclink-go/base.device.service/report"
	"queclink-go/queclinkreport/fields"
)

//GV350LocationReportParser parses location packets
type GV350LocationReportParser struct {
	BaseLocationReportParser
}

//Parse location reports packet
func (parser *GV350LocationReportParser) Parse(packet []byte, config *Report) []report.IMessage {
	parser.Init(&packet)
	messages := parser.BaseLocationReportParser.Parse(packet, config)
	parser.InvalidateFields(config, &messages)
	return messages
}

//InvalidateFields makes fields human readable
func (parser *GV350LocationReportParser) InvalidateFields(config *Report, messages *[]report.IMessage) {
	if len(*messages) > 0 {
		parser.SetAck((*messages)[0])
	}
	for _, message := range *messages {
		parser.HumanizeUniqueID(message.(IQueclinkMessage), fields.UniqueID)
		parser.HumanizeSpeed(message.(IQueclinkMessage))
		parser.HumanizeDateTime(message.(IQueclinkMessage), fields.GPSUTCTime)
		parser.HumanizeDateTime(message.(IQueclinkMessage), fields.SendTime)
		parser.HumanizeCurrentMileage(message.(IQueclinkMessage))
		parser.HumanizeTotalMileage(message.(IQueclinkMessage))
		parser.HumanizeCurrentHourMeterCount(message.(IQueclinkMessage))
		parser.HumanizeTotalHourMeterCount(message.(IQueclinkMessage))
		parser.HumanizeTimestamp(message.(IQueclinkMessage))
		parser.HumanizeLngLat(message.(IQueclinkMessage), fields.Latitude)
		parser.HumanizeLngLat(message.(IQueclinkMessage), fields.Longitude)
		parser.HumanizeInputs(message.(IQueclinkMessage))
		parser.HumanizeValidity(config, message.(IQueclinkMessage))
		parser.HumanizeInt32ToFloat32(fields.Altitude, message.(IQueclinkMessage))
		parser.HumanizeInt32ToFloat32(fields.Heading, message.(IQueclinkMessage))
		parser.HumanizeSatellites(message.(IQueclinkMessage))
		parser.HumanizeBatteryPercentage(message.(IQueclinkMessage))
		parser.DecodeReason(message.(IQueclinkMessage))
		message.SetValue(fields.LocationMessage, config.Location)
	}
}

//HumanizeSatellites converts satellites
func (parser *GV350LocationReportParser) HumanizeSatellites(message IQueclinkMessage) {
	if iv, f := message.GetValue(fields.Satellites); f {
		if v, valid := iv.(byte); valid {
			message.SetValue(fields.Satellites, int32(v&0x0F))
			message.SetValue("ExternalGPSAntennaStatus", int32((v>>6)&0x03))
		}
	}
}
