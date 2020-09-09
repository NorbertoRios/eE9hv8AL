package core

import (
	"container/list"
	"log"
	"sync"
	"time"

	"queclink-go/base.device.service/comm"
)

//Workers are list of workers
type Workers struct {
	Workers *list.List
}

//InitializeWorkers pool
func InitializeWorkers(count int) *Workers {
	workers := &Workers{
		Workers: list.New(),
	}
	for i := 0; i < count; i++ {
		worker := InitializeWorker()
		workers.Workers.PushBack(worker)
		worker.Start()
	}
	return workers
}

//AddDevice to less loaded worker
func (w *Workers) AddDevice(device IDevice) {
	w.RemoveDeviceByIdentity(device.GetIdentity())
	worker := w.Workers.Front()
	value := worker.Value.(*Worker).Devices.Count()
	for wr := w.Workers.Front(); wr != nil; wr = wr.Next() {
		if wr.Value.(*Worker).Devices.Count() < value {
			worker = wr
			value = worker.Value.(*Worker).Devices.Count()
		}
	}

	worker.Value.(*Worker).AddDevice(device)
}

//RemoveDevice from worker
func (w *Workers) RemoveDevice(device IDevice) {
	for wr := w.Workers.Front(); wr != nil; wr = wr.Next() {
		wr.Value.(*Worker).Remove(device)
	}
}

//RemoveDeviceByIdentity from worker
func (w *Workers) RemoveDeviceByIdentity(identity string) {
	for wr := w.Workers.Front(); wr != nil; wr = wr.Next() {
		wr.Value.(*Worker).RemoveDeviceByIdentity(identity)
	}
}

//DevicesCount from worker
func (w *Workers) DevicesCount() int {
	var cnt int
	for wr := w.Workers.Front(); wr != nil; wr = wr.Next() {
		cnt += wr.Value.(*Worker).DevicesCount()
	}
	return cnt
}

//Signalize worker new message was received
func (w *Workers) Signalize(identity string) {
	for wr := w.Workers.Front(); wr != nil; wr = wr.Next() {
		wr.Value.(*Worker).ProcessChannel <- identity
	}
}

//Worker devices thread
type Worker struct {
	Devices        *MangedConnections
	mutex          *sync.Mutex
	ProcessChannel chan string
}

//AddDevice to worker
func (w *Worker) AddDevice(device IDevice) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.Devices.Add(device)
}

//Remove from worker
func (w *Worker) Remove(device IDevice) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.Devices.Remove(device)
}

//RemoveDeviceByIdentity find and remove device by identity
func (w *Worker) RemoveDeviceByIdentity(identity string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.Devices.RemoveKey(identity)
}

//DevicesCount return count of devices assigned to worker
func (w *Worker) DevicesCount() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.Devices.Count()
}

func (w *Worker) processDevices() {

	for {
		w.process()
		time.Sleep(10 * time.Second)
	}
}

func (w *Worker) process() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in worker process:", r)
		}
	}()

	w.mutex.Lock()
	defer w.mutex.Unlock()
	connections := *w.Devices.GetConnections()
	for _, d := range connections {
		w.processDevice(d)
	}
}

func (w *Worker) processDevice(d IDevice) {
	if d.GetChannel().ConnectionMode() != "UDP" {
		return
	}
	d.GetChannel().(*comm.UDPChannel).Flushin()
}

func (w *Worker) waitForDevice() {
	for {
		select {
		case identity := <-w.ProcessChannel:
			w.mutex.Lock()
			if device, found := w.Devices.GetDeviceByDeviceIdentity(identity); found {
				w.processDevice(device)
			}
			w.mutex.Unlock()
		}
	}
}

//Start 2 goroutine. 1 backup thread. 2. For immediately processing
func (w *Worker) Start() {
	go w.processDevices()
	go w.waitForDevice()
}

//InitializeWorker new instance
func InitializeWorker() *Worker {
	return &Worker{
		Devices:        newManagedConnection(),
		mutex:          &sync.Mutex{},
		ProcessChannel: make(chan string),
	}
}
