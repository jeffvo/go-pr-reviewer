package adapters

import (
	"encoding/json"
	"fmt"

	"github.com/jeffvo/go-pr-reviewer/domain/entities"
	"github.com/jeffvo/go-pr-reviewer/internal/clients"
)

type GeminiAdapter struct {
	token        string
	version      string
	geminiClient clients.GeminiClientInterface
}

func NewGeminiAdapter(token string, version string) *GeminiAdapter {
	client := clients.NewGeminiClient(token, version)

	return &GeminiAdapter{token: token, version: version, geminiClient: client}
}

func NewGeminiAdapterWithClient(token string, version string, client clients.GeminiClientInterface) *GeminiAdapter {
	return &GeminiAdapter{token: token, version: version, geminiClient: client}
}

type Changes struct {
	FileName string
	Changes  []string
}

func (g *GeminiAdapter) GetCodeSuggestions(pullRequestFiles []*entities.PullRequestChanges) ([]entities.Suggestion, error) {

	text, err := g.geminiClient.GetSuggestions(pullRequestFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggestions from Gemini: %w", err)
	}

	parsedResponse, err := parseJSONResponse(*text)

	if err != nil {
		return nil, err
	}

	return parsedResponse, nil
}

func parseJSONResponse(jsonString string) ([]entities.Suggestion, error) {
	var responses []entities.Suggestion
	err := json.Unmarshal([]byte(jsonString), &responses)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return responses, nil
}
