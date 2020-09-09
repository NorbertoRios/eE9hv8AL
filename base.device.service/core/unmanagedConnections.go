package core

import (
	"log"
	"sync"
	"time"

	"queclink-go/base.device.service/comm"
)

//UnmangedConnections type to stare all unmanaged connections
type UnmangedConnections struct {
	connections map[string]comm.IChannel
	mutex       *sync.Mutex
}

//GetConnections thread save connections
func (c *UnmangedConnections) GetConnections() *map[string]comm.IChannel {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return &c.connections
}

//RemoveConnections remove connections
func (c *UnmangedConnections) GarbageConnections() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, channel := range c.connections {
		if time.Now().UTC().Sub(channel.ConnectedAtTs()).Seconds() >= 60 {
			channel.Close()
			delete(c.connections, key)
		}
	}
}

//Add new connection
func (c *UnmangedConnections) Add(client comm.IChannel) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	strAddr := client.RemoteAddr()
	log.Println("ReceivedNewConnection from :", strAddr)
	c.connections[strAddr] = client
}

//Remove connection
func (c *UnmangedConnections) Remove(client comm.IChannel) {
	c.mutex.Lock()
	strAddr := client.RemoteAddr()
	delete(c.connections, strAddr)
	log.Println("Removed undefined connection from :", strAddr, ";Undefined count:", len(c.connections))
	c.mutex.Unlock()
}

//RemoveKey removes connection by key
func (c *UnmangedConnections) RemoveKey(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.connections, key)
	log.Println("Removed undefined connection from :", key, ";Undefined count:", len(c.connections))
}

//Count return count of connections
func (c *UnmangedConnections) Count() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.connections)
}

//NewUnmanagedConnection returns ner struct of unmanagedConnection
func newUnmanagedConnection() *UnmangedConnections {
	return &UnmangedConnections{
		connections: make(map[string]comm.IChannel),
		mutex:       &sync.Mutex{},
	}
}
