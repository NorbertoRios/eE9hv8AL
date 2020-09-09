package controller

import (
	"fmt"
	"net/http"
	"time"

	"queclink-go/base.device.service/api/model"
	"queclink-go/base.device.service/core"
	"github.com/gin-gonic/gin"
)

// PostConfiguration godoc
// @Summary Send configuration to device
// @Description Enqueue configuration to device
// @Tags device
// @Accept  multipart/form-data
// @Produce  json
// @Param identity formData string true "identity"
// @Param source formData int true "source"
// @Param config formData array true "trigger"
// @Param safety_option formData bool true "safety_option"
// @Success 200 {object} model.Response
// @Failure 404 {object} model.Response
// @Router /device/update_config [post]
func (c *Controller) PostConfiguration(ctx *gin.Context) {
	identity := ctx.Request.PostFormValue("identity")
	response := &model.Response{
		CreatedAt: time.Now().UTC(),
		Code:      "",
		Success:   false,
	}

	device, found := core.InstanceDM.GetConectedDeviceByIdentity(identity)
	if !found {
		response.Code = fmt.Sprintf("Cant send config  to %v", identity)
		response.Success = false
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	device.GetConfiguration().Load()
	device.GetConfiguration().SyncDeviceConfig()
	response.Success = true
	ctx.JSON(http.StatusFound, response)
}
