package core

import (
	"container/list"
	"sync"
)

func BuildParsingWorkers(count int, manager IDeviceManager) *ParsingWorkers {
	workers := &ParsingWorkers{
		Workers: list.New(),
		mutex:   &sync.Mutex{},
		Devices: make(map[string]string, 0),
	}
	for i := 0; i < count; i++ {
		worker := InitializeParsingWorker(manager)
		workers.Workers.PushBack(worker)
		worker.Start()
	}
	return workers
}

//ParsingWorkers parse
type ParsingWorkers struct {
	Workers *list.List
	Devices map[string]string
	mutex   *sync.Mutex
}

//Signalize worker new message was received
func (w *ParsingWorkers) Signalize(data *WorkerData) {
	w.AddDevice(data.DevID)
	for wr := w.Workers.Front(); wr != nil; wr = wr.Next() {
		wr.Value.(*ParsingWorker).MessageChannel <- data
	}
}

//AddDevice add new device
func (w *ParsingWorkers) AddDevice(devID string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	_, f := w.Devices[devID]
	if f {
		return
	}
	worker := w.Workers.Front()
	value := worker.Value.(*ParsingWorker).Count()
	for wr := w.Workers.Front(); wr != nil; wr = wr.Next() {
		if wr.Value.(*ParsingWorker).Count() < value {
			worker = wr
			value = worker.Value.(*ParsingWorker).Count()
		}
	}
	worker.Value.(*ParsingWorker).AddDevice(devID)
	w.Devices[devID] = devID
}
