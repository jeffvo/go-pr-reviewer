package handlers

import (
	"encoding/json"
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

	var payload dto.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = wbh.usecase.ProcessWebhook(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
