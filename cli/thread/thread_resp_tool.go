package thread

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/util/errorsx"
)

func (t *Thread) processToolMessage(ctx context.Context, response openai.ChatCompletionMessage, retries *int) error {
	r := NewToolRunner(t.host)

	for _, call := range response.ToolCalls {
		if err := r.Add(&call); err != nil {
			if errorsx.IsCanceled(err) {
				return err
			}
			if errorsx.IsRetryable(err) && *retries < t.client.Config.Retries {
				*retries++

				// A returned error indicates that the retry loop must continue.
				// Otherwise the error gets straight into the messages and supposed
				// to be handled by the model.
				return fmt.Errorf("tool %s call failed: %w", call.Function.Name, err)
			}
		}
	}

	// The next iteration is charged with all retries.
	*retries = 0

	if response.Content != "" {
		fmt.Println(response.Content)
	}

	messages, err := r.Complete(ctx)
	if err != nil {
		return err
	}

	t.history.Messages = append(t.history.Messages, response)
	t.history.Messages = append(t.history.Messages, messages...)

	return nil
}
