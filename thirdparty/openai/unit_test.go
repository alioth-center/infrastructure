package openai

import (
	"github.com/alioth-center/infrastructure/logger"
	"os"
	"strings"
	"testing"
)

func initMockingClient(t *testing.T) Client {
	t.Helper()
	return nil
}

func TestOpenAiClient(t *testing.T) {
	// uses real openai endpoint to test, because mocking the endpoint cannot find issues in the client
	var client Client
	apiKey, baseUrl := os.Getenv("OPENAI_API_KEY"), os.Getenv("OPENAI_BASE_URL")
	if apiKey == "" || baseUrl == "" {
		t.Log("OPENAI_API_KEY or OPENAI_BASE_URL is not set, mock it, but it is not recommended")
		client = initMockingClient(t)
	} else {
		client = NewClient(Config{ApiKey: apiKey, BaseUrl: baseUrl}, logger.Default())
	}

	t.Run("CompleteChat", func(t *testing.T) {
		response, err := client.CompleteChat(CompleteChatRequest{
			Body: CompleteChatRequestBody{
				Model: "gpt-4o",
				Messages: []ChatMessageObject{
					{
						Role:    ChatRoleEnumSystem,
						Content: "now testing api is working, please echo any input",
					},
					{
						Role:    ChatRoleEnumUser,
						Content: "testing",
					},
				},
				N: 1,
			},
		})

		if err != nil {
			t.Error(err)
		}

		if len(response.Choices) == 0 || !strings.Contains(response.Choices[0].Message.Content, "testing") {
			t.Error("response is not as expected")
		}
	})
}
