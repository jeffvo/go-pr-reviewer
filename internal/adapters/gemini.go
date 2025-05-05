package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"github.com/jeffvo/go-pr-reviewer/domain/entities"
	"google.golang.org/api/option"
)

type GeminiAdapter struct {
	token string
}

func NewGeminiAdapter(token string) *GeminiAdapter {
	return &GeminiAdapter{token: token}
}

type Changes struct {
	FileName string
	Changes  []string
}

func (g *GeminiAdapter) GetCodeSuggestions(pullRequestFiles []*entities.PullRequestChanges) (*[]entities.Suggestion, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(g.token))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	var text string
	for _, file := range pullRequestFiles {
		text += fmt.Sprintf("%s\n%s", file.FileName, file.Changed)
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx,

		genai.Text("Could you review the following code changes? the file name will be shown first then the changes. Could you return it in json format which is the following { startLine: int, endLine:int, suggestion: string, additionalInformation: string, fileName: string}? Please only fill the suggestion field with actual code suggestions. Also, use the additionalInformation field to explain why this code suggestion is recommended. Thank you. \n"+text))

	if err != nil {
		return nil, err
	}

	// For now I directly unmarshal the result to the PullRequestInfo struct
	// This is done as we do not want to manipulate the data so having an extra layer feels unnecessary
	textString := convertGeminiResponeToString(resp)

	parsedResponse, err := parseJSONResponse(textString)

	if err != nil {
		return nil, err
	}

	return &parsedResponse, nil
}

func convertGeminiResponeToString(resp *genai.GenerateContentResponse) string {
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
