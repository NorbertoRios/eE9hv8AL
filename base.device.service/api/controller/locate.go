package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"queclink-go/base.device.service/api/model"
	"queclink-go/base.device.service/core"
	"github.com/gin-gonic/gin"
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
func (c *Controller) GetLocateCommand(ctx *gin.Context) {
	uniqueID := strings.Replace(ctx.Query("identity"), "xirgo_", "", -1)
	response := &model.Response{
		CreatedAt:       time.Now().UTC(),
		Code:            fmt.Sprintf("+XT:%v,7001,1", uniqueID),
		ExecutedCommand: fmt.Sprintf("+XT:%v,7001,1", uniqueID),
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
func (c *Controller) PostLocateCommand(ctx *gin.Context) {
	identity := ctx.Request.PostFormValue("identity")
	uniqueID := strings.Replace(identity, "xirgo_", "", -1)
	command := fmt.Sprintf("+XT:%v,7001,1", uniqueID)
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
	confItem.MessageType = "7001"
	confItem.Type = "api_request"

	device.GetConfiguration().AddCommand(confItem, true)
	device.GetConfiguration().SyncDeviceConfig()
	response.Success = true
	ctx.JSON(http.StatusFound, response)
}
