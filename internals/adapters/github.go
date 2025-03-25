package adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (g *GithubAdapter) GetPullRequest(repository string) ([]*PullRequestFiles, error) {
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

	var pullRequestInfo []*PullRequestFiles
	if err := json.Unmarshal(result, &pullRequestInfo); err != nil {
		return nil, fmt.Errorf("failed to parse the API request result: %w", err)
	}

	return pullRequestInfo, nil

}
