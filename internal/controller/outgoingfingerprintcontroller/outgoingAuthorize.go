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
	"strings"

	"github.com/tidwall/gjson"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/requestobjects"
)

func (controller *OutgoingFingerprintController) outgoingAuthorize(_reqObj requestobjects.SubmitForIdentifyReqObj) (isMatched bool, err error) {
	data := url.Values{}
	data.Set("client_id", os.Getenv("CONSUMER_KEY_FOR_OUTGOING"))
	data.Set("response_type", "code")
	data.Set("redirect_uri", os.Getenv("IDP_APP_REDIRECT_URI"))
	data.Set("state", "Authentication via fingerprint")
	data.Set("scope", "openid internal_login FingerprintAuth")
	data.Set("response_mode", "direct")

	requestBody := bytes.NewBufferString(data.Encode())
	urlString, err := url.JoinPath(controller.targetAdress, AuthorizeEndpoint)
	if err != nil {
		isMatched = false
		err = fmt.Errorf("error occured while trying to url string using url.JoinPath , %w", err)
		return
	}
	initialreq, err := http.NewRequest("POST", urlString, requestBody)
	if err != nil {
		isMatched = false
		err = fmt.Errorf("error occured while trying to make a new http post request : %w", err)
	}
	initialreq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	initialreq.Header.Add("Accept", "application/json")

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	internalclient := &http.Client{Transport: tr}
	res, errReqSend := internalclient.Do(initialreq)
	if errReqSend != nil {
		isMatched = false
		err = fmt.Errorf("error while sending or recieving post request : %w", errReqSend)
		return
	}
	defer res.Body.Close()
	bodybytes, errReadAll := io.ReadAll(res.Body)
	if errReadAll != nil {
		isMatched = false
		err = fmt.Errorf("error while reading bytes of response body : %w", errReadAll)
		return
	}
	if res.StatusCode != 200 {
		isMatched = false
		err = fmt.Errorf("the idp responded with an error message for first request for authorization : %s", string(bodybytes))
		return
	}
	initialResult := gjson.ParseBytes(bodybytes)

	authenticators := initialResult.Get("nextStep.authenticators").Array()
	authenticatorID := authenticators[0].Get("authenticatorId").String()
	flowID := initialResult.Get("flowId").String()
	biometricKey := base64.StdEncoding.EncodeToString(_reqObj.ProbeCbor)

	secondReqObj := requestobjects.AuthReqObj{
		FlowId: flowID,
		SelectedAuthenticar: requestobjects.AuthObj_SelectedAuthenticator{
			AuthenticationID: authenticatorID,
			Params: requestobjects.AuthObj_Params{
				BiometricKey: biometricKey,
			},
		}}
	jsonobj, err := json.Marshal(secondReqObj)
	if err != nil {
		isMatched = false
		err = fmt.Errorf("error occured while trying to convert reqObj to json, %w", err)
		return
	}
	secondRequestBody := bytes.NewBuffer(jsonobj)

	urlString, err = url.JoinPath(controller.targetAdress, AuthnEndpoint)
	if err != nil {
		isMatched = false
		err = fmt.Errorf("error occured while trying to make second url string using url.JoinPath , %w", err)
		return
	}

	secondReq, err := http.NewRequest("POST", urlString, secondRequestBody)
	if err != nil {
		isMatched = false
		err = fmt.Errorf("error occured while trying to make second http post request : %w", err)
	}

	secondReq.Header.Add("Content-Type", "application/json")

	secondRes, errReqSend := internalclient.Do(secondReq)
	if errReqSend != nil {
		isMatched = false
		err = fmt.Errorf("error while sending or recieving post request : %w", errReqSend)
		return
	}
	defer secondRes.Body.Close()

	secondBodybytes, errReadAll := io.ReadAll(secondRes.Body)
	if errReadAll != nil {
		isMatched = false
		err = fmt.Errorf("error while reading bytes of response body : %w", errReadAll)
		return
	}
	if res.StatusCode != 200 {
		isMatched = false
		err = fmt.Errorf("the idp responded with an error message for second request for authorization : %s", string(bodybytes))
		return
	}
	secondResult := gjson.ParseBytes(secondBodybytes)
	FlowStatus := secondResult.Get("flowStatus").String()
	var message string
	if FlowStatus == "SUCCESS_COMPLETED" {
		message = "Access Granted"
		code := secondResult.Get("authData.code").String()
		trData := url.Values{}
		trData.Set("grant_type", "authorization_code")
		trData.Set("code", code)
		trData.Set("redirect_uri", os.Getenv("IDP_APP_REDIRECT_URI"))

		urlString, err = url.JoinPath(controller.targetAdress, TokenEndpoint)
		if err != nil {
			isMatched = false
			err = fmt.Errorf("error occured while trying to make token endpoint url string using url.JoinPath , %w", err)
			return
		}
		requestBody = bytes.NewBufferString(trData.Encode())
		tokenReq, err := http.NewRequest("POST", urlString, requestBody)
		if err != nil {
			isMatched = false
			err = fmt.Errorf("error occured while trying to make a new http post request : %w", err)
		}
		consumerKey := os.Getenv("CONSUMER_KEY_FOR_OUTGOING")
		consumerSecret := os.Getenv("CONSUMER_SECRET_FOR_OUTGOING")
		authHeadervalue := base64.StdEncoding.EncodeToString([]byte(consumerKey + ":" + consumerSecret))
		tokenReq.Header.Add("Authorization", "Basic "+authHeadervalue)
		tokenReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		tokenRes, errReqSend := internalclient.Do(tokenReq)
		if errReqSend != nil {
			isMatched = false
			err = fmt.Errorf("error while sending or recieving post request : %w", errReqSend)
			// return
		}
		defer secondRes.Body.Close()

		tokenResBodyBytes, errReadAll := io.ReadAll(tokenRes.Body)
		if errReadAll != nil {
			isMatched = false
			err = fmt.Errorf("error while reading bytes of response body : %w", errReadAll)
			// return
		}
		if res.StatusCode != 200 {
			isMatched = false
			err = fmt.Errorf("the idp responded with an error message for second request for authorization : %s", string(bodybytes))
			// return
		}
		tokenReqResult := gjson.ParseBytes(tokenResBodyBytes)
		idToken := tokenReqResult.Get("id_token").String()
		tokenParts := strings.Split(idToken, ".")
		log.Println("token Parts len", len(tokenParts))
		payloadencoded := tokenParts[1]
		payloadDecoded, errDecode := base64.RawURLEncoding.DecodeString(payloadencoded)
		if errDecode != nil {
			err = fmt.Errorf("error while decoding jwt token payload : %w", errDecode)
		}
		log.Println("token payload : ", string(payloadDecoded))

		payloadResult := gjson.ParseBytes(payloadDecoded)
		log.Println(payloadResult.Get("DiscoveredID").String())
		log.Println(payloadResult.Get("IsMatched").String())
		isMatched = true

	} else {
		message = "Access Denied"
		isMatched = false
	}

	err = nil
	log.Println("message: ", message)
	return

}
