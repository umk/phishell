package thread

import (
	"fmt"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/util/termx"
)

func (t *Thread) processChatMessage(response openai.ChatCompletionMessage) error {
	if response.Refusal != "" {
		return fmt.Errorf("refused: %s", response.Refusal)
	}

	t.history.Messages = append(t.history.Messages, response)

	termx.Response.Println(response.Content)

	return nil
}
