package incomingfingerprintcontroller

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/requestobjects"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/responseobjects"
)

func (controller *IncomingFingerprintController) incomingMatchHandler(c *gin.Context) {
	authorizedClientIdsString := os.Getenv("FINGERPRINT_AUTHORIZED_GENERAL_PURPOSE_CLIENT_IDS")
	authorizedClientIds := strings.Split(authorizedClientIdsString, ",")
	var reqObj requestobjects.SubmitForMatchReqObj
	err := c.BindJSON(&reqObj)
	if err != nil {
		resObj := responseobjects.ErrorResObj{Message: "Error when running BindJSON check response body contents, " + err.Error()}
		c.IndentedJSON(http.StatusInternalServerError, resObj)
		return
	}
	isAuthorized := false

	for _, authorizedClientId := range authorizedClientIds {
		if authorizedClientId == reqObj.DeviceId {
			isAuthorized = true
		}
	}
	if !isAuthorized {
		resObj := responseobjects.ErrorResObj{Message: "Permission Denied client that submit for enrollment is not authorized to use this feature "}
		c.IndentedJSON(http.StatusUnauthorized, resObj)
		return
	}

	response, err := controller.outgoingfingerprintcontroller.OutgoingMatchHandler(reqObj)
	if err != nil {
		resObj := responseobjects.ErrorResObj{Message: "error occured in outgoing to target: " + err.Error()}
		c.IndentedJSON(http.StatusInternalServerError, resObj)
		return
	}
	c.IndentedJSON(http.StatusOK, response)

}
