package authorizationcalls

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/requestobjects"
)

func CredentialSubmissionCall(urlString string, internalClient *http.Client, initialResult *gjson.Result, _reqObj requestobjects.SubmitForIdentifyReqObj) (secondResult gjson.Result, err error) {
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
