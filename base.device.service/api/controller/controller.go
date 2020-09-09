package controller

import "github.com/gin-gonic/gin"

//IApiController interface for service api controller
type IApiController interface {
	DeviceIdentityExists(ctx *gin.Context)
	GetLocateCommand(ctx *gin.Context)
	PostLocateCommand(ctx *gin.Context)
	SendCommand(ctx *gin.Context)
	PostConfiguration(ctx *gin.Context)
	GetServiceStats(ctx *gin.Context)
	GetImmobilizerCommand(ctx *gin.Context)
	PostImmobilizerCommand(ctx *gin.Context)
	SendCommandDirect(ctx *gin.Context)
}

//Controller for gin server
type Controller struct {
}

//NewController returns new instanse of Controller
func NewController() *Controller {
	return &Controller{}
}
