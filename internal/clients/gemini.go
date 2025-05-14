package clients

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/jeffvo/go-pr-reviewer/domain/entities"
	"google.golang.org/api/option"
)

type GeminiClientInterface interface {
	GetSuggestions(pullRequestFiles []*entities.PullRequestChanges) (*string, error)
}

type GeminiClient struct {
	token   string
	version string
	client  *genai.Client
}

func NewGeminiClient(token string, version string) *GeminiClient {
	ctx := context.Background()
	client, _ := genai.NewClient(ctx, option.WithAPIKey(token))
	defer client.Close()

	return &GeminiClient{token: token, version: version}
}

func (g *GeminiClient) GetSuggestions(pullRequestFiles []*entities.PullRequestChanges) (*string, error) {
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

	textString := convertGeminiResponseToString(resp)

	return &textString, nil
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
