package qconfig

import "queclink-go/base.device.service/config"

//QConfiguration of service
type QConfiguration struct {
	config.Configuration
	RTLThreshold int
}
