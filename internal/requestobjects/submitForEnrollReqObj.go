package requestobjects

type SubmitForEnrollReqObj struct {
	Data     []byte `json:"data"`
	Id       string `json:"id"`
	ClientId string `json:"clientid"`
}
