package comm

import (
	"time"
)

//IChannel connection client interface
type IChannel interface {
	Send(message string) error
	SendBytes(message []byte) error
	OnDisconnect(callback func(c IChannel))
	OnPacket(callback func(c IChannel, packet []byte))
	Close()
	RemoteAddr() string
	RemoteIP() string
	RemotePort() int
	ConnectedAtTs() time.Time
	LastActivity() time.Time
	ReceivedBytes() int64
	TransmittedBytes() int64
	ConnectionMode() string
}
