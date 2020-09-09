package queclinkreport

import (
	"testing"

	"queclink-go/queclinkreport/fields"
)

func TestSAck(t *testing.T) {
	message := NewQueclinkMessage()
	message.SetValue(fields.CountNumber, int32(12))
	message.SetValue(fields.MessageHeader, "+RSP")

	parser := &ReportParser{}
	parser.SetAck(message)
	ack, found := message.GetValue(fields.Ack)
	if !found {
		t.Error("[TestSAck]Ack not found")
	}
	if string(ack.([]byte)) != "+SACK:000C$" {
		t.Error("[TestSAck]Invalid ack")
	}
}

func TestHBDSAck(t *testing.T) {
	message := NewQueclinkMessage()
	message.SetValue(fields.CountNumber, int32(12))
	message.SetValue(fields.MessageHeader, "+HBD")

	parser := &ReportParser{}
	parser.SetAck(message)
	ack, found := message.GetValue(fields.Ack)
	if !found {
		t.Error("[TestSAck]Ack not found")
	}
	if string(ack.([]byte)) != "+SACK:GTHBD,,000C$" {
		t.Error("[TestSAck]Invalid ack")
	}
}

func TestBBDSAck(t *testing.T) {
	message := NewQueclinkMessage()
	message.SetValue(fields.CountNumber, int32(12))
	message.SetValue(fields.MessageHeader, "+BBD")

	parser := &ReportParser{}
	parser.SetAck(message)
	ack, found := message.GetValue(fields.Ack)
	if !found {
		t.Error("[TestSAck]Ack not found")
	}
	if string(ack.([]byte)) != "+SACK:GTBBD,,000C$" {
		t.Error("[TestSAck]Invalid ack")
	}
}
