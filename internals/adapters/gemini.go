package adapters

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiAdapter struct {
	token string
}

func NewGeminiAdapter(token string) *GeminiAdapter {
	return &GeminiAdapter{token: token}
}

type ChangedFiles struct {
	FileName string
	Changes  []string
}

func (g *GeminiAdapter) GetCodeSuggestions(pullRequestFiles []*PullRequestFiles) error {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(g.token))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	var text string
	for _, file := range pullRequestFiles {
		text += fmt.Sprintf("%s\n%s", file.Filename, file.Patch)
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx,

		genai.Text("Could you review the following code changes? the file name will be shown first then the changes please respond in a way I can paste it as suggestions to the PR. Could you return it in json format which is the following { startLine: int, endline:int, suggestions: string, additionalInformation: string}. In the suggestion field could you post the code suggestion that will be interperted correctly by github? \n"+text))

	if err != nil {
		return err
	}

	fmt.Printf("Suggestions: %v\n", resp.Candidates[0].Content.Parts)

	return nil
}
