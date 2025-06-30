package requestobjects

type SubmitForEnrollReqObj struct {
	NewEntryCbor []byte `json:"newentrycbor"`
	UserId       string `json:"userid"`
	ClientId     string `json:"clientid"`
}
