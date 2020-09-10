package queclinkreport

import (
	"testing"

	"queclink-go/queclinkreport/fields"
)

func TestTextReport(t *testing.T) {
	config := &ReportConfiguration{}
	LoadReportConfiguration("..", "/ReportConfiguration.xml", config)
	packet := []byte("+RESP:GTDOG,2F0500,135790246811220,,,01,1,1,4.3,92,70.0,121.354335,31.222073,20090214013254,0460,0000,18d8,6141,00,2000.0,20090214093254,11F0$")
	parser := &TextParser{}
	message := parser.Parse(packet, nil)
	if message == nil {
		t.Error("[TestTextReport]message is nil")
	}
	if message[0].MessageType() != "+RSP" {
		t.Error("[TestTextReport]Invalid Message Header")
	}
	if message[0].EventCode() != 12 {
		t.Error("[TestTextReport]Invalid Message type")
	}
	iack, found := message[0].GetValue(fields.Ack)
	ack := iack.([]byte)
	if !found || string(ack) != "+SACK:11F0$" {
		t.Error("[TestTextReport]Invalid ack")
	}

}
