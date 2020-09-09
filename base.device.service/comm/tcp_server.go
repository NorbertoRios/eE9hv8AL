package comm

import (
	"fmt"
	"log"
	"net"
	"time"
)

//TCPServer struct
type TCPServer struct {
	Host                string
	Port                int
	onNewClientCallback func(c *Client)
	Listener            net.Listener
}

//Listen new connetcions
func (server *TCPServer) Listen() {

	defer server.Listener.Close()

	for {
		conn, _ := server.Listener.Accept()
		client := &Client{
			Connection:  conn,
			ConnectedAt: time.Now().UTC(),
		}
		server.onNewClientCallback(client)
		go client.Listen()
	}
}

//NewTCPServer creates new instance of tcp server
func NewTCPServer(host string, port int) *TCPServer {
	l, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
	log.Println("Creating server with address ", host, ":", port, "Error:", err)

	server := &TCPServer{
		Host:     host,
		Port:     port,
		Listener: l,
	}
	ServerCounters.AddFloat("Transmitted", 0)
	ServerCounters.AddFloat("Received", 0)
	server.OnNewClient(func(c *Client) {})
	return server
}

//OnNewClient indicates new client connected
func (server *TCPServer) OnNewClient(callback func(c *Client)) {
	server.onNewClientCallback = callback
}
