package core

import (
	"log"
	"sync"
	"time"
)

//MangedConnections type to store all managed connections
type MangedConnections struct {
	connections map[string]IDevice
	mutex       *sync.Mutex
}

//Add new connection
func (c *MangedConnections) Add(device IDevice) {
	log.Println("New registered device :", device.GetIdentity())
	c.mutex.Lock()
	defer c.mutex.Unlock()
	d, found := c.connections[device.GetIdentity()]
	if found {
		d.GetChannel().Close()
		c.Remove(d)
	}
	c.connections[device.GetIdentity()] = device
}

//Remove connection
func (c *MangedConnections) Remove(device IDevice) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.connections, device.GetIdentity())
	log.Println("Removed connection for :", device.GetIdentity(), ";Total device count:", len(c.connections))
}

//RemoveKey removes connection by key
func (c *MangedConnections) RemoveKey(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.connections, key)
	log.Println("Removed connection :", key, ";Total device count:", len(c.connections))
}

//GetDeviceByDeviceIdentity returns device by device id
func (c *MangedConnections) GetDeviceByDeviceIdentity(identity string) (device IDevice, found bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	device, found = c.connections[identity]
	return
}

//RemoveConnections remove managed connections
func (c *MangedConnections) GarbageConnections(callback func(string)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, device := range c.connections {
		if time.Now().UTC().Sub(device.GetChannel().LastActivity()).Seconds() >= 3600 {
			device.GetChannel().Close()
			delete(c.connections, key)
			callback(key)
		}
	}
}

//LastActivity return last activity time
func (c *MangedConnections) LastActivity(key string) time.Time {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.connections[key].GetChannel().LastActivity()
}

//GetConnections returns connected devices
func (c *MangedConnections) GetConnections() *map[string]IDevice {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return &c.connections
}

//Count return count of connections
func (c *MangedConnections) Count() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.connections)
}

//GetTypedConnectionCount returns connections by count
func (c *MangedConnections) GetTypedConnectionCount(connectionMode string) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var cnt int
	for _, d := range c.connections {
		if d.GetChannel().ConnectionMode() == connectionMode {
			cnt++
		}
	}
	return cnt
}

//NewUnmanagedConnection returns ner struct of unmanagedConnection
func newManagedConnection() *MangedConnections {
	return &MangedConnections{
		connections: make(map[string]IDevice),
		mutex:       &sync.Mutex{},
	}
}
