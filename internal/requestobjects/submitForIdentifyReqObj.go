package requestobjects

type SubmitForIdentifyReqObj struct {
	ProbeCbor []byte `json:"probecbor"`
	ClientId  string `json:"clientid"`
}
