package controller

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"queclink-go/base.device.service/api/model"
	"github.com/gin-gonic/gin"
)

// SendCommandDirect godoc
// @Summary Send command to IP and port using UDP protocol
// @Description Send packet to IP:port
// @Tags device
// @Accept  json
// @Produce  json
// @Param ip formData string true "ip"
// @Param port formData int true "port"
// @Param command formData string true "command"
// @Success 200 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /device/directcommand [post]
func (c *Controller) SendCommandDirect(ctx *gin.Context) {
	ipAddress := ctx.Request.PostFormValue("ip")
	port := ctx.Request.PostFormValue("port")
	command := ctx.Request.PostFormValue("command")

	conn, err := net.Dial("udp", fmt.Sprintf("%v:%v", ipAddress, port))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "")
		return
	}	
	defer conn.Close()
	fmt.Fprintf(conn, command)
	response := &model.Response{
		CreatedAt:       time.Now().UTC(),
		Code:            fmt.Sprintf("Address:%v:%v; Command:%v", ipAddress, port, command),
		ExecutedCommand: command,
		Success:         true,
	}
	ctx.JSON(http.StatusOK, response)
}
