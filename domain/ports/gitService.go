package ports

import "github.com/jeffvo/go-pr-reviewer/domain/entities"

type GitService interface {
	GetPullRequest(url string) ([]*entities.PullRequestChanges, error)
	PostCodeSuggestions(url string, suggestions []entities.Suggestion, metadata *entities.PullRequestMetadata) error
	GetPullRequestMetadata(url string) (*entities.PullRequestMetadata, error)
}
