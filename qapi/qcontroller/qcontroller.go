package qcontroller

import "queclink-go/base.device.service/api/controller"

//QController for gin server
type QController struct {
	controller.Controller
}

//NewController returns new instanse of QController
func NewController() *QController {
	return &QController{}
}
