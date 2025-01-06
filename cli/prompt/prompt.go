package prompt

import (
	"fmt"

	"github.com/openai/openai-go"
)

func getCompletionContent(c *openai.ChatCompletion) (string, error) {
	if len(c.Choices) != 1 {
		return "", fmt.Errorf("unable to choose from %d choices", len(c.Choices))
	}

	choice := c.Choices[0]

	if choice.Message.Refusal != "" {
		return "", fmt.Errorf("refused: %s", choice.Message.Refusal)
	}

	return choice.Message.Content, nil
}
