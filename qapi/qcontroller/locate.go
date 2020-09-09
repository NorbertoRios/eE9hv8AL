package qcontroller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"queclink-go/base.device.service/api/model"
	"queclink-go/base.device.service/core"
)

// GetLocateCommand godoc
// @Summary Get locate command
// @Description Returns locate command
// @Tags device
// @Accept  json
// @Produce  json
// @Param identity query string true "identity"
// @Success 302 {object} model.Response
// @Failure 404 {object} model.Response
// @Router /device/locate [get]
func (c *QController) GetLocateCommand(ctx *gin.Context) {
	response := &model.Response{
		CreatedAt:       time.Now().UTC(),
		Code:            "AT+GTRTO=gv55,1,,,,,,FFFF$",
		ExecutedCommand: "AT+GTRTO=gv55,1,,,,,,FFFF$",
		Success:         true,
	}
	ctx.JSON(http.StatusFound, response)
}

// PostLocateCommand godoc
// @Summary Send locate request
// @Description Enqueue location request to device
// @Tags device
// @Accept  multipart/form-data
// @Produce  json
// @Param identity formData string true "identity"
// @Param callback_id formData string true "callback_id"
// @Param ttl formData int true "ttl"
// @Success 302 {object} model.Response
// @Failure 404 {object} model.Response
// @Router /device/locate [post]
func (c *QController) PostLocateCommand(ctx *gin.Context) {
	identity := ctx.Request.PostFormValue("identity")
	command := "AT+GTRTO=gv55,1,,,,,,FFFF$"
	t := ctx.Request.PostFormValue("ttl")
	callbackID := ctx.Request.PostFormValue("callback_id")

	ttl, err := strconv.Atoi(t)
	if err != nil || callbackID == "" || identity == "" {
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
		Success:         false,
	}

	device, found := core.InstanceDM.GetConectedDeviceByIdentity(identity)
	if !found {
		response.Code = fmt.Sprintf("Device with 'identity'=%v is offline", identity)
		response.Success = false
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	confItem := &core.ConfigurationItemAPI{
		CallbackID:   callbackID,
		CreationTime: time.Now().UTC(),
		TTL:          ttl,
	}
	confItem.OnCommandSent(core.SendConfigurationAPIResponse)
	confItem.Command = command
	confItem.MessageType = "16"
	confItem.Type = "api_request"

	device.GetConfiguration().AddCommand(confItem, true)
	device.GetConfiguration().SyncDeviceConfig()
	response.Success = true
	ctx.JSON(http.StatusFound, response)
}
