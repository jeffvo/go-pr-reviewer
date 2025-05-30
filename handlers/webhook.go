package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jeffvo/go-pr-reviewer/internal/usecases"
	"github.com/jeffvo/go-pr-reviewer/internal/usecases/dto"
)

type WebhookHandler struct {
	usecase *usecases.WebhookProcessor
}

func NewWebhookHandler(usecase *usecases.WebhookProcessor) *WebhookHandler {
	return &WebhookHandler{usecase}
}

func (wbh *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var webhook dto.WebhookPayload
	err = json.Unmarshal([]byte(body), &webhook)
	if err != nil {
		fmt.Printf("Failed to unmarshal webhook: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	url, err := webhook.GetPullRequestURL()
	if err != nil {
		fmt.Printf("Failed to get pull request URL: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = wbh.usecase.ProcessWebhook(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
