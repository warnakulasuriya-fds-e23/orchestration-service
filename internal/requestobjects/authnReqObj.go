package requestobjects

type AuthReqObj struct {
	FlowId              string                        `json:"flowId"`
	SelectedAuthenticar AuthObj_SelectedAuthenticator `json:"selectedAuthenticator"`
}

type AuthObj_SelectedAuthenticator struct {
	AuthenticationID string         `json:"authenticatorId"`
	Params           AuthObj_Params `json:"params"`
}

type AuthObj_Params struct {
	BiometricKey string `json:"biometric-key"`
}
