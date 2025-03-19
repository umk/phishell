package thread

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/tool"
	"github.com/umk/phishell/util/errorsx"
)

func (t *Thread) processToolMessage(ctx context.Context, response openai.ChatCompletionMessage) error {
	r := NewToolRunner(t.host)

	for _, call := range response.ToolCalls {
		var functionDescr string
		if t, ok := t.host.Tool(call.Function.Name); ok {
			if f, ok := t.Function.Raw.(tool.Function); ok {
				functionDescr = f.Description
			}
		}

		if err := r.Add(&call, functionDescr); err != nil {
			if errorsx.IsCanceled(err) {
				return err
			}
		}
	}

	if response.Content != "" {
		fmt.Println(response.Content)
	}

	messages, err := r.Complete(ctx)
	if err != nil {
		return err
	}

	t.frame.Messages = append(t.frame.Messages, response)
	t.frame.Messages = append(t.frame.Messages, messages...)

	return nil
}
