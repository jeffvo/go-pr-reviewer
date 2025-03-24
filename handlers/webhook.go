package handlers

import (
	"io"
	"net/http"

	"github.com/jeffvo/go-pr-reviewer/usecases"
)

type WebhookHandler struct {
	usecase *usecases.WebhookProcessor
}

func NewWebhookHandler(usecase *usecases.WebhookProcessor) *WebhookHandler {
	return &WebhookHandler{usecase}
}

func (wbh *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Check if the request is a POST request
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = wbh.usecase.ProcessWebhook(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
