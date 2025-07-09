package customstorage

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type tokenResponseObject struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	Scope            string `json:"scope,omitempty"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type TokenStorage struct {
	accessToken string
	expiryTime  time.Time
}

func NewTokenStorage() (tokenStorage *TokenStorage, err error) {
	tokenStorage = &TokenStorage{accessToken: "", expiryTime: time.Now()}
	_, err = tokenStorage.GetAccessToken()
	if err != nil {
		tokenStorage = nil
		err = fmt.Errorf("error while making a new token storage : %w", err)
		return
	}

	err = nil
	return
}

func (tokenStorage *TokenStorage) GetAccessToken() (token string, err error) {
	if tokenStorage.accessToken == "" || tokenStorage.expiryTime.Equal(time.Now()) || tokenStorage.expiryTime.Before(time.Now().Add(5*time.Second)) {

		tokenEndpoint := os.Getenv("TOKEN_ENDPOINT_FOR_OUTGOING")
		data := url.Values{}
		data.Set("grant_type", "client_credentials")
		requestBody := bytes.NewBufferString(data.Encode())
		req, errNewReq := http.NewRequest("POST", tokenEndpoint, requestBody)
		if errNewReq != nil {
			token = ""
			err = fmt.Errorf("error while creating a post request for the tokenEndpoint : %w", errNewReq)
			return
		}
		consumerKey := os.Getenv("CONSUMER_KEY_FOR_OUTGOING")
		consumerSecret := os.Getenv("CONSUMER_SECRET_FOR_OUTGOING")
		authHeadervalue := base64.StdEncoding.EncodeToString([]byte(consumerKey + ":" + consumerSecret))
		req.Header.Add("Authorization", "Basic "+authHeadervalue)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		internalclient := &http.Client{}
		res, errReqSend := internalclient.Do(req)
		if errReqSend != nil {
			token = ""
			err = fmt.Errorf("error while sending or recieving post request : %w", errReqSend)
			return
		}
		defer res.Body.Close()
		bodybytes, errReadAll := io.ReadAll(res.Body)
		if errReadAll != nil {
			token = ""
			err = fmt.Errorf("error while reading bytes of response body : %w", errReadAll)
			return
		}

		var resObj tokenResponseObject
		errUnMarshal := json.Unmarshal(bodybytes, &resObj)
		if errUnMarshal != nil {
			token = ""
			err = fmt.Errorf("error while running json unmarshal for the read bytes of the response body : %w", errUnMarshal)
			return
		}

		tokenStorage.expiryTime = time.Now().Add(time.Duration(resObj.ExpiresIn) * time.Second)
		tokenStorage.accessToken = resObj.AccessToken
		log.Println("obtained new access token from : ", tokenEndpoint)

	}
	token = tokenStorage.accessToken
	err = nil
	return
}
