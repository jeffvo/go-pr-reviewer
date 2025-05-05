package entities

type Suggestion struct {
	StartLine             int    `json:"startLine"`
	EndLine               int    `json:"endLine"`
	Suggestion            string `json:"suggestion"`
	AdditionalInformation string `json:"additionalInformation"`
	FileName              string `json:"fileName"`
}
