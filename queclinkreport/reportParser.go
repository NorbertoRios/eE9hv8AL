package queclinkreport

import (
	"fmt"
	"strings"

	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
	"queclink-go/queclinkreport/fields"
)

//ReportParser parser
type ReportParser struct {
	Packet              []byte
	ReportConfiguration *Report
	StartByte           int
	MessageHeader       string
	MessageType         int32
	ReportMask          int32
}

//Init base parser
func (parser *ReportParser) Init(packet *[]byte) {
	parser.MessageHeader = string((*packet)[0:4])
	parser.MessageType = int32((*packet)[4])
	ireportMask, _ := utils.GetBytesValue(*packet, 5, 4)
	parser.ReportMask = ireportMask.(int32)
	parser.StartByte = 9
}

//Parse packet
func (parser *ReportParser) Parse(packet []byte, config *Report) []report.IMessage {
	return nil
}

//SetAck sets ack for message
func (parser *ReportParser) SetAck(message report.IMessage) {

	icn, found := message.GetValue(fields.CountNumber)
	if !found {
		return
	}

	var hexCn = ""

	if icn == nil {
		ipacket, f := message.GetValue("RawData")
		if !f {
			return
		}
		if packet, valid := ipacket.([]byte); valid {
			countNumber, _ := utils.GetBytesValue(packet, len(packet)-6, 2)
			hexCn = fmt.Sprintf("%04X", countNumber.(int32))
		}
	} else {
		hexCn = fmt.Sprintf("%04X", icn.(int32))
	}
	cn := ""
	if message.MessageType() == "+HBD" {
		cn = fmt.Sprintf("GTHBD,%v,%v", "", hexCn)
	} else if message.MessageType() == "+BBD" {
		cn = fmt.Sprintf("GTBBD,%v,%v", "", hexCn)
	} else {
		cn = fmt.Sprintf("%v", hexCn)
	}
	ack := strings.ToUpper(fmt.Sprintf("+SACK:%v$", cn))
	message.SetValue(fields.Ack, []byte(ack))
	message.SetValue("SACK", ack)
}
