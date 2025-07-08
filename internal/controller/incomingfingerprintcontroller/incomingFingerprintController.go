package incomingfingerprintcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/config"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/controller/outgoingfingerprintcontroller"
)

type IncomingFingerprintController struct {
	outgoingfingerprintcontroller *outgoingfingerprintcontroller.OutgoingFingerprintController
	devicesConfig                 config.DevicesConfigJSON
}

func NewIncomingFingerprintController(devicesConfig config.DevicesConfigJSON, outcntrl *outgoingfingerprintcontroller.OutgoingFingerprintController) *IncomingFingerprintController {
	return &IncomingFingerprintController{
		outgoingfingerprintcontroller: outcntrl,
		devicesConfig:                 devicesConfig,
	}
}

func (controller *IncomingFingerprintController) IncomingIdentifyHandler(c *gin.Context) {
	// controller.incomingIdentifyHandler(c)
	controller.IncomingAuthorize(c)
}
func (controller *IncomingFingerprintController) IncomingMatchHandler(c *gin.Context) {
	controller.incomingMatchHandler(c)
}
func (controller *IncomingFingerprintController) IncomingEnrollHandler(c *gin.Context) {
	controller.incomingEnrollHandler(c)
}
func (controller *IncomingFingerprintController) IncomingAuthorize(c *gin.Context) {
	controller.incomingAuthorize(c)
}
