package usecases

import (
	"fmt"

	"github.com/jeffvo/go-pr-reviewer/domain/ports"
)

type WebhookProcessor struct {
	gitService ports.GitService
	aiService  ports.AIService
}

func NewWebhookProcessor(githubAdapter ports.GitService, aiService ports.AIService) *WebhookProcessor {
	return &WebhookProcessor{gitService: githubAdapter, aiService: aiService}
}

func (p *WebhookProcessor) ProcessWebhook(url string) error {

	pullRequestFiles, err := p.gitService.GetPullRequest(url)
	if err != nil {
		fmt.Printf("Error getting pull request: %v\n", err)
		return err
	}

	metadata, err := p.gitService.GetPullRequestMetadata(url)
	if err != nil {
		fmt.Printf("Error getting pull request metadata: %v\n", err)
		return err
	}

	suggestions, err := p.aiService.GetCodeSuggestions(pullRequestFiles)

	if err != nil {
		fmt.Printf("Error getting code suggestions: %v\n", err)
		return err
	}

	err = p.gitService.PostCodeSuggestions(url, suggestions, metadata)
	if err != nil {
		fmt.Printf("Error posting code suggestions: %v\n", err)
		return err
	}

	return nil
}
