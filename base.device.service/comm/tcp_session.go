package comm

import (
	"fmt"
	"net"
	"time"
)

//Client struct for tcp connection
type Client struct {
	Connection     net.Conn
	ConnectedAt    time.Time
	onDisconnect   func(c IChannel)
	onPacket       func(c IChannel, packet []byte)
	LastActivityTs time.Time
	Received       int64
	Transmitted    int64
}

// Send text message to client
func (c *Client) Send(message string) error {
	count, err := c.Connection.Write([]byte(message))
	ServerCounters.AddFloat("Transmitted", float64(count))
	c.Transmitted += int64(count)
	return err
}

//SendBytes packet to client
func (c *Client) SendBytes(message []byte) error {
	count, err := c.Connection.Write(message)
	ServerCounters.AddFloat("Transmitted", float64(count))
	c.Transmitted += int64(count)
	return err
}

//Close current connection
func (c *Client) Close() {
	c.Connection.Close()
}

//ConnectedAtTs returns connected timestamp
func (c *Client) ConnectedAtTs() time.Time {
	return c.ConnectedAt
}

//ReceivedBytes total received bytes
func (c *Client) ReceivedBytes() int64 {
	return c.Received
}

//TransmittedBytes total transmitted bytes
func (c *Client) TransmittedBytes() int64 {
	return c.Transmitted
}

//RemoteAddr client's ip address
func (c *Client) RemoteAddr() string {
	return c.Connection.RemoteAddr().String()
}

//RemoteIP indicates device remote IP address
func (c *Client) RemoteIP() string {
	return fmt.Sprintf("%v", c.Connection.RemoteAddr().String())
}

//RemotePort indicates device remote port
func (c *Client) RemotePort() int {
	return 0
}

//Listen client data from channel
func (c *Client) Listen() {
	buffer := make([]byte, 4096)
	for {
		count, err := c.Connection.Read(buffer)
		if err != nil {
			c.Connection.Close()
			c.onDisconnect(c)
			return
		}
		c.LastActivityTs = time.Now().UTC()
		c.Received = c.Received + int64(count)
		ServerCounters.AddFloat("Received", float64(count))
		c.onPacket(c, buffer[:count])
	}
}

//OnDisconnect indicates client disconnected
func (c *Client) OnDisconnect(callback func(c IChannel)) {
	c.onDisconnect = callback
}

//OnPacket indicates client reseived new packet
func (c *Client) OnPacket(callback func(c IChannel, packet []byte)) {
	c.onPacket = callback
}

//ConnectionMode indicates connection mode
func (c *Client) ConnectionMode() string {
	return "TCP"
}

//LastActivity indicates last device activity
func (c *Client) LastActivity() time.Time {
	return c.LastActivityTs
}
