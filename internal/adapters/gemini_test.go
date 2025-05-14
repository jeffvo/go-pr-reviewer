package adapters_test

import (
	"errors"
	"testing"

	"github.com/jeffvo/go-pr-reviewer/domain/entities"
	"github.com/jeffvo/go-pr-reviewer/internal/adapters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGeminiClient struct {
	mock.Mock
}

func (m *MockGeminiClient) GetSuggestions(pullRequestFiles []*entities.PullRequestChanges) (*string, error) {
	args := m.Called(pullRequestFiles)
	if args.Get(0) != nil {
		result := args.Get(0).(string)
		return &result, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGeminiAdapterWithMockClient(t *testing.T) {
	t.Run("given_valid_pull_request_when_gemini_returns_suggestions_then_returns_suggestions", func(t *testing.T) {
		mockClient := new(MockGeminiClient)

		adapter := adapters.NewGeminiAdapterWithClient("test-api-key", "test-version", mockClient)

		pullRequestChanges := []*entities.PullRequestChanges{
			{
				FileName: "test.go",
				Changed:  "@@ -1,1 +1,1 @@\n-test\n+test2",
			},
		}

		mockResponse := `[{"FileName": "test.go", "Suggestion": "Use test2 instead of test"}]`
		mockClient.On("GetSuggestions", pullRequestChanges).Return(mockResponse, nil)

		suggestions, err := adapter.GetCodeSuggestions(pullRequestChanges)

		assert.NoError(t, err)
		assert.NotNil(t, suggestions)
		assert.Len(t, suggestions, 1)
		assert.Equal(t, "test.go", suggestions[0].FileName)
		assert.Equal(t, "Use test2 instead of test", suggestions[0].Suggestion)

		mockClient.AssertExpectations(t)
	})

	t.Run("given_valid_pull_request_when_gemini_returns_error_then_returns_error", func(t *testing.T) {
		mockClient := new(MockGeminiClient)

		adapter := adapters.NewGeminiAdapterWithClient("test-api-key", "test-version", mockClient)

		pullRequestChanges := []*entities.PullRequestChanges{
			{
				FileName: "test.go",
				Changed:  "@@ -1,1 +1,1 @@\n-test\n+test2",
			},
		}

		mockClient.On("GetSuggestions", pullRequestChanges).Return("", errors.New("mock error"))

		suggestions, err := adapter.GetCodeSuggestions(pullRequestChanges)

		assert.Error(t, err)
		assert.Nil(t, suggestions)

		mockClient.AssertExpectations(t)
	})
}
