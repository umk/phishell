package thread

import (
	"github.com/openai/openai-go"
	"github.com/umk/phishell/util/termx"
)

func (t *Thread) processChatMessage(response openai.ChatCompletionMessage) error {
	t.frame.Messages = append(t.frame.Messages, response)
	t.frame.Response = response.Content

	termx.MD.Print(response.Content)

	return nil
}
