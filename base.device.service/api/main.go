package api

import (
	"fmt"

	"queclink-go/base.device.service/api/controller"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Device-service swagger example API
// @version 1.0
// @description This is device-service exaple API.
// @BasePath /

//Server struct
type Server struct {
	Server *gin.Engine
	Port   int
}

//StartNewAPIServer starts new http server and provides api access
func StartNewAPIServer(c controller.IApiController, port int, mode string) *Server {
	gin.SetMode(mode)
	r := gin.Default()
	pprof.Register(r)
	v1 := r.Group("/")
	{
		device := v1.Group("/device")
		{
			device.GET("identity_exists", c.DeviceIdentityExists)
			//
			device.GET("locate", c.GetLocateCommand)
			device.POST("locate", c.PostLocateCommand)
			//
			device.POST("command", c.SendCommand)
			//
			device.POST("update_config", c.PostConfiguration)

			device.GET("stats", c.GetServiceStats)
			//
			device.GET("immobilizer", c.GetImmobilizerCommand)
			device.POST("immobilizer", c.PostImmobilizerCommand)

			device.POST("directcommand", c.SendCommandDirect)
		}
		v1.GET("/debug/vars", controller.MetricsHandler)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	go r.Run(fmt.Sprintf(":%v", port))
	return &Server{
		Server: r,
		Port:   port,
	}

}
