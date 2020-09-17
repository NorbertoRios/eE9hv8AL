package queclinkreport

import (
	"fmt"
	"strings"

	"queclink-go/queclinkreport/devicetypes"
	"queclink-go/queclinkreport/gv300/resp300"

	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
)

const (
	hbdDeviceTypePos       = 6
	ackDeviceTypePos       = 7
	reportDeviceTypePos    = 11
	reportDeviceTypeEriPos = 15
	reportDeviceTypeCraPos = 8
	reportDeviceTypeInfPos = 17
)

//Parser queclink's service for
type Parser struct {
}

//Parse queclink packets
func (parser *Parser) Parse(packet []byte) ([]report.IMessage, error) {
	if parser.IsTextReport(packet) {
		textParser := &TextParser{}
		messages := textParser.Parse(packet, nil)
		value := []report.IMessage{}
		if messages != nil && len(messages) > 0 {
			value = append(value, messages[0])
		}
		return value, nil
	}
	deviceType, messageHeader, messageType, err := parser.GetHeader(packet)
	if err != nil {
		return nil, err
	}

	switch deviceType {
	case devicetypes.GV55,
		devicetypes.GV55Lite,
		devicetypes.GV55N,
		devicetypes.GV55NLite:
		return parser.ParseGV55(deviceType, messageHeader, messageType, packet)
	case devicetypes.GV320:
		//parse gv320
		return nil, nil
	case devicetypes.GV75,
		devicetypes.GV75W:
		return parser.ParseGV75(deviceType, messageHeader, messageType, packet)
	case devicetypes.GV55W:
		return parser.ParseGV55w(deviceType, messageHeader, messageType, packet)
	case devicetypes.GV600W, devicetypes.GV600Fake, devicetypes.GV600MG:
		return parser.ParseGV600(deviceType, messageHeader, messageType, packet)
	case devicetypes.GV350MA:
		return parser.ParseGV350(deviceType, messageHeader, messageType, packet)
	case devicetypes.GV300W:
		return parser.ParseGV300w(deviceType, messageHeader, messageType, packet)
	default:
		return parser.ParseGV55(deviceType, messageHeader, messageType, packet)
	}
}

//GetUnknownAck returns ack for packet if possible
func (parser *Parser) GetUnknownAck(packet []byte) []byte {
	if packet == nil || len(packet) < 10 {
		return nil
	}
	countNumber, _ := utils.GetBytesValue(packet, len(packet)-6, 2)
	hexCn := fmt.Sprintf("%04X", countNumber)
	ack := strings.ToUpper(fmt.Sprintf("+SACK:%v$", hexCn))
	return []byte(ack)
}

//GetHeader parses packet and return device type, message header, message type
func (parser *Parser) GetHeader(packet []byte) (deviceType int, messageHeader string, messageType int, err error) {
	messageHeader = string(packet[0:4])
	messageType = int(packet[4]) //1 byte for message type
	deviceType = 16

	switch messageHeader {
	case "+HBD",
		"+BBD":
		{
			if len(packet) < hbdDeviceTypePos+1 {
				err = fmt.Errorf("Invalid packet length")
				return 0, "", 0, err
			}
			deviceType = int(packet[hbdDeviceTypePos])
			break
		}

	case "+ACK":
		{
			if len(packet) < ackDeviceTypePos+1 {
				err = fmt.Errorf("Invalid packet length")
				return 0, "", 0, err
			}
			deviceType = int(packet[ackDeviceTypePos])
			break
		}
	case "+CRD",
		"+BRD":
		{
			if len(packet) < reportDeviceTypeCraPos+1 {
				err = fmt.Errorf("Invalid packet length")
				return 0, "", 0, err
			}
			deviceType = int(packet[reportDeviceTypeCraPos])
			break
		}

	case "+INF",
		"+BNF":
		{
			if len(packet) < reportDeviceTypeInfPos+1 {
				err = fmt.Errorf("Invalid packet length")
				return 0, "", 0, err
			}
			deviceType = int(packet[reportDeviceTypeInfPos])
			break
		}

	default:
		{
			if (messageType == resp300.GTERI) &&
				(messageHeader == "+BSP" || messageHeader == "+RSP") {
				if len(packet) < reportDeviceTypeEriPos+1 {
					err = fmt.Errorf("Invalid packet length")
					return 0, "", 0, err
				}
				deviceType = int(packet[reportDeviceTypeEriPos])
			} else {
				if len(packet) < reportDeviceTypePos+1 {
					err = fmt.Errorf("Invalid packet length")
					return 0, "", 0, err
				}
				deviceType = int(packet[reportDeviceTypePos])
			}
			break
		}
	}
	return deviceType, messageHeader, messageType, nil
}

//IsTextReport checks packet is text report
func (parser *Parser) IsTextReport(packet []byte) bool {
	reportHeader := string(packet[:5])
	return reportHeader == "+RESP" || reportHeader == "+BUFF"
}

