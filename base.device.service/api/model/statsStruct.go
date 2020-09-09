package model

import (
	"queclink-go/base.device.service/utils"
)

//ServiceStatistics describes count of connected devices
type ServiceStatistics struct {
	TotalDeviceCount             int
	UDPConnectionsCount          int
	TCPConnectionsCount          int
	UnregisteredConnectionsCount int
	TotalCountByWorkers          int
	ProcessInfo                  *utils.ProcessInfo
}
