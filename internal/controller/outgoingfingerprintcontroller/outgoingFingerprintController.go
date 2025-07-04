package outgoingfingerprintcontroller

import (
	"os"

	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/customstorage"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/requestobjects"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/responseobjects"
)

const (
	MatchTemplatesEndpoint    = "/api/fingerprint/match"
	IdentifyTemplateEndpoint  = "/api/fingerprint/identify"
	EnrollTemplateEndpoint    = "/api/fingerprint/enroll"
	UploadCborZipFileEndpoint = "/api/gallery/upload-cbor-zip"
	AuthorizeEndpoint         = "/oauth2/authorize"
	AuthnEndpoint             = "/oauth2/authn"
	TokenEndpoint             = "/oauth2/token"
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
func (controller *OutgoingFingerprintController) OutgoingIdentifyHandler(reqObj requestobjects.SubmitForIdentifyReqObj) (response responseobjects.IdentifyTemplateResObj, err error) {
	response, err = controller.outgoingIdentifyHandler(reqObj)
	return
}
func (controller *OutgoingFingerprintController) OutgoingMatchHandler(reqObj requestobjects.SubmitForMatchReqObj) (response responseobjects.MatchTemplatesResObj, err error) {
	response, err = controller.outgoingMatchHandler(reqObj)
	return
}
func (controller *OutgoingFingerprintController) OutgoingAuthorize(reqObj requestobjects.SubmitForIdentifyReqObj) (Status string, err error) {
	Status, err = controller.outgoingAuthorize(reqObj)
	return
}
