package outgoingfingerprintcontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/requestobjects"
	"github.com/warnakulasuriya-fds-e23/orchestration-service/internal/responseobjects"
)

func (controller *OutgoingFingerprintController) outgoingEnrollHandler(_reqObj requestobjects.SubmitForEnrollReqObj) (message string, err error) {
	accessToken, err := controller.tokenStorage.GetAccessToken()
	if err != nil {
		message = ""
		err = fmt.Errorf("error occured while trying to get access token: %w", err)
		return
	}
	newEntryData := _reqObj.Data
	userId := _reqObj.Id
	reqObj := requestobjects.EnrollTemplateReqObj{Data: newEntryData, Id: userId}
	jsonobj, err := json.Marshal(reqObj)
	if err != nil {
		message = ""
		err = fmt.Errorf("error occured while trying to convert reqObj to json, %w", err)
		return
	}

	urlString, err := url.JoinPath(controller.targetAdress, EnrollTemplateEndpoint)
	if err != nil {
		message = ""
		err = fmt.Errorf("error occured while trying to url string using url.JoinPath , %w", err)
		return
	}
	requestBody := bytes.NewBuffer(jsonobj)
	req, err := http.NewRequest("POST", urlString, requestBody)
	if err != nil {
		message = ""
		// TODO: properly handle this error
		log.Fatal(err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	internalClient := &http.Client{}
	resp, err := internalClient.Do(req)
	if err != nil {
		message = ""
		err = fmt.Errorf("error occured while using http.Post , %w", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		message = ""
		err = fmt.Errorf("error occured while reading response bytes using io.ReadAll , %w", err)
		return
	}

	if resp.StatusCode != 200 {
		var resObj responseobjects.ErrorResObj
		err = json.Unmarshal(bodyBytes, &resObj)
		if err != nil {
			err = fmt.Errorf("error occured while runnig json.Unmarshal on response bytes , %w", err)
			return
		}
		message = ""
		err = fmt.Errorf("error occured in bio-sdk-service , %s", resObj.Message)
		return
	}
	var resObj responseobjects.EnrollTemplateResObj
	err = json.Unmarshal(bodyBytes, &resObj)
	if err != nil {
		err = fmt.Errorf("error occured while runnig json.Unmarshal on response bytes , %w", err)
		return
	}
	message = resObj.Message
	err = nil
	return
}
