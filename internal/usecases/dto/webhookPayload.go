package dto

type WebhookPayload struct {
	PullRequest struct {
		URL string `json:"url"`
	} `json:"pull_request"`
}
