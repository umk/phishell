package prompt

import (
	"fmt"

	"github.com/openai/openai-go"
)

type Completion openai.ChatCompletion

func (c *Completion) Content() (string, error) {
	if len(c.Choices) != 1 {
		return "", fmt.Errorf("unable to choose from %d choices", len(c.Choices))
	}

	choice := c.Choices[0]

	if choice.Message.Refusal != "" {
		return "", fmt.Errorf("refused: %s", choice.Message.Refusal)
	}

	return choice.Message.Content, nil
}

func (c *Completion) Toks() int {
	return int(c.Usage.CompletionTokens)
}

func (c *Completion) TotalToks() int {
	return int(c.Usage.TotalTokens)
}
