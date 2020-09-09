package controller

import (
	"math"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"queclink-go/base.device.service/api/model"
	"queclink-go/base.device.service/core"
	"queclink-go/base.device.service/utils"
)

// GetServiceStats godoc
// @Summary Get service statistics
// @Description Returns service statistics
// @Tags device
// @Accept  json
// @Produce  json
// @Success 200 {object} model.ServiceStatistics
// @Router /device/stats [get]
func (c *Controller) GetServiceStats(ctx *gin.Context) {
	p := &utils.ProcessInfo{}
	p.Pid = int32(os.Getpid())
	if percentage, err := p.Process.CPUPercent(); err == nil {
		percentage = percentage * 10
		percentage = math.Round(percentage) / 10
		p.CPUPercent = percentage
	}
	response := &model.ServiceStatistics{
		TotalDeviceCount:             core.InstanceDM.GetManagedConnections().Count(),
		TotalCountByWorkers:          core.InstanceDM.GetWorkers().DevicesCount(),
		UnregisteredConnectionsCount: core.InstanceDM.GetUnManagedConnections().Count(),
		UDPConnectionsCount:          core.InstanceDM.GetManagedConnections().GetTypedConnectionCount("UDP"),
		TCPConnectionsCount:          core.InstanceDM.GetManagedConnections().GetTypedConnectionCount("TCP"),
		ProcessInfo:                  p,
	}
	ctx.JSON(http.StatusOK, response)
}
