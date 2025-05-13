package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/jeffvo/go-pr-reviewer/domain/entities"
	"google.golang.org/api/option"
)

type GeminiAdapter struct {
	token   string
	version string
	client  *genai.Client
}

func NewGeminiAdapter(token string, version string) *GeminiAdapter {
	ctx := context.Background()
	client, _ := genai.NewClient(ctx, option.WithAPIKey(token))
	defer client.Close()
	return &GeminiAdapter{token: token, version: version, client: client}
}

// For testing purposes
func NewGeminiAdapterWithClient(token string, version string, client *genai.Client) *GeminiAdapter {
	return &GeminiAdapter{token: token, version: version, client: client}
}

type Changes struct {
	FileName string
	Changes  []string
}

func (g *GeminiAdapter) GetCodeSuggestions(pullRequestFiles []*entities.PullRequestChanges) ([]entities.Suggestion, error) {
	ctx := context.Background()

	var text strings.Builder
	for _, file := range pullRequestFiles {
		fmt.Fprintf(&text, "%s\n%s", file.FileName, file.Changed)
	}

	model := g.client.GenerativeModel(g.version)
	resp, err := model.GenerateContent(ctx,

		genai.Text("Could you review the following code changes? the file name will be shown first then the changes. Could you return it in json format which is the following { startLine: int, endLine:int, suggestion: string, additionalInformation: string, fileName: string}? Please only fill the suggestion field with actual code suggestions. Also, use the additionalInformation field to explain why this code suggestion is recommended. Thank you. \n"+text.String()))

	if err != nil {
		return nil, err
	}

	// For now I directly unmarshal the result to the PullRequestInfo struct
	// This is done as we do not want to manipulate the data so having an extra layer feels unnecessary
	textString := convertGeminiResponseToString(resp)

	parsedResponse, err := parseJSONResponse(textString)

	if err != nil {
		return nil, err
	}

	return parsedResponse, nil
}

func convertGeminiResponseToString(resp *genai.GenerateContentResponse) string {
	stringValue := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if textPart, ok := part.(genai.Text); ok {
			stringValue = string(textPart)
		}
	}

	//Removing the first 8 and last 4 characters as gemini adds some extra characters
	return stringValue[8 : len(stringValue)-4]
}

func parseJSONResponse(jsonString string) ([]entities.Suggestion, error) {
	var responses []entities.Suggestion
	err := json.Unmarshal([]byte(jsonString), &responses)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return responses, nil
}
