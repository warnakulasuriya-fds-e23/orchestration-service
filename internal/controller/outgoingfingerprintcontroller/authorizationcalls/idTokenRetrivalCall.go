package authorizationcalls

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/tidwall/gjson"
)

func IdTokenRetrivalCall(accessToken string, urlString string, internalClient *http.Client, secondResult gjson.Result) (idToken string, err error) {
	code := secondResult.Get("authData.code").String()
	trData := url.Values{}
	trData.Set("grant_type", "authorization_code")
	trData.Set("code", code)
	trData.Set("redirect_uri", os.Getenv("IDP_APP_REDIRECT_URI"))

	if err != nil {
		err = fmt.Errorf("error occured while trying to make token endpoint url string using url.JoinPath , %w", err)
		return
	}
	requestBody := bytes.NewBufferString(trData.Encode())
	tokenReq, err := http.NewRequest("POST", urlString, requestBody)
	if err != nil {
		err = fmt.Errorf("error occured while trying to make a new http post request : %w", err)
		return
	}
	consumerKey := os.Getenv("CONSUMER_KEY_FOR_OUTGOING")
	consumerSecret := os.Getenv("CONSUMER_SECRET_FOR_OUTGOING")
	authHeadervalue := base64.StdEncoding.EncodeToString([]byte(consumerKey + ":" + consumerSecret))
	tokenReq.Header.Add("Authorization", "Basic "+authHeadervalue)
	tokenReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	tokenReq.Header.Add("Authorization", "Bearer "+accessToken)
	tokenRes, errReqSend := internalClient.Do(tokenReq)
	if errReqSend != nil {
		err = fmt.Errorf("error while sending or recieving post request : %w", errReqSend)
		return
	}
	defer tokenRes.Body.Close()

	tokenResBodyBytes, errReadAll := io.ReadAll(tokenRes.Body)
	if errReadAll != nil {
		err = fmt.Errorf("error while reading bytes of response body : %w", errReadAll)
		return
	}
	if tokenRes.StatusCode != 200 {
		err = fmt.Errorf("the idp responded with an error message for second request for authorization : %s", string(tokenResBodyBytes))
		return
	}
	tokenReqResult := gjson.ParseBytes(tokenResBodyBytes)
	idToken = tokenReqResult.Get("id_token").String()
	err = nil
	return
}
