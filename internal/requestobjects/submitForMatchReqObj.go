package requestobjects

type SubmitForMatchReqObj struct {
	ProbeCbor     []byte `json:"probecbor"`
	CandidateCbor []byte `json:"candidatecbor"`
	DeviceId      string `json:"deviceid"`
}
