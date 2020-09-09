package queclinkreport

import (
	"encoding/json"
	"sync"

	"queclink-go/base.device.service/report"
)

//IQueclinkMessage interface for queclink message type
type IQueclinkMessage interface {
	report.IMessage
	IsText() bool
}

//QueclinkMessage type for queclink message
type QueclinkMessage struct {
	report.Message
}

//IsText returns report is text
func (m *QueclinkMessage) IsText() bool {
	mt := m.MessageType()
	return mt == "+RESP" || mt == "+BESP"
}

//MessageType return type of message
func (m *QueclinkMessage) MessageType() string {
	h, found := m.GetValue("MessageHeader")
	if found {
		return h.(string)
	}
	return ""
}

//EventCode returns report type
func (m *QueclinkMessage) EventCode() int32 {
	//t, found := m.Data["MessageType"]
	t, found := m.GetValue("MessageType")
	if found {
		return t.(int32)
	}
	return -1
}

//DeviceType returns device type
func (m *QueclinkMessage) DeviceType() int32 {
	t, found := m.GetValue("DeviceType")
	//t, found := m.Data["DeviceType"]
	if found {
		return t.(int32)
	}
	return -1
}

//NewQueclinkMessage returns new instance of QueclinkMessage
func NewQueclinkMessage() *QueclinkMessage {
	message := &QueclinkMessage{}
	message.Data = make(map[string]interface{})
	message.Mtx = &sync.Mutex{}
	return message
}

//UnMarshalMessage given string to Message struct
func UnMarshalMessage(str string) (*QueclinkMessage, error) {
	message := &QueclinkMessage{}
	err := json.Unmarshal([]byte(str), message)
	if err != nil {
		return &QueclinkMessage{}, err
	}
	return message, err
}
