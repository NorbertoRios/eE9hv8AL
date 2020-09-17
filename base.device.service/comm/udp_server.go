package comm

import (
	"fmt"
	"log"
	"net"
)

//UDPServer implementation of server
type UDPServer struct {
	Port       int
	Connection *net.UDPConn
	onPacket   func(server *UDPServer, addr *net.UDPAddr, packet []byte)
}

//Listen for incoming packet
func (server *UDPServer) Listen() {
	for {
		var buf [4096]byte
		n, addr, err := server.Connection.ReadFromUDP(buf[0:])
		if err != nil {
			log.Fatalf("Error Reading from udp connection:%v", err.Error())
			return
		}		
		ServerCounters.AddFloat("Received", float64(n))
		server.onPacket(server, addr, buf[0:n])
	}
}

//OnPacket assing callback for incoming packets
func (server *UDPServer) OnPacket(callback func(*UDPServer, *net.UDPAddr, []byte)) {
	server.onPacket = callback
}

//SendBytes bytes via udp
func (server *UDPServer) SendBytes(addr *net.UDPAddr, packet []byte) (int, error) {
	ServerCounters.AddFloat("Transmitted", float64(len(packet)))
	return server.Connection.WriteToUDP(packet, addr)
}

//Send string via udp
func (server *UDPServer) Send(addr *net.UDPAddr, packet string) (int, error) {
	ServerCounters.AddFloat("Transmitted", float64(len(packet)))
	return server.Connection.WriteToUDP([]byte(packet), addr)
}

//NewUDPServer creates new UDP server
func NewUDPServer(host string, port int) *UDPServer {
	addr := fmt.Sprintf("%v:%v", host, port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("Wrong UDP Address:%v", addr)
		return nil
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		log.Fatalf("Create udp server error:%v", err.Error())
		return nil
	}
	udpServer := &UDPServer{
		Port:       port,
		Connection: udpConn,
	}
	ServerCounters.AddFloat("Transmitted", 0)
	ServerCounters.AddFloat("Received", 0)
	log.Println("UDP server created on:", addr)
	return udpServer
}
