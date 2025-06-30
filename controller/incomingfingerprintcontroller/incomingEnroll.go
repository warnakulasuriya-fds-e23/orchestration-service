package incomingfingerprintcontroller

import (
	"os"

	"github.com/gin-gonic/gin"
)

func IncomingEnrollHandler(c *gin.Context) {
	// check whether the client-id is authorized to be performing the enroll procedure
	os.Getenv("FINGERPRINT_AUTHORIZED_ENROLLMENT_CLIENT_ID")
}
