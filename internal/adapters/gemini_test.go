package adapters_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeffvo/go-pr-reviewer/domain/entities"
	"github.com/jeffvo/go-pr-reviewer/internal/adapters"
	"github.com/stretchr/testify/assert"
)

func TestNewGeminiAdapter(t *testing.T) {
	t.Run("given_valid_api_key_when_creating_adapter_then_returns_non_nil_adapter", func(t *testing.T) {

		adapter := adapters.NewGeminiAdapter("test-api-key", "test-version")

		assert.NotNil(t, adapter)
	})
}

func TestReviewCode(t *testing.T) {
	t.Run("given_valid_pull_request_when_gemini_returns_suggestions_then_returns_suggestions", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"suggestions": [{"file": "test.go", "suggestion": "Use test2 instead of test"}]}`))
		}))
		defer server.Close()
		adapter := adapters.NewGeminiAdapter("test-api-key", "test-version")
		pullRequestChanges := []*entities.PullRequestChanges{
			{
				FileName: "test.go",
				Changed:  "@@ -1,1 +1,1 @@\n-test\n+test2",
			},
		}
		suggestions, err := adapter.GetCodeSuggestions(pullRequestChanges)
		assert.NoError(t, err)
		assert.NotNil(t, suggestions)
		assert.Len(t, suggestions, 1)
		assert.Equal(t, "test.go", suggestions[0].FileName)
		assert.Equal(t, "Use test2 instead of test", suggestions[0].Suggestion)
	})

	t.Run("given_valid_pull_request_when_gemini_returns_error_then_returns_error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": {"message": "Invalid request", "code": 400}}`))
		}))
		defer server.Close()

		adapter := adapters.NewGeminiAdapter("test-api-key", "test-version")
		pullRequestChanges := []*entities.PullRequestChanges{
			{
				FileName: "test.go",
				Changed:  "@@ -1,1 +1,1 @@\n-test\n+test2",
			},
		}

		suggestions, err := adapter.GetCodeSuggestions(pullRequestChanges)

		assert.Error(t, err)
		assert.Nil(t, suggestions)
	})

	t.Run("given_valid_pull_request_when_server_connection_fails_then_returns_error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		}))
		server.Close()

		adapter := adapters.NewGeminiAdapter("test-api-key", "test-version")

		pullRequestChanges := []*entities.PullRequestChanges{
			{
				FileName: "test.go",
				Changed:  "@@ -1,1 +1,1 @@\n-test\n+test2",
			},
		}

		suggestions, err := adapter.GetCodeSuggestions(pullRequestChanges)

		assert.Error(t, err)
		assert.Nil(t, suggestions)
	})

	t.Run("given_valid_pull_request_when_gemini_returns_invalid_json_then_returns_error", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{invalid json`))
		}))
		defer server.Close()

		adapter := adapters.NewGeminiAdapter("test-api-key", "test-version")

		pullRequestChanges := []*entities.PullRequestChanges{
			{
				FileName: "test.go",
				Changed:  "@@ -1,1 +1,1 @@\n-test\n+test2",
			},
		}

		suggestions, err := adapter.GetCodeSuggestions(pullRequestChanges)

		assert.Error(t, err)
		assert.Nil(t, suggestions)
	})

	t.Run("given_valid_pull_request_when_gemini_returns_malformed_response_then_returns_error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"something": "unexpected"}`))
		}))
		defer server.Close()

		adapter := adapters.NewGeminiAdapter("test-api-key", "test-version")

		pullRequestChanges := []*entities.PullRequestChanges{
			{
				FileName: "test.go",
				Changed:  "@@ -1,1 +1,1 @@\n-test\n+test2",
			},
		}

		suggestions, err := adapter.GetCodeSuggestions(pullRequestChanges)

		assert.Error(t, err)
		assert.Nil(t, suggestions)
	})
}
