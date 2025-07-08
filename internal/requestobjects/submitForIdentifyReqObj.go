package requestobjects

type SubmitForIdentifyReqObj struct {
	ProbeCbor []byte `json:"probecbor"`
	DeviceId  string `json:"deviceid"`
}
