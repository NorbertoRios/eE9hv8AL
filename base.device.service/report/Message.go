package report

import (
	"encoding/json"
	"sync"

	"queclink-go/base.device.service/utils"
)

//IMessage interface for Message
type IMessage interface {
	EventCode() int32
	MessageType() string
	UniqueID() string
	LocationMessage() bool
	ReceivedTime() *utils.JSONTime
	SetValue(key string, value interface{})
	GetValue(key string) (value interface{}, found bool)
	GetStringValue(key string, def string) (value string)
	GetIntValue(key string, def int32) (value int32)
	GetFloatValue(key string, def float32) (value float32)
	RemoveKey(fn string)
	SourceID() uint64
	SetSourceID(uint64)
	GetData() *map[string]interface{}
	AppendRange(data map[string]interface{})
	GetTemperatureSensors() *TemperatureSensors
	SetTemperatureSensors(temperatureSensors *TemperatureSensors)
}

//Message struct for parsed messages
type Message struct {
	Mtx                *sync.Mutex         `json:"-"`
	SID                uint64              `json:"sid"`
	TemperatureSensors *TemperatureSensors `json:"ts,omitempty"`
	Data               map[string]interface{}
}

//GetData returns Data map
func (m *Message) GetData() *map[string]interface{} {
	m.Mtx.Lock()
	defer m.Mtx.Unlock()
	return &m.Data
}

//EventCode of message
func (m *Message) EventCode() int32 {
	m.Mtx.Lock()
	data := m.Data["EventCode"]
	m.Mtx.Unlock()
	return data.(int32)
}

//MessageType of message
func (m *Message) MessageType() string {
	m.Mtx.Lock()
	tp := m.Data["MessageType"]
	m.Mtx.Unlock()
	return tp.(string)
}

//ReceivedTime of message
func (m *Message) ReceivedTime() *utils.JSONTime {
	m.Mtx.Lock()
	tiime := m.Data["ReceivedTime"]
	m.Mtx.Unlock()
	return tiime.(*utils.JSONTime)
}

//UniqueID of message
func (m *Message) UniqueID() string {
	m.Mtx.Lock()
	uId := m.Data["UniqueId"]
	m.Mtx.Unlock()
	return uId.(string)
}

//LocationMessage indicate message is location
func (m *Message) LocationMessage() bool {
	m.Mtx.Lock()
	isLocation := m.Data["LocationMessage"]
	m.Mtx.Unlock()
	return isLocation.(bool)
}

//SetValue to Data field
func (m *Message) SetValue(key string, value interface{}) {
	m.Mtx.Lock()
	m.Data[key] = value
	m.Mtx.Unlock()
}

//GetValue from Data field
func (m *Message) GetValue(key string) (value interface{}, found bool) {
	m.Mtx.Lock()
	value, found = m.Data[key]
	m.Mtx.Unlock()
	return value, found
}

//GetStringValue returns string value
func (m *Message) GetStringValue(key string, def string) (value string) {
	v, found := m.GetValue(key)
	if found {
		return v.(string)
	}
	return def
}

//GetIntValue returns int32 value
func (m *Message) GetIntValue(key string, def int32) (value int32) {
	if iv, found := m.GetValue(key); found {
		if v, valid := iv.(int32); valid {
			return v
		}
	}
	return def
}

//GetFloatValue returns float32 value
func (m *Message) GetFloatValue(key string, def float32) (value float32) {
	v, found := m.GetValue(key)
	if found {
		return v.(float32)
	}

	return def
}

//RemoveKey from data collection
func (m *Message) RemoveKey(key string) {
	m.Mtx.Lock()
	delete(m.Data, key)
	m.Mtx.Unlock()
}

//SourceID returns message source id
func (m *Message) SourceID() uint64 {
	return m.SID
}

//SetSourceID assign id to message
func (m *Message) SetSourceID(sourceID uint64) {
	m.SID = sourceID
}

//AppendRange append data fields to current Data
func (m *Message) AppendRange(data map[string]interface{}) {
	m.Mtx.Lock()
	for k, v := range data {
		m.Data[k] = v
	}
	m.Mtx.Unlock()
}

//GetTemperatureSensors data
func (m *Message) GetTemperatureSensors() *TemperatureSensors {
	return m.TemperatureSensors
}

//SetTemperatureSensors data
func (m *Message) SetTemperatureSensors(temperatureSensors *TemperatureSensors) {
	m.TemperatureSensors = temperatureSensors
}

//NewMessage returns new struct of  message
func NewMessage() *Message {
	return &Message{Data: make(map[string]interface{}), Mtx: &sync.Mutex{}}
}

//UnMarshalMessage given string to Message struct
func UnMarshalMessage(str string) (*Message, error) {
	message := &Message{Mtx: &sync.Mutex{}}
	err := json.Unmarshal([]byte(str), message)
	if err != nil {
		return &Message{Mtx: &sync.Mutex{}}, err
	}
	return message, err
}
