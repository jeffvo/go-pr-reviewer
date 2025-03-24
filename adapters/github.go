package adapters

import (
	"fmt"
	"net/http"
)

type GithubAdapter struct {
	token string
}

func NewGithubAdapter(token string) *GithubAdapter {
	return &GithubAdapter{token: token}
}

func (g *GithubAdapter) GetPullRequest(repository string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/files", repository), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Add("Accept", `application/vnd.github+json`)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", g.token))
	req.Header.Add("X-Github-Api-Version", "2022-11-28")
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Sprintf("Received a webhook: %s", resp)

	return err

}
