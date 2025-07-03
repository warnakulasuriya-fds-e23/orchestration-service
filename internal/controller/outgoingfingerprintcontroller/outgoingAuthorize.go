package outgoingfingerprintcontroller

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/tidwall/gjson"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/requestobjects"
)

func (controller *OutgoingFingerprintController) outgoingAuthorize(_reqObj requestobjects.SubmitForIdentifyReqObj) (message string, err error) {
	data := url.Values{}
	data.Set("client_id", os.Getenv("CONSUMER_KEY_FOR_OUTGOING"))
	data.Set("response_type", "code")
	data.Set("redirect_uri", os.Getenv("IDP_APP_REDIRECT_URI"))
	data.Set("state", "Authentication via fingerprint")
	data.Set("scope", "openid internal_login")
	data.Set("response_mode", "direct")

	requestBody := bytes.NewBufferString(data.Encode())
	urlString, err := url.JoinPath(controller.targetAdress, AuthorizeEndpoint)
	if err != nil {
		message = ""
		err = fmt.Errorf("error occured while trying to url string using url.JoinPath , %w", err)
		return
	}
	req, err := http.NewRequest("POST", urlString, requestBody)
	if err != nil {
		message = ""
		err = fmt.Errorf("error occured while trying to make a new http post request : %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	internalclient := &http.Client{Transport: tr}
	res, errReqSend := internalclient.Do(req)
	if errReqSend != nil {
		message = ""
		err = fmt.Errorf("error while sending or recieving post request : %w", errReqSend)
		return
	}
	defer res.Body.Close()
	bodybytes, errReadAll := io.ReadAll(res.Body)
	if errReadAll != nil {
		message = ""
		err = fmt.Errorf("error while reading bytes of response body : %w", errReadAll)
		return
	}
	initialResult := gjson.ParseBytes(bodybytes)

	authenticators := initialResult.Get("nextStep.authenticators").Array()
	authenticatorID := authenticators[0].Get("authenticatorId").String()
	flowID := initialResult.Get("flowId").String()
	biometricKey := base64.StdEncoding.EncodeToString(_reqObj.ProbeCbor)

	secondReqObj := requestobjects.AuthReqObj{
		FlowId: flowID,
		SelectedAuthenticar: requestobjects.SelectedAuthenticator{
			AuthenticationID: authenticatorID,
			Params: requestobjects.Params{
				BiometricKey: biometricKey,
			},
		}}
	jsonobj, err := json.Marshal(secondReqObj)
	if err != nil {
		message = ""
		err = fmt.Errorf("error occured while trying to convert reqObj to json, %w", err)
		return
	}
	secondRequestBody := bytes.NewBuffer(jsonobj)

	urlString, err = url.JoinPath(controller.targetAdress, AuthnEndpont)
	if err != nil {
		message = ""
		err = fmt.Errorf("error occured while trying to make second url string using url.JoinPath , %w", err)
		return
	}

	secondReq, err := http.NewRequest("POST", urlString, secondRequestBody)
	if err != nil {
		message = ""
		err = fmt.Errorf("error occured while trying to make second http post request : %w", err)
	}

	secondReq.Header.Add("Content-Type", "application/json")

	secondRes, errReqSend := internalclient.Do(secondReq)
	if errReqSend != nil {
		message = ""
		err = fmt.Errorf("error while sending or recieving post request : %w", errReqSend)
		return
	}
	defer secondRes.Body.Close()

	secondBodybytes, errReadAll := io.ReadAll(secondRes.Body)
	if errReadAll != nil {
		message = ""
		err = fmt.Errorf("error while reading bytes of response body : %w", errReadAll)
		return
	}
	secondResult := gjson.ParseBytes(secondBodybytes)
	FlowStatus := secondResult.Get("flowStatus").String()
	if FlowStatus == "SUCCESS_COMPLETED" {
		message = "Access Granted"
	} else {
		message = "Access Denied"
	}
	err = nil
	log.Println("message: ", message)
	return

}
