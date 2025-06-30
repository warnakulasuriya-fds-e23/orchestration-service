package requestobjects

type SubmitForMatchReqObj struct {
	ProbeCbor     []byte `json:"probecbor"`
	CandidateCbor []byte `json:"candidatecbor"`
	ClientId      string `json:"clientid"`
}
