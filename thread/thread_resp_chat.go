package thread

import (
	"github.com/openai/openai-go"
	"github.com/umk/phishell/util/termx"
)

func (t *Thread) processChatMessage(response openai.ChatCompletionMessage) error {
	message := response.ToAssistantMessageParam()
	t.frame.Messages = append(t.frame.Messages, openai.ChatCompletionMessageParamUnion{
		OfAssistant: &message,
	})
	t.frame.Response = response.Content

	termx.MD.Print(response.Content)

	return nil
}
