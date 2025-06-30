package incomingfingerprintcontroller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/controller/outgoingfingerprintcontroller"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/requestobjects"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/responseobjects"
)

func IncomingEnrollHandler(c *gin.Context) {
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

	go outgoingfingerprintcontroller.OutgoingEnrollHandler(reqObj)

	c.IndentedJSON(http.StatusOK, "Successfully Forwarded data for enrollment")
}
