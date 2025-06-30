package incomingfingerprintcontroller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/requestobjects"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/responseobjects"
)

func (controller *IncomingFingerprintController) incomingEnrollHandler(c *gin.Context) {
	// check whether the client-id is authorized to be performing the enroll procedure
	authorizedClientId := os.Getenv("FINGERPRINT_AUTHORIZED_ENROLLMENT_CLIENT_ID")
	var reqObj requestobjects.SubmitForEnrollReqObj
	err := c.BindJSON(&reqObj)
	if err != nil {
		resObj := responseobjects.ErrorResObj{Message: "Error when running BindJSON check response body contents, " + err.Error()}
		c.IndentedJSON(http.StatusInternalServerError, resObj)
		return
	}
	if reqObj.ClientId != authorizedClientId {
		resObj := responseobjects.ErrorResObj{Message: "Permission Denied client that submit for enrollment is not authorized to use this feature "}
		c.IndentedJSON(http.StatusUnauthorized, resObj)
		return
	}

	message, err := controller.outgoingfingerprintcontroller.OutgoingEnrollHandler(reqObj)
	if err != nil {
		resObj := responseobjects.ErrorResObj{Message: "error occured in outgoing to target: " + err.Error()}
		c.IndentedJSON(http.StatusUnauthorized, resObj)
		return
	}
	c.IndentedJSON(http.StatusOK, message)
}
