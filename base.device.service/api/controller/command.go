package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"queclink-go/base.device.service/api/model"
	"queclink-go/base.device.service/core"
)

// SendCommand godoc
// @Summary Send command to device
// @Description Enqueue command to device
// @Tags device
// @Accept  json
// @Produce  json
// @Param identity formData string true "identity"
// @Param command formData string true "command"
// @Param callback_id formData string true "callback_id"
// @Param ttl formData int true "ttl"
// @Success 302 {object} model.Response
// @Failure 404 {object} model.Response
// @Router /device/command [post]
func (c *Controller) SendCommand(ctx *gin.Context) {
	identity := ctx.Request.PostFormValue("identity")
	t := ctx.Request.PostFormValue("ttl")
	callbackID := ctx.Request.PostFormValue("callback_id")
	command := ctx.Request.PostFormValue("command")

	_, err := strconv.Atoi(t)
	if err != nil || callbackID == "" || identity == "" || command == "" {
		response := &model.Response{
			CreatedAt: time.Now().UTC(),
			Success:   false,
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := &model.Response{
		CreatedAt:       time.Now().UTC(),
		Code:            command,
		ExecutedCommand: command,
		Success:         true,
	}

	device, found := core.InstanceDM.GetConectedDeviceByIdentity(identity)
	if !found {
		response.Code = fmt.Sprintf("Device with 'identity'=%v is offline", identity)
		response.Success = false
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	device.SendString(command)
	core.SendConfigurationAPIResponse(callbackID, "Done", true)
	ctx.JSON(http.StatusFound, response)
}
