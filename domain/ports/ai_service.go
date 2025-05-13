package ports

import "github.com/jeffvo/go-pr-reviewer/domain/entities"

type AIService interface {
	GetCodeSuggestions(pullRequestFiles []*entities.PullRequestChanges) ([]entities.Suggestion, error)
}
