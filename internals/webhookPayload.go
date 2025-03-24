package internals

import (
	"encoding/json"
	"fmt"
)

// WebhookPayload represents the structure of the webhook payload.
type WebhookPayload struct {
	PullRequest struct {
		URL string `json:"url"`
	} `json:"pull_request"`
}

// ParseWebhookPayload parses the raw JSON payload into a WebhookPayload struct.
func ParseWebhookPayload(rawPayload []byte) (*WebhookPayload, error) {
	var payload WebhookPayload
	if err := json.Unmarshal(rawPayload, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}
	return &payload, nil
}
