package dto

import "encoding/json"

type WebhookPayload struct {
	Payload string `json:"payload"`
}

type PullRequestPayload struct {
	PullRequest struct {
		URL string `json:"url"`
	} `json:"pull_request"`
}

func (w *WebhookPayload) GetPullRequestURL() (string, error) {
	var prPayload PullRequestPayload
	err := json.Unmarshal([]byte(w.Payload), &prPayload)
	if err != nil {
		return "", err
	}
	return prPayload.PullRequest.URL, nil
}
