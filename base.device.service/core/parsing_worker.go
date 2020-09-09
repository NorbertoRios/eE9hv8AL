package core

import (
	"log"
	"net"
	"sync"

	"queclink-go/base.device.service/comm"
	"queclink-go/base.device.service/report"
)

//WorkerData data
type WorkerData struct {
	Message report.IMessage
	DevID   string
	Server  *comm.UDPServer
	Addr    *net.UDPAddr
}

//InitializeParsingWorker new instance
func InitializeParsingWorker(manager IDeviceManager) *ParsingWorker {
	return &ParsingWorker{
		Devices:        make(map[string]string),
		Mutex:          &sync.Mutex{},
		MessageChannel: make(chan *WorkerData, 100000),
		Manager:        manager,
	}
}

//ParsingWorker workers
type ParsingWorker struct {
	Mutex          *sync.Mutex
	MessageChannel chan *WorkerData
	Devices        map[string]string
	Manager        IDeviceManager
}

//AddDevice add new devices
func (worker *ParsingWorker) AddDevice(devID string) {
	worker.Mutex.Lock()
	defer worker.Mutex.Unlock()
	worker.Devices[devID] = devID
}

//Count device count
func (worker *ParsingWorker) Count() int {
	worker.Mutex.Lock()
	defer worker.Mutex.Unlock()
	return len(worker.Devices)
}

//Start start prase worker
func (worker *ParsingWorker) Start() {
	go worker.process()
}

func (worker *ParsingWorker) process() {
	for {
		select {
		case data := <-worker.MessageChannel:
			worker.Mutex.Lock()
			_, found := worker.Devices[data.DevID]
			worker.Mutex.Unlock()
			if found {
				worker.processData(data)
			}
		}
	}
}

func (worker *ParsingWorker) processData(data *WorkerData) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in processData:", r)
		}
	}()
	d, found := worker.Manager.GetManagedConnections().GetDeviceByDeviceIdentity(data.DevID)
	var err error
	if !found {
		d, err = worker.Manager.InitializeUDPDevice(comm.InitializeUDPChannel(data.Server, data.Addr), data.Message)
		if err != nil {
			log.Println("[Worker | processData] Cant create device. Error: ", err)
			return
		}
	}
	ch, _ := d.GetChannel().(*comm.UDPChannel)
	ch.ClientAddr = data.Addr
	ch.Enqueue(data.Message)
	worker.Manager.Signalize(data.DevID)
}
