package entities

import "fmt"

type CommentPayload struct {
	Body      string `json:"body"`
	CommitID  string `json:"commit_id"`
	Path      string `json:"path"`
	StartLine int    `json:"start_line,omitempty"`
	EndLine   int    `json:"line"`
}

func (s *Suggestion) ToCommentPayload(commitId string) CommentPayload {
	body := ""
	if s.StartLine == s.EndLine {
		body = fmt.Sprintf("%s\n\n%s", s.AdditionalInformation, s.Suggestion)
	} else {
		body = fmt.Sprintf("```suggestion\n%s\n```\n%s", s.Suggestion, s.AdditionalInformation)
	}

	payload := CommentPayload{
		Body:     body,
		CommitID: commitId,
		Path:     s.FileName,
		EndLine:  s.EndLine,
	}

	// This is needed for a single line suggestion
	if s.StartLine != s.EndLine {
		payload.StartLine = s.StartLine
	}

	return payload
}
