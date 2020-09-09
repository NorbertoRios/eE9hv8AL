package controller

import (
	"expvar"
	"fmt"

	"queclink-go/base.device.service/core"
	"github.com/gin-gonic/gin"
)

var (
	serviceInfo = expvar.NewMap("Service")
)

//MetricsHandler handle metrics request
func MetricsHandler(c *gin.Context) {

	updateServiceInfo()
	w := c.Writer
	c.Header("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte("{\n"))
	first := true
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			w.Write([]byte(",\n"))
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	w.Write([]byte("\n}\n"))
	c.AbortWithStatus(200)
}

func updateServiceInfo() {
	serviceInfo.Set("ManagedConnections", MetricIntValue{V: core.InstanceDM.GetManagedConnections().Count()})
	serviceInfo.Set("TotalCountByWorkers", MetricIntValue{V: core.InstanceDM.GetWorkers().DevicesCount()})
	serviceInfo.Set("UnregisteredConnectionsCount", MetricIntValue{core.InstanceDM.GetUnManagedConnections().Count()})
	serviceInfo.Set("UDPConnectionsCount", MetricIntValue{core.InstanceDM.GetManagedConnections().GetTypedConnectionCount("UDP")})
	serviceInfo.Set("TCPConnectionsCount", MetricIntValue{core.InstanceDM.GetManagedConnections().GetTypedConnectionCount("TCP")})
	//serviceInfo.Set("OnUDPMessageTime", MetricStrValue{core.InstanceDM.TimeStamp()})
}

//MetricIntValue type adapter
type MetricIntValue struct {
	V int
}

func (m MetricIntValue) String() string {
	return fmt.Sprint(m.V)
}

//MetricStrValue type adapter
type MetricStrValue struct {
	V string
}

func (m MetricStrValue) String() string {
	return fmt.Sprint("\"", m.V, "\"")
}
