package prompt

import (
	"errors"
	"fmt"

	"github.com/openai/openai-go"
)

type Completion openai.ChatCompletion

func (c *Completion) Content() (string, error) {
	message, err := c.Message()
	if err != nil {
		return "", err
	}

	if len(message.ToolCalls) > 0 {
		return "", errors.New("unexpected tool calls")
	}

	return message.Content, nil
}

func (c *Completion) Message() (openai.ChatCompletionMessage, error) {
	if len(c.Choices) != 1 {
		return openai.ChatCompletionMessage{}, fmt.Errorf("unable to choose from %d choices", len(c.Choices))
	}

	choice := c.Choices[0]

	if choice.Message.Refusal != "" {
		return openai.ChatCompletionMessage{}, fmt.Errorf("refused: %s", choice.Message.Refusal)
	}

	return choice.Message, nil
}

func (c *Completion) RequestToks() int {
	return int(c.Usage.PromptTokens)
}

func (c *Completion) ResponseToks() int {
	return int(c.Usage.CompletionTokens)
}

func (c *Completion) TotalToks() int {
	return int(c.Usage.TotalTokens)
}
