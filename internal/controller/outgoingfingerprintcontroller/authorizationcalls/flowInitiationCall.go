package authorizationcalls

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/tidwall/gjson"
)

func FlowInitiationCall(urlString string, internalClient *http.Client) (initialResult gjson.Result, err error) {
	data := url.Values{}
	data.Set("client_id", os.Getenv("CONSUMER_KEY_FOR_OUTGOING"))
	data.Set("response_type", "code")
	data.Set("redirect_uri", os.Getenv("IDP_APP_REDIRECT_URI"))
	data.Set("state", "Authentication via fingerprint")
	data.Set("scope", "openid internal_login FingerprintAuth roles")
	data.Set("response_mode", "direct")

	requestBody := bytes.NewBufferString(data.Encode())

	initialreq, errNewRequest := http.NewRequest("POST", urlString, requestBody)
	if errNewRequest != nil {
		err = fmt.Errorf("error occured while trying to make a new http post request : %w", errNewRequest)
	}
	initialreq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	initialreq.Header.Add("Accept", "application/json")

	// tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	// internalclient := &http.Client{Transport: tr}
	res, errReqSend := internalClient.Do(initialreq)
	if errReqSend != nil {
		err = fmt.Errorf("error while sending or recieving post request : %w", errReqSend)
		return
	}
	defer res.Body.Close()
	bodybytes, errReadAll := io.ReadAll(res.Body)
	if errReadAll != nil {
		err = fmt.Errorf("error while reading bytes of response body : %w", errReadAll)
		return
	}
	if res.StatusCode != 200 {
		err = fmt.Errorf("the idp responded with an error message for first request for authorization : %s", string(bodybytes))
		return
	}
	err = nil
	initialResult = gjson.ParseBytes(bodybytes)
	return

}
