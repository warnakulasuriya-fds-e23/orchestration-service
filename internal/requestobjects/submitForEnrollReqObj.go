package requestobjects

type SubmitForEnrollReqObj struct {
	Data     []byte `json:"data"`
	Id       string `json:"id"`
	DeviceId string `json:"deviceid"`
}
