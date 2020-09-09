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

// GetImmobilizerCommand godoc
// @Summary Get immobilizer command
// @Description Returns immobilizer command
// @Tags device
// @Accept  json
// @Produce  json
// @Param identity query string true "identity"
// @Param port query string true "port"
// @Param state query string true "state"
// @Param trigger query string true "trigger"
// @Success 302 {object} model.Response
// @Failure 404 {object} model.Response
// @Router /device/immobilizer [get]
func (c *Controller) GetImmobilizerCommand(ctx *gin.Context) {
	identity := ctx.Query("identity")
	port := ctx.Query("port")
	state := ctx.Query("state")
	trigger := ctx.Query("trigger")
	method := ctx.Query("method")

	command := ""
	switch strings.ToLower(method) {
	case "sms_only":
		command = core.InstanceDM.GetSmsImmobilizerCommand(identity, port, state, trigger)
	default:
		command = core.InstanceDM.GetImmobilizerCommand(identity, port, state, trigger)
	}

	response := &model.Response{
		CreatedAt:       time.Now().UTC(),
		Code:            "Success",
		ExecutedCommand: command,
		Success:         true,
	}
	ctx.JSON(http.StatusFound, response)
}

// PostImmobilizerCommand godoc
// @Summary Send immobilizer state change request
// @Description Enqueue immobilizer state change request to device
// @Tags device
// @Accept  multipart/form-data
// @Produce  json
// @Param identity formData string true "identity"
// @Param callback_id formData string true "callback_id"
// @Param ttl formData int true "ttl"
// @Param port formData string true "port"
// @Param state formData string true "state"
// @Param trigger formData string true "trigger"
// @Param safety_option formData bool true "safety_option"
// @Success 302 {object} model.Response
// @Failure 404 {object} model.Response
// @Router /device/immobilizer [post]
func (c *Controller) PostImmobilizerCommand(ctx *gin.Context) {
	identity := ctx.Request.PostFormValue("identity")
	callbackID := ctx.Request.PostFormValue("callback_id")
	t := ctx.Request.PostFormValue("ttl")
	ttl, err := strconv.Atoi(t)
	if err != nil || callbackID == "" || identity == "" {
		response := &model.Response{
			CreatedAt: time.Now().UTC(),
			Success:   false,
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	port := ctx.Request.PostFormValue("port")
	state := ctx.Request.PostFormValue("state")
	trigger := ctx.Request.PostFormValue("trigger")
	safety := strings.ToUpper(ctx.Request.PostFormValue("safety_option")) == "TRUE"
	command := core.InstanceDM.GetImmobilizerCommand(identity, port, state, trigger)

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
	device.GetImmobilizer().SendAPIImmobilizerCommand(callbackID, port, state, ttl, trigger, safety)
	response.Success = true
	ctx.JSON(http.StatusFound, response)
}
