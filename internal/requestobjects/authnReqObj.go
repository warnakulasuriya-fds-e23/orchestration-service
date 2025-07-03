package requestobjects

type AuthReqObj struct {
	FlowId              string                `json:"flowId"`
	SelectedAuthenticar SelectedAuthenticator `json:"selectedAuthenticator"`
}

type SelectedAuthenticator struct {
	AuthenticationID string `json:"authenticatorId"`
	Params           Params `json:"params"`
}

type Params struct {
	BiometricKey string `json:"biometric-key"`
}
