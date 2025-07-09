package authorizationcalls

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/requestobjects"
)

func CredentialSubmissionCall(accessToken string, urlString string, internalClient *http.Client, initialResult *gjson.Result, _reqObj requestobjects.SubmitForIdentifyReqObj) (secondResult gjson.Result, err error) {
	authenticators := initialResult.Get("nextStep.authenticators").Array()
	var authenticatorId string
	for _, authenticator := range authenticators {
		if authenticator.Get("authenticator").String() == os.Getenv("IDP_BIO_SDK_AUTHENTICATOR_DISPLAY_NAME") {
			authenticatorId = authenticator.Get("authenticatorId").String()
		}
	}

	if authenticatorId == "" {
		err = fmt.Errorf("error occured authenticator under given displayname was not found")
		return
	}
	flowID := initialResult.Get("flowId").String()
	BiometricTemplate := base64.StdEncoding.EncodeToString(_reqObj.ProbeCbor)

	secondReqObj := requestobjects.AuthReqObj{
		FlowId: flowID,
		SelectedAuthenticar: requestobjects.AuthObj_SelectedAuthenticator{
			AuthenticationID: authenticatorId,
			Params: requestobjects.AuthObj_Params{
				BiometricTemplate: BiometricTemplate,
			},
		}}
	jsonobj, err := json.Marshal(secondReqObj)
	if err != nil {
		err = fmt.Errorf("error occured while trying to convert reqObj to json, %w", err)
		return
	}
	secondRequestBody := bytes.NewBuffer(jsonobj)

	if err != nil {
		err = fmt.Errorf("error occured while trying to make second url string using url.JoinPath , %w", err)
		return
	}

	secondReq, err := http.NewRequest("POST", urlString, secondRequestBody)
	if err != nil {
		err = fmt.Errorf("error occured while trying to make second http post request : %w", err)
	}

	secondReq.Header.Add("Content-Type", "application/json")
	secondReq.Header.Add("Authorization", "Bearer "+accessToken)

	secondRes, errReqSend := internalClient.Do(secondReq)
	if errReqSend != nil {
		err = fmt.Errorf("error while sending or recieving post request : %w", errReqSend)
		return
	}
	defer secondRes.Body.Close()

	secondBodybytes, errReadAll := io.ReadAll(secondRes.Body)
	if errReadAll != nil {
		err = fmt.Errorf("error while reading bytes of response body : %w", errReadAll)
		return
	}
	if secondRes.StatusCode != 200 {
		err = fmt.Errorf("the idp responded with an error message for second request for authorization : %s", string(secondBodybytes))
		return
	}
	secondResult = gjson.ParseBytes(secondBodybytes)
	err = nil
	return
}
