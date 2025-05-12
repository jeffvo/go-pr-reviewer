package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jeffvo/go-pr-reviewer/domain/entities"
)

type GithubAdapter struct {
	token string
}

func NewGithubAdapter(token string) *GithubAdapter {
	return &GithubAdapter{token: token}
}

type PullRequestFiles struct {
	Filename string `json:"filename"`
	Patch    string `json:"patch"`
}

func (g *GithubAdapter) GetPullRequest(repository string) ([]*entities.PullRequestChanges, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/files", repository), nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Add("Accept", `application/vnd.github+json`)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", g.token))
	req.Header.Add("X-Github-Api-Version", "2022-11-28")
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// For now I directly unmarshal the result to the PullRequestInfo struct
	// This is done as we do not want to manipulate the data so having an extra layer feels unnecessary
	var pullRequestInfo []*entities.PullRequestChanges
	if err := json.Unmarshal(result, &pullRequestInfo); err != nil {
		return nil, fmt.Errorf("failed to parse the API request result: %w", err)
	}

	return pullRequestInfo, nil
}

func (g *GithubAdapter) GetPullRequestMetadata(repository string) (*entities.PullRequestMetadata, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s", repository), nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Add("Accept", `application/vnd.github+json`)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", g.token))
	req.Header.Add("X-Github-Api-Version", "2022-11-28")
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// For now I directly unmarshal the result to the PullRequestInfo struct
	// This is done as we do not want to manipulate the data so having an extra layer feels unnecessary
	var metadata *entities.PullRequestMetadata
	if err := json.Unmarshal(result, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse the API request result: %w", err)
	}

	return metadata, nil
}

func (g *GithubAdapter) PostCodeSuggestions(url string, suggestions []entities.Suggestion, metadata *entities.PullRequestMetadata) error {
	client := &http.Client{}

	for _, suggestion := range suggestions {
		body, err := json.Marshal(suggestion.ToCommentPayload(metadata.Head.Sha))
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/comments", url), bytes.NewReader(body))
		if err != nil {
			return err
		}

		fmt.Println("", string(body))

		req.Header.Add("Accept", `application/json`)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", g.token))
		req.Header.Add("X-Github-Api-Version", "2022-11-28")
		resp, err := client.Do(req)

		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusCreated {
			resp.Body.Close()
			return fmt.Errorf("failed to post code suggestions: %s", resp.Status)
		}

		resp.Body.Close()
	}

	return nil
}
