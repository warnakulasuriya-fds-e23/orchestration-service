package responseobjects

type IdentifyTemplateResObj struct {
	IsMatched    bool   `json:"ismatched"`
	DiscoveredId string `json:"discoveredid"`
}
