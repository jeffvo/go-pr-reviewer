package entities

type PullRequestChanges struct {
	FileName string `json:"filename"`
	Changed  string `json:"patch"`
}
