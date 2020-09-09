package queclinkreport

import (
	"fmt"
	"strconv"

	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
	"queclink-go/queclinkreport/fields"
)

//TextParser parser
type TextParser struct {
	ReportParser
}

//Parse text report
func (parser *TextParser) Parse(packet []byte, reportConfig *Report) []IQueclinkMessage {
	message := NewQueclinkMessage()
	parser.SetMessageType(packet, message)
	input := string(packet)
	strPacket := utils.SplitString(input, []string{",", "$"})

	countNumber, err := strconv.ParseUint(strPacket[len(strPacket)-1], 16, 32)
	if err != nil {
		return nil
	}
	message.SetValue(fields.CountNumber, int32(countNumber))
	message.SetValue(fields.DeviceType, 0)

	parser.SetAck(message)

	message.SetValue(fields.UniqueID, strPacket[2])
	message.SetValue("DevId", "queclink_"+message.GetStringValue("UniqueId", ""))
	message.SetValue("Value", input)
	return []IQueclinkMessage{message}
}

//SetMessageType sets message type for packet
func (parser *TextParser) SetMessageType(packet []byte, message *QueclinkMessage) error {
	input := string(packet)
	var mt = input[6:11]
	deviceType := 14
	deviceConfigs := report.ReportConfiguration.(*ReportConfiguration)
	deviceConfig, err := deviceConfigs.FindDeviceType(deviceType)
	if err != nil {
		return fmt.Errorf("Device config not foundd")
	}

	item, err := deviceConfig.FindRespMessagesType(mt)
	if err != nil || item == nil {
		return fmt.Errorf("Resp report config not foundd")

	}
	message.SetValue(fields.MessageHeader, item.Class)
	message.SetValue(fields.MessageType, int32(item.ID))
	return nil
}
