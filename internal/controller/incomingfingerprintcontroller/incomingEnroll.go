package incomingfingerprintcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/requestobjects"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/responseobjects"
)

func (controller *IncomingFingerprintController) incomingEnrollHandler(c *gin.Context) {
	var reqObj requestobjects.SubmitForEnrollReqObj
	err := c.BindJSON(&reqObj)
	if err != nil {
		resObj := responseobjects.ErrorResObj{Message: "Error when running BindJSON check response body contents, " + err.Error()}
		c.IndentedJSON(http.StatusInternalServerError, resObj)
		return
	}
	isAuthorized := false

	// checking through configured devices
	for _, deviceid := range controller.devicesConfig.EnrollmentDevices {
		if deviceid == reqObj.DeviceId {
			isAuthorized = true
		}
	}
	if !isAuthorized {
		resObj := responseobjects.ErrorResObj{Message: "Permission Denied Device that submit for Access is unregistered in orchestration server "}
		c.IndentedJSON(http.StatusUnauthorized, resObj)
		return
	}

	message, err := controller.outgoingfingerprintcontroller.OutgoingEnrollHandler(reqObj)
	if err != nil {
		resObj := responseobjects.ErrorResObj{Message: "error occured in outgoing to target: " + err.Error()}
		c.IndentedJSON(http.StatusInternalServerError, resObj)
		return
	}

	c.IndentedJSON(http.StatusOK, responseobjects.EnrollTemplateResObj{Message: message})
}
