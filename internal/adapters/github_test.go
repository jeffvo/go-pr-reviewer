package adapters_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeffvo/go-pr-reviewer/domain/entities"
	"github.com/jeffvo/go-pr-reviewer/internal/adapters"
	"github.com/stretchr/testify/assert"
)

func TestGetPullRequest(t *testing.T) {
	t.Run("given_valid_repository_when_getting_pull_request_then_returns_success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/files", r.URL.Path)
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "application/vnd.github+json", r.Header.Get("Accept"))
			assert.Contains(t, r.Header.Get("Authorization"), "Bearer ")
			assert.Equal(t, "2022-11-28", r.Header.Get("X-Github-Api-Version"))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"filename": "test.go", "patch": "@@ -1 +1 @@\n+Test" }]`))
		}))
		defer server.Close()

		adapter := adapters.NewGithubAdapter("test-token")
		repository := server.URL

		pullRequestFiles, err := adapter.GetPullRequest(repository)

		assert.NoError(t, err)
		assert.Len(t, pullRequestFiles, 1)
		assert.Equal(t, "test.go", pullRequestFiles[0].FileName)
	})

	t.Run("given_invalid_repository_when_getting_pull_request_then_returns_error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "Bad request"}`))
		}))
		defer server.Close()

		adapter := adapters.NewGithubAdapter("test-token")
		repository := server.URL

		pullRequestFiles, err := adapter.GetPullRequest(repository)

		assert.Error(t, err)
		assert.Nil(t, pullRequestFiles)
	})
}

func TestGetPullRequestMetadata(t *testing.T) {
	t.Run("given_valid_repository_when_getting_pull_request_metadata_then_returns_success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/", r.URL.Path)
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "application/vnd.github+json", r.Header.Get("Accept"))
			assert.Contains(t, r.Header.Get("Authorization"), "Bearer ")
			assert.Equal(t, "2022-11-28", r.Header.Get("X-Github-Api-Version"))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"head": {"sha": "abc123"}}`))
		}))
		defer server.Close()

		adapter := adapters.NewGithubAdapter("test-token")
		repository := server.URL

		pullRequestMetadata, err := adapter.GetPullRequestMetadata(repository)

		assert.NoError(t, err)
		assert.Equal(t, "abc123", pullRequestMetadata.Head.Sha)
	})

	t.Run("given_invalid_repository_when_getting_pull_request_metadata_then_returns_error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "Bad request"}`))
		}))
		defer server.Close()

		adapter := adapters.NewGithubAdapter("test-token")
		repository := server.URL

		pullRequestMetadata, err := adapter.GetPullRequestMetadata(repository)

		assert.Error(t, err)
		assert.Nil(t, pullRequestMetadata)
	})
}

func TestPostCodeSuggestions(t *testing.T) {

	metadata := &entities.PullRequestMetadata{
		Head: struct {
			Sha string "json:\"sha\""
		}{
			Sha: "abc123",
		},
	}

	t.Run("given_valid_suggestions_when_posting_to_github_then_returns_success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/comments", r.URL.Path)
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Accept"))
			assert.Contains(t, r.Header.Get("Authorization"), "Bearer ")
			assert.Equal(t, "2022-11-28", r.Header.Get("X-Github-Api-Version"))

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id": 123}`))
		}))
		defer server.Close()

		adapter := adapters.NewGithubAdapter("test-token")
		suggestions := []entities.Suggestion{
			{
				StartLine:             10,
				EndLine:               15,
				Suggestion:            "Test suggestion",
				AdditionalInformation: "Test additional information",
			},
		}

		err := adapter.PostCodeSuggestions(server.URL, suggestions, metadata)

		assert.NoError(t, err)
	})

	t.Run("given_valid_suggestions_when_github_returns_non_201_status_then_returns_error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "Bad request"}`))
		}))
		defer server.Close()

		adapter := adapters.NewGithubAdapter("test-token")
		suggestions := []entities.Suggestion{
			{
				StartLine:             10,
				EndLine:               15,
				Suggestion:            "Test suggestion",
				AdditionalInformation: "Test additional information",
			},
		}

		err := adapter.PostCodeSuggestions(server.URL, suggestions, metadata)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to post code suggestions")
	})

	t.Run("given_valid_suggestions_when_server_connection_fails_then_returns_error", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		}))
		server.Close()

		adapter := adapters.NewGithubAdapter("test-token")
		suggestions := []entities.Suggestion{
			{
				StartLine:             10,
				EndLine:               15,
				Suggestion:            "Test suggestion",
				AdditionalInformation: "Test additional information",
			},
		}

		err := adapter.PostCodeSuggestions(server.URL, suggestions, metadata)

		assert.Error(t, err)
	})

	t.Run("given_multiple_suggestions_when_posting_to_github_then_makes_separate_requests_for_each", func(t *testing.T) {

		requestCount := 0

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id": 123}`))
		}))
		defer server.Close()

		adapter := adapters.NewGithubAdapter("test-token")
		suggestions := []entities.Suggestion{
			{
				StartLine:             10,
				EndLine:               15,
				Suggestion:            "Test suggestion 1",
				AdditionalInformation: "Test additional information 1",
			},
			{
				StartLine:             20,
				EndLine:               25,
				Suggestion:            "Test suggestion 2",
				AdditionalInformation: "Test additional information 2",
			},
		}

		err := adapter.PostCodeSuggestions(server.URL, suggestions, metadata)

		assert.NoError(t, err)
		assert.Equal(t, 2, requestCount, "Should have made exactly 2 requests")
	})
}
