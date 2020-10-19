package comm

import (
	"container/list"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"queclink-go/base.device.service/report"
)

//UDPChannel is the type of IChannel
type UDPChannel struct {
	Messages         *list.List
	ConnectedAt      time.Time
	onDisconnect     func(c IChannel)
	onPacket         func(c IChannel, packet []byte)
	onProcessMessage func(message report.IMessage) error
	LastActivityTs   time.Time
	Received         int64
	Transmitted      int64
	Server           *UDPServer
	ClientAddr       *net.UDPAddr
	mutex            *sync.Mutex
}

//InitializeUDPChannel returns new udp channel
func InitializeUDPChannel(server *UDPServer, addr *net.UDPAddr) *UDPChannel {
	return &UDPChannel{
		Messages:    list.New(),
		ConnectedAt: time.Now().UTC(),
		Server:      server,
		ClientAddr:  addr,
		mutex:       &sync.Mutex{},
	}
}

// Send text message to client
func (c *UDPChannel) Send(message string) error {
	count, err := c.Server.Send(c.ClientAddr, message)
	c.Transmitted += int64(count)
	log.Println("[UDPChannel] Message: ", message, "Sended to ", c.ClientAddr.String())
	return err
}

//SendBytes packet to client
func (c *UDPChannel) SendBytes(message []byte) error {
	count, err := c.Server.SendBytes(c.ClientAddr, message)
	c.Transmitted += int64(count)
	log.Println("[SendBytes] ", string(message), "Sended to ", c.ClientAddr.String())
	return err
}

//Enqueue packet to queue
func (c *UDPChannel) Enqueue(message report.IMessage) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Messages.PushBack(message)
	c.LastActivityTs = time.Now().UTC()
}

//Flushin messages to device
func (c *UDPChannel) Flushin() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	count := c.Messages.Len()
	for count > 0 {
		imessage := c.Messages.Front()
		message, valid := imessage.Value.(report.IMessage)
		if valid {
			err := c.onProcessMessage(message)
			if err != nil {
				return
			}
			c.Messages.Remove(imessage)
			count = c.Messages.Len()
		}

	}
}

//Close method closes current channel
func (c *UDPChannel) Close() {

}

//ConnectedAtTs indicates create date of channel
func (c *UDPChannel) ConnectedAtTs() time.Time {
	return c.ConnectedAt
}

//ReceivedBytes indicates count of received bytes
func (c *UDPChannel) ReceivedBytes() int64 {
	return c.Received
}

//TransmittedBytes indicates count of transmitted bytes
func (c *UDPChannel) TransmittedBytes() int64 {
	return c.Transmitted
}

//RemoteAddr indicates device remote address
func (c *UDPChannel) RemoteAddr() string {
	return c.ClientAddr.String()
}

//RemoteIP indicates device remote IP address
func (c *UDPChannel) RemoteIP() string {
	return fmt.Sprintf("%v", c.ClientAddr.IP)
}

//RemotePort indicates device remote port
func (c *UDPChannel) RemotePort() int {
	return c.ClientAddr.Port
}

//OnDisconnect indicates client disconnected
func (c *UDPChannel) OnDisconnect(callback func(c IChannel)) {
	c.onDisconnect = callback
}

//OnPacket indicates client received new packet
func (c *UDPChannel) OnPacket(callback func(c IChannel, packet []byte)) {
	c.onPacket = callback
}

//OnProcessMessage assign callback
func (c *UDPChannel) OnProcessMessage(callback func(message report.IMessage) error) {
	c.onProcessMessage = callback
}

//ConnectionMode indicates connection mode
func (c *UDPChannel) ConnectionMode() string {
	return "UDP"
}

//LastActivity indicates last device activity
func (c *UDPChannel) LastActivity() time.Time {
	return c.LastActivityTs
}
