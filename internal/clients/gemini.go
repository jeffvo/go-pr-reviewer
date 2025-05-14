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
	GetSuggestions(pullRequestFiles []*entities.PullRequestChanges) (string, error)
}

type GeminiClient struct {
	token   string
	model   *genai.GenerativeModel
	context context.Context
}

func NewGeminiClient(token string, version string) *GeminiClient {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(token))

	if err != nil {
		fmt.Printf("Error creating Gemini client: %v\n", err)
		return nil
	}

	model := client.GenerativeModel(version)

	return &GeminiClient{token: token, model: model, context: ctx}
}

func (g *GeminiClient) GetSuggestions(pullRequestFiles []*entities.PullRequestChanges) (string, error) {
	var text strings.Builder
	for _, file := range pullRequestFiles {
		text.WriteString(fmt.Sprintf("%s\n%s", file.FileName, file.Changed))
	}

	genText := genai.Text(fmt.Sprintf("Could you review the following code changes? the file name will be shown first then the changes. Could you return it in json format which is the following { startLine: int, endLine:int, suggestion: string, additionalInformation: string, fileName: string}? Please only fill the suggestion field with actual code suggestions. Also, use the additionalInformation field to explain why this code suggestion is recommended. Thank you. \n%s", text.String()))

	ctx := context.Background()
	resp, err := g.model.GenerateContent(ctx, genText)

	if err != nil {
		fmt.Printf("Error generating content: %v\n", err)
		return "", err
	}

	textString := convertGeminiResponseToString(resp)

	return textString, nil
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
