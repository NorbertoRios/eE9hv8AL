package queclinkreport

import (
	"fmt"

	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
	"queclink-go/queclinkreport/fields"
)

//LBCReportParser location by call message parser
type LBCReportParser struct {
	BaseLocationReportParser
}

//Parse LBC packet
func (parser *LBCReportParser) Parse(packet []byte, config *Report) []report.IMessage {
	parser.Init(&packet)
	messages := []report.IMessage{}
	readFrom := parser.StartByte
	header := parser.ParseHeader(parser.ReportMask, packet, config.Header, &readFrom)
	header.SetValue(fields.MessageHeader, parser.MessageHeader)
	header.SetValue(fields.MessageType, parser.MessageType)
	header.SetValue(fields.ReportMask, parser.ReportMask)
	header.SetValue("RawData", packet)

	number := byte(0)
	if v, found := header.GetValue(fields.Number); found {
		number = v.(byte)
	}

	for i := 1; i <= int(number); i++ {
		body := parser.ParseBody(parser.ReportMask, packet, config.MultiPosition, &readFrom)
		messages = append(messages, body)
	}
	var tail = parser.ParseTail(parser.ReportMask, packet, config.Tail, &readFrom)

	for _, report := range messages {
		report.AppendRange(*header.GetData())
		report.AppendRange(*tail.GetData())
	}

	if len(messages) == 0 && (len(*header.GetData()) != 0 || len(*tail.GetData()) != 0) {
		header.AppendRange(*tail.GetData())
		messages = append(messages, header)
	}
	parser.InvalidateFields(config, &messages)
	return messages
}

//ParseHeader of message
func (parser *LBCReportParser) ParseHeader(reportMask int32, packet []byte, confItems []Item, readFrom *int) IQueclinkMessage {
	message := NewQueclinkMessage()
	//in config.Where(reportItem => ByteIsSet(ReportMask, reportItem.MaskBit))
	for _, item := range confItems {
		if item.MaskBit != nil && !utils.BitIsSet(int64(reportMask), uint(*item.MaskBit)) {
			continue
		}
		fieldSize := item.Size
		if item.ItemName == fields.PhoneNumber {
			inumber, found := message.GetValue(fields.NumberLengthType)
			if !found {
				continue
			}
			number, valid := inumber.(byte)
			if !valid {
				continue
			}
			high := number >> 4
			fieldSize = int(high - 1)
		}

		value, _ := utils.GetBytesValue(packet, *readFrom, fieldSize)
		parser.OnValueParsed(confItems, item.ItemName, value)
		message.SetValue(item.ItemName, value)
		*readFrom = *readFrom + fieldSize
	}
	return message
}

//InvalidateFields makes fields human readable
func (parser *LBCReportParser) InvalidateFields(config *Report, messages *[]report.IMessage) {
	parser.BaseLocationReportParser.InvalidateFields(config, messages)
	for _, message := range *messages {
		parser.HumanizePhone(message.(IQueclinkMessage))
	}
}

//HumanizePhone makes phone radable
func (parser *LBCReportParser) HumanizePhone(message IQueclinkMessage) {
	pn, found := message.GetValue(fields.PhoneNumber)
	if !found {
		return
	}
	bPn := pn.([]byte)
	var phoneNumber string
	for i, v := range bPn {
		high := v >> 4
		low := v & 0x0F
		if low == 0xF && i == len(bPn)-1 {
			break
		}
		phoneNumber = fmt.Sprintf("%v%v%v", phoneNumber, high, low)
	}
	inl, found := message.GetValue(fields.NumberLengthType)
	if found && (inl.(byte)&15) == 1 {
		phoneNumber = fmt.Sprintf("+%v", phoneNumber)
	}
	message.SetValue(fields.PhoneNumber, phoneNumber)
}
