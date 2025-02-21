package thread

import (
	"github.com/openai/openai-go"
)

func (t *Thread) processChatMessage(response openai.ChatCompletionMessage) error {
	t.frame.Messages = append(t.frame.Messages, response)
	t.frame.Response = response.Content

	t.printer.Print(response.Content)

	return nil
}
