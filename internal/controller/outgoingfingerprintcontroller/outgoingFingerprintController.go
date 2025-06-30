package outgoingfingerprintcontroller

import (
	"os"

	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/customstorage"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/requestobjects"
)

const (
	MatchTemplatesEndpoint    = "/api/fingerprint/match"
	IdentifyTemplateEndpoint  = "/api/fingerprint/identify"
	EnrollTemplateEndpoint    = "/api/fingerprint/enroll"
	UploadCborZipFileEndpoint = "/api/gallery/upload-cbor-zip"
)

type OutgoingFingerprintController struct {
	tokenStorage customstorage.TokenStorage
	targetAdress string
}

func NewOutgoingFingerprintController(tstorage customstorage.TokenStorage) (controller *OutgoingFingerprintController) {
	return &OutgoingFingerprintController{
		tokenStorage: tstorage,
		targetAdress: os.Getenv("ADRESS_FOR_OUTGOING"),
	}
}

func (controller *OutgoingFingerprintController) OutgoingEnrollHandler(reqObj requestobjects.SubmitForEnrollReqObj) (message string, err error) {
	message, err = controller.outgoingEnrollHandler(reqObj)
	return
}
