package outgoingfingerprintcontroller

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/controller/outgoingfingerprintcontroller/authorizationcalls"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/requestobjects"
)

func (controller *OutgoingFingerprintController) outgoingAuthorize(_reqObj requestobjects.SubmitForIdentifyReqObj) (Status string, err error) {
	Status = "processing Urls"
	//TODO : Handle Url concatanation at outgoingFingerprintController Startup
	flowInitiationUrl, errflInitUrl := url.JoinPath(controller.targetAdress, AuthorizeEndpoint)
	if errflInitUrl != nil {
		err = fmt.Errorf("error encountered in authorization flow while at %s  : %w", Status, errflInitUrl)
		return
	}
	credentialSubmissionUrl, errCredSubmitUrl := url.JoinPath(controller.targetAdress, AuthnEndpoint)
	if errCredSubmitUrl != nil {
		err = fmt.Errorf("error encountered in authorization flow while at %s  : %w", Status, errCredSubmitUrl)
		return
	}
	// TODO: Initialize http client at outgoingFingerprintController Startup
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	internalclient := &http.Client{Transport: tr}

	Status = "Flow Initiation"
	intialResult, errflowItinit := authorizationcalls.FlowInitiationCall(flowInitiationUrl, _reqObj.DeviceId, internalclient)
	if errflowItinit != nil {
		err = fmt.Errorf("error encountered in authorization flow while at %s  : %w", Status, errflowItinit)
		return
	}

	Status = "Credential Submission"
	secondResult, errCredSubmit := authorizationcalls.CredentialSubmissionCall(credentialSubmissionUrl, internalclient, &intialResult, _reqObj)
	if errCredSubmit != nil {
		err = fmt.Errorf("error encountered in authorization flow while at %s  : %w", Status, errCredSubmit)
		return
	}

	FlowStatus := secondResult.Get("flowStatus").String()
	switch FlowStatus {
	case "SUCCESS_COMPLETED":
		deviceDetails, ok := controller.devicesConfig.DeviceDetails[_reqObj.DeviceId]
		if !ok {
			Status = "device was identified as unregistered within orchestration server"
			err = fmt.Errorf("device id is not registered")
			return
		}
		door := deviceDetails.Door
		floor := deviceDetails.Floor
		Status = fmt.Sprintf("Access Granted for floor: %s door: %s", floor, door)
		err = nil
		return
	case "INCOMPLETE":
		deviceDetails, ok := controller.devicesConfig.DeviceDetails[_reqObj.DeviceId]
		if !ok {
			Status = "device was identified as unregistered within orchestration server"
			err = fmt.Errorf("device id is not registered")
			return
		}
		door := deviceDetails.Door
		floor := deviceDetails.Floor
		Status = fmt.Sprintf("Access Denied for floor: %s door: %s", floor, door)
		err = nil
		return
	case "":
		Status = "unregistered biometric data"
		err = nil
		return
	default:

		Status = "Identifaction Failed : " + FlowStatus
		errorMessage := "device is not properly cofigured in idp" + secondResult.Get("message").String()
		err = fmt.Errorf("%s", errorMessage)
		return
	}

}
