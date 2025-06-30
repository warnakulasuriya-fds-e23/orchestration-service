package outgoingfingerprintcontroller

import (
	"os"

	"github.com/warnakulasuriya-fds-e23/orchestration-service/requestobjects"
)

func OutgoingEnrollHandler(reqObj requestobjects.SubmitForEnrollReqObj) {
	os.Getenv("ADRESS_FOR_OUTGOING")
}
