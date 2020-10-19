package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"queclink-go/base.device.service/comm"
	"queclink-go/base.device.service/config"
	"queclink-go/base.device.service/report"
	"queclink-go/base.device.service/utils"
)

//InstanceDM of device manager
var InstanceDM IDeviceManager

//IDeviceManager interface for devices
type IDeviceManager interface {
	Initialize()
	Stop()
	ReceivedNewConnection(client *comm.Client)
	AddRegisteredDevice(device IDevice)
	GetConectedDeviceByIdentity(identity string) (device IDevice, found bool)
	OnUDPPacket(server *comm.UDPServer, addr *net.UDPAddr, packet []byte)
	InitializeDevice(c comm.IChannel, message report.IMessage) (IDevice, error)
	InitializeUDPDevice(c comm.IChannel, message report.IMessage) (IDevice, error)
	GetManagedConnections() *MangedConnections
	GetUnManagedConnections() *UnmangedConnections
	GetWorkers() *Workers
	SetParser(parser report.IParser)
	GetParser() report.IParser
	GetImmobilizerCommand(identity string, port string, state string, trigger string) string
	GetSmsImmobilizerCommand(identity string, port string, state string, trigger string) string
	DisconnectUnmanagedClient(c comm.IChannel)
	DisconnectManagedClient(device IDevice)
	Signalize(string)
}

//DeviceManager implementation of IDeviceManager
type DeviceManager struct {
	PrasingWorkers              *ParsingWorkers
	UnmanagedConnections        *UnmangedConnections
	ManagedConnections          *MangedConnections
	DeviceWorkers               *Workers
	quit                        chan struct{}
	Parser                      report.IParser
	InitializeDeviceCallback    func(c comm.IChannel, message report.IMessage) (IDevice, error)
	InitializeUDPDeviceCallback func(c comm.IChannel, message report.IMessage) (IDevice, error)
	OnUDPDuration               string
	mutex                       *sync.Mutex
}

func (manager *DeviceManager) InitializeDevice(c comm.IChannel, message report.IMessage) (IDevice, error) {
	return manager.InitializeDeviceCallback(c, message)
}
func (manager *DeviceManager) InitializeUDPDevice(c comm.IChannel, message report.IMessage) (IDevice, error) {
	return manager.InitializeUDPDeviceCallback(c, message)
}

//ReceiveUnmanagedPacket receive packet from unmanaged connection. Must be implemented in derived struct
func (manager *DeviceManager) ReceiveUnmanagedPacket(c comm.IChannel, packet []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in ReceiveUnmanagedPacket:", r)
		}
	}()

	strAddr := c.RemoteAddr()
	log.Println("Received packet:", utils.InsertNth(utils.ByteToString(packet), 2, ' '), "from:", strAddr)
	messages, err := manager.Parser.Parse(packet)
	if err != nil {
		return
	}

	var device IDevice

	if len(messages) > 0 {
		var derr error
		device, derr = manager.InitializeDevice(c, messages[0])
		if derr != nil {
			c.Close()
			manager.DisconnectUnmanagedClient(c)
			return
		}
	}
	for _, message := range messages {
		ack, found := message.GetValue("Ack")
		device.ProcessMessage(message)
		if found {
			c.SendBytes(ack.([]byte))
		}
	}
	jMessage, jerr := json.Marshal(messages)
	if jerr == nil {
		log.Println("Received packet:", string(jMessage), "from:", strAddr)
	}
}

func (manager *DeviceManager) OnUDPPacket(server *comm.UDPServer, addr *net.UDPAddr, packet []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in OnUDPPacket:", r)
		}
	}()
	ts := time.Now()
	log.Println("Received UDP packet:", utils.InsertNth(utils.ByteToString(packet), 2, ' '), "from:", addr.String())
	messages, err := manager.Parser.Parse(packet)
	if err != nil || messages == nil || len(messages) == 0 {
		ack := manager.Parser.GetUnknownAck(packet)
		if ack != nil && len(ack) > 0 {
			server.SendBytes(addr, ack)
		}
		return
	}
	for _, message := range messages {
		devID, f := message.GetValue("DevId")
		if !f {
			continue
		}
		data := &WorkerData{
			Addr:    addr,
			Server:  server,
			Message: message,
			DevID:   devID.(string),
		}
		// ack, found := message.GetValue("Ack")
		// if found {
		// 	server.SendBytes(addr, ack.([]byte))
		// }
		manager.PrasingWorkers.Signalize(data)
	}
	log.Println("[OnUDPPacket] Processing time : ", fmt.Sprint(time.Since(ts)))
}

