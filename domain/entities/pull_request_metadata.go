package entities

type PullRequestMetadata struct {
	Head struct {
		Sha string `json:"sha"`
	} `json:"head"`
}
