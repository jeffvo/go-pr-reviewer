package usecases

import (
	"fmt"

	"github.com/jeffvo/go-pr-reviewer/adapters"
	"github.com/jeffvo/go-pr-reviewer/internals"
)

type WebhookProcessor struct {
	githubAdapter *adapters.GithubAdapter
}

func NewWebhookProcessor(githubAdapter *adapters.GithubAdapter) *WebhookProcessor {
	return &WebhookProcessor{githubAdapter: githubAdapter}
}

func (p *WebhookProcessor) ProcessWebhook(body []byte) error {
	webHookResult, err := internals.ParseWebhookPayload(body)
	if err != nil {
		return err
	}

	fmt.Printf("Received a webhook: %s", webHookResult)
	p.githubAdapter.GetPullRequest(webHookResult.PullRequest.URL)

	return nil
}
