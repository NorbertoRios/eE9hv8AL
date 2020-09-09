package controller

import (
	"fmt"
	"net/http"
	"time"

	"queclink-go/base.device.service/api/model"
	"queclink-go/base.device.service/core"
	"github.com/gin-gonic/gin"
)

// DeviceIdentityExists godoc
// @Summary Checks device is currently connected to service
// @Description Checks device by device identity
// @Tags device
// @Accept  json
// @Produce  json
// @Param identity query string true "identity"
// @Success 302 {object} model.Response
// @Failure 404 {object} model.Response
// @Router /device/identity_exists [get]
func (c *Controller) DeviceIdentityExists(ctx *gin.Context) {
	identity := ctx.Query("identity")
	if identity == "" {
		response := &model.Response{
			CreatedAt: time.Now().UTC(),
			Success:   false,
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	_, found := core.InstanceDM.GetConectedDeviceByIdentity(identity)

	response := &model.Response{
		CreatedAt: time.Now().UTC(),
	}

	if found {
		response.Code = fmt.Sprintf("Device with 'identity'=%v online", identity)
		response.Success = true
		ctx.JSON(http.StatusFound, response)
		return
	}
	response.Code = fmt.Sprintf("Device with 'identity'=%v offline", identity)
	response.Success = false
	ctx.JSON(http.StatusNotFound, response)
}
