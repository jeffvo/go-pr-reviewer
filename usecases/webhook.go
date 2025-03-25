package usecases

import (
	"github.com/jeffvo/go-pr-reviewer/internals"
	"github.com/jeffvo/go-pr-reviewer/internals/adapters"
)

type WebhookProcessor struct {
	githubAdapter *adapters.GithubAdapter
	geminiAdapter *adapters.GeminiAdapter
}

func NewWebhookProcessor(githubAdapter *adapters.GithubAdapter, geminiAdapter *adapters.GeminiAdapter) *WebhookProcessor {
	return &WebhookProcessor{githubAdapter: githubAdapter, geminiAdapter: geminiAdapter}
}

func (p *WebhookProcessor) ProcessWebhook(body []byte) error {
	webHookResult, err := internals.ParseWebhookPayload(body)
	if err != nil {
		return err
	}

	pullRequestFiles, err := p.githubAdapter.GetPullRequest(webHookResult.PullRequest.URL)
	if err != nil {
		return err
	}

	p.geminiAdapter.GetCodeSuggestions(pullRequestFiles)

	return nil
}