func (manager *DeviceManager) Signalize(devID string) {
	manager.DeviceWorkers.Signalize(devID)
}

//SetParser assing parser for reports
func (manager *DeviceManager) SetParser(parser report.IParser) {
	manager.Parser = parser
}

//GetParser returns parser for reports
func (manager *DeviceManager) GetParser() report.IParser {
	return manager.Parser
}

//GetManagedConnections returns managed connections
func (manager *DeviceManager) GetManagedConnections() *MangedConnections {
	return manager.ManagedConnections
}

//GetUnManagedConnections returns unmanaged connections
func (manager *DeviceManager) GetUnManagedConnections() *UnmangedConnections {
	return manager.UnmanagedConnections
}

//GetWorkers returns workers
func (manager *DeviceManager) GetWorkers() *Workers {
	return manager.DeviceWorkers
}

//GetImmobilizerCommand return immobilizer command for network usage
func (manager *DeviceManager) GetImmobilizerCommand(identity string, port string, state string, trigger string) string {
	return ""
}

//GetSmsImmobilizerCommand return immobilizer command for cell usage
func (manager *DeviceManager) GetSmsImmobilizerCommand(identity string, port string, state string, trigger string) string {
	return ""
}

//Initialize device manager
func (manager *DeviceManager) Initialize() {
	manager.mutex = &sync.Mutex{}
	manager.UnmanagedConnections = newUnmanagedConnection()
	manager.ManagedConnections = newManagedConnection()
	manager.DeviceWorkers = InitializeWorkers(config.Config.GetBase().WorkersCount)
	manager.PrasingWorkers = BuildParsingWorkers(config.Config.GetBase().ParsingWorkersCount, manager)
	go func() {
		managedticker := time.NewTicker(300 * time.Second)
		ticker := time.NewTicker(5 * time.Second)
		manager.quit = make(chan struct{})
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in watchdog function:", r)
			}
		}()
		for {
			select {
			case <-ticker.C:
				manager.GetUnManagedConnections().GarbageConnections()
				time.Sleep(2 * time.Second)
			case <-managedticker.C:
				manager.GetManagedConnections().GarbageConnections(manager.DeviceWorkers.RemoveDeviceByIdentity)
			case <-manager.quit:
				ticker.Stop()
				managedticker.Stop()
				return
			}
		}
	}()
}

//Stop stop cleaner timer for device manager
func (manager *DeviceManager) Stop() {
	close(manager.quit)
}

//ReceivedNewConnection handle new tcp connection
func (manager *DeviceManager) ReceivedNewConnection(client *comm.Client) {
	client.OnDisconnect(manager.DisconnectUnmanagedClient)
	client.OnPacket(manager.ReceiveUnmanagedPacket)
	manager.UnmanagedConnections.Add(client)
}

//DisconnectUnmanagedClient finalizes managed device
func (manager *DeviceManager) DisconnectUnmanagedClient(c comm.IChannel) {
	manager.UnmanagedConnections.Remove(c)
}

//DisconnectManagedClient finalizes managed device
func (manager *DeviceManager) DisconnectManagedClient(device IDevice) {
	device.GetChannel().Close()
	manager.ManagedConnections.Remove(device)
}

//AddRegisteredDevice to managed device collection
func (manager *DeviceManager) AddRegisteredDevice(device IDevice) {
	device.GetChannel().OnPacket(device.OnReceivePacket)
	device.GetChannel().OnDisconnect(device.OnChannelDisconnected)
	device.SetOnDisconnect(manager.DisconnectManagedClient)
	manager.DeviceWorkers.RemoveDeviceByIdentity(device.GetIdentity())
	manager.UnmanagedConnections.Remove(device.GetChannel())
	manager.ManagedConnections.Add(device)
}

//GetConectedDeviceByIdentity returns device by device identity from managed connections
func (manager *DeviceManager) GetConectedDeviceByIdentity(identity string) (device IDevice, found bool) {
	return manager.ManagedConnections.GetDeviceByDeviceIdentity(identity)
}

//InitializeDeviceManager initialization and assigment instance to single InstanceDM
func InitializeDeviceManager(instance IDeviceManager) {
	InstanceDM = instance
	instance.Initialize()
}
