package usecases

import (
	"fmt"

	"github.com/jeffvo/go-pr-reviewer/domain/ports"
	"github.com/jeffvo/go-pr-reviewer/internal/usecases/dto"
)

type WebhookProcessor struct {
	gitService ports.GitService
	aiService  ports.AIService
}

func NewWebhookProcessor(githubAdapter ports.GitService, aiService ports.AIService) *WebhookProcessor {
	return &WebhookProcessor{gitService: githubAdapter, aiService: aiService}
}

func (p *WebhookProcessor) ProcessWebhook(payload dto.WebhookPayload) error {

	pullRequestFiles, err := p.gitService.GetPullRequest(payload.PullRequest.URL)
	if err != nil {
		fmt.Printf("Error getting pull request: %v\n", err)
		return err
	}

	metadata, err := p.gitService.GetPullRequestMetadata(payload.PullRequest.URL)
	if err != nil {
		fmt.Printf("Error getting pull request metadata: %v\n", err)
		return err
	}

	suggestions, err := p.aiService.GetCodeSuggestions(pullRequestFiles)

	if err != nil {
		fmt.Printf("Error getting code suggestions: %v\n", err)
		return err
	}

	err = p.gitService.PostCodeSuggestions(payload.PullRequest.URL, suggestions, metadata)
	if err != nil {
		fmt.Printf("Error posting code suggestions: %v\n", err)
		return err
	}

	return nil
}