//ParseGV55 packet and create array of messages
func (parser *Parser) ParseGV55(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	switch messageHeader {
	case "+RSP",
		"+BSP":
		return parser.ParseRespGV55(deviceType, messageHeader, messageType, packet)
	case "+BVT",
		"+EVT":
		return parser.ParseEventGV55(deviceType, messageHeader, messageType, packet)
	case "+CRD",
		"+BRD":
		return parser.ParseCrashReport(deviceType, messageHeader, packet)
	case "+ACK":
		return parser.ParseAckReport(deviceType, messageHeader, messageType, packet)
	case "+HBD",
		"+BBD":
		return parser.ParseHbdReport(deviceType, messageHeader, packet)
	case "+INF",
		"+BNF":
		return parser.ParseInfReport(deviceType, messageHeader, messageType, packet)
	}
	return nil, nil
}

//ParseGV75 packet and create array of messages
func (parser *Parser) ParseGV75(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	switch messageHeader {
	case "+RSP",
		"+BSP":
		return parser.ParseRespGV75(deviceType, messageHeader, messageType, packet)
	case "+BVT",
		"+EVT":
		return parser.ParseEventGV75(deviceType, messageHeader, messageType, packet)
	case "+CRD",
		"+BRD":
		return parser.ParseCrashReport(deviceType, messageHeader, packet)
	case "+ACK":
		return parser.ParseAckReport(deviceType, messageHeader, messageType, packet)
	case "+HBD",
		"+BBD":
		return parser.ParseHbdReport(deviceType, messageHeader, packet)
	case "+INF",
		"+BNF":
		return parser.ParseInfReport(deviceType, messageHeader, messageType, packet)
	}
	return nil, nil
}

//ParseGV55w packet and create array of messages
func (parser *Parser) ParseGV55w(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	switch messageHeader {
	case "+RSP",
		"+BSP":
		return parser.ParseRespGV55w(deviceType, messageHeader, messageType, packet)
	case "+BVT",
		"+EVT":
		return parser.ParseEventGV55w(deviceType, messageHeader, messageType, packet)
	case "+CRD",
		"+BRD":
		return parser.ParseCrashReport(deviceType, messageHeader, packet)
	case "+ACK":
		return parser.ParseAckReport(deviceType, messageHeader, messageType, packet)
	case "+HBD",
		"+BBD":
		return parser.ParseHbdReport(deviceType, messageHeader, packet)
	case "+INF",
		"+BNF":
		return parser.ParseInfReport(deviceType, messageHeader, messageType, packet)
	}
	return nil, nil
}

//ParseGV600 packet and create array of messages
func (parser *Parser) ParseGV600(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	switch messageHeader {
	case "+RSP",
		"+BSP":
		return parser.ParseRespGV600(deviceType, messageHeader, messageType, packet)
	case "+BVT",
		"+EVT":
		return parser.ParseEventGV600(deviceType, messageHeader, messageType, packet)
	case "+CRD",
		"+BRD":
		return parser.ParseCrashReport(deviceType, messageHeader, packet)
	case "+ACK":
		return parser.ParseAckReport(deviceType, messageHeader, messageType, packet)
	case "+HBD",
		"+BBD":
		return parser.ParseHbdReport(deviceType, messageHeader, packet)
	case "+INF",
		"+BNF":
		return parser.ParseInfReport(deviceType, messageHeader, messageType, packet)
	}
	return nil, nil
}

//ParseGV350 packet and create array of messages
func (parser *Parser) ParseGV350(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	switch messageHeader {
	case "+RSP",
		"+BSP":
		return parser.ParseRespGV350(deviceType, messageHeader, messageType, packet)
	case "+BVT",
		"+EVT":
		return parser.ParseEventGV350(deviceType, messageHeader, messageType, packet)
	case "+CRD",
		"+BRD":
		return parser.ParseCrashReport(deviceType, messageHeader, packet)
	case "+ACK":
		return parser.ParseAckReport(deviceType, messageHeader, messageType, packet)
	case "+HBD",
		"+BBD":
		return parser.ParseHbdReport(deviceType, messageHeader, packet)
	case "+INF",
		"+BNF":
		return parser.ParseInfReport(deviceType, messageHeader, messageType, packet)
	}
	return nil, nil
}

//ParseGV300w packet and create array of messages
func (parser *Parser) ParseGV300w(deviceType int, messageHeader string, messageType int, packet []byte) ([]report.IMessage, error) {
	switch messageHeader {
	case "+RSP",
		"+BSP":
		return parser.ParseRespGV300W(deviceType, messageHeader, messageType, packet)
	case "+BVT",
		"+EVT":
		return parser.ParseEventGV300W(deviceType, messageHeader, messageType, packet)
	case "+CRD",
		"+BRD":
		return parser.ParseCrashReport(deviceType, messageHeader, packet)
	case "+ACK":
		return parser.ParseAckReport(deviceType, messageHeader, messageType, packet)
	case "+HBD",
		"+BBD":
		return parser.ParseHbdReport(deviceType, messageHeader, packet)
	case "+INF",
		"+BNF":
		return parser.ParseInfReport(deviceType, messageHeader, messageType, packet)
	}
	return nil, nil
}
