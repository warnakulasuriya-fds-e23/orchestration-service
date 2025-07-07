package outgoingfingerprintcontroller

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
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
	idTokenRetrivalUrl, errIdTokRetUrl := url.JoinPath(controller.targetAdress, TokenEndpoint)
	if errIdTokRetUrl != nil {
		err = fmt.Errorf("error encountered in authorization flow while at %s  : %w", Status, errIdTokRetUrl)
		return
	}

	// TODO: Initialize http client at outgoingFingerprintController Startup
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	internalclient := &http.Client{Transport: tr}

	Status = "Flow Initiation"
	intialResult, errflowItinit := authorizationcalls.FlowInitiationCall(flowInitiationUrl, _reqObj.ClientId, internalclient)
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
	if FlowStatus == "SUCCESS_COMPLETED" {
		Status = "Token Retrival"
		idToken, errIdToken := authorizationcalls.IdTokenRetrivalCall(idTokenRetrivalUrl, internalclient, secondResult)
		if errIdToken != nil {
			err = fmt.Errorf("error encountered in authorization flow while at %s  : %w", Status, errIdToken)
			return
		}

		Status = "Token Processing"
		tokenParts := strings.Split(idToken, ".")
		payloadencoded := tokenParts[1]
		payloadDecoded, errDecode := base64.RawURLEncoding.DecodeString(payloadencoded)
		if errDecode != nil {
			err = fmt.Errorf("error encountered in authorization flow while at %s while decoding jwt token payload : %w", Status, errDecode)
			return
		}

		Status = "Clearence Verification"
		payloadResult := gjson.ParseBytes(payloadDecoded)
		discoveredID := payloadResult.Get("DiscoveredID").String()
		rolesArray := payloadResult.Get("roles").Array()
		for _, role := range rolesArray {
			log.Println(role.String())
		}

		if discoveredID != "none" {
			Status = "Granting Access"
			// checkClearenceLevelAndUnlockDoor()
			Status = "Access Granted"
		} else {
			Status = "Access Denied"
		}

	} else {

		Status = "Identifaction Failed : " + FlowStatus
	}

	err = nil
	return

}
