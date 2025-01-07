package thread

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/cli/msg"
	"github.com/umk/phishell/cli/prompt"
	"github.com/umk/phishell/cli/tool/host"
	"github.com/umk/phishell/util/errorsx"
	"github.com/umk/phishell/util/termx"
)

type Thread struct {
	history *History

	client *bootstrap.ClientRef
	host   *host.Host

	tools []openai.ChatCompletionToolParam
}

func NewThread(
	history *History,
	client *bootstrap.ClientRef,
	host *host.Host,
) (*Thread, error) {
	return &Thread{
		history: history.Clone(),

		client: client,
		host:   host,

		// Save for consistency across the rounds of LLM calls even if some of
		// the tools become unavailable.
		tools: host.Tools(),
	}, nil
}

func (t *Thread) Post(ctx context.Context, message string) (*History, error) {
	if err := t.compactHistory(ctx); err != nil {
		return nil, fmt.Errorf("history compaction failed: %w", err)
	}

	message, err := msg.FormatUserMessage(&msg.UserMessageParams{
		Request: message,
		Context: t.history.Pending,
	})
	if err != nil {
		return nil, err
	}

	t.history.Messages = append(t.history.Messages, openai.UserMessage(message))
	t.history.Pending = nil

	for retries := 0; ; {
		message, err := prompt.PromptUser(ctx, &prompt.UserPromptParams{
			Client:   t.client,
			Messages: t.history.Messages,
			Tools:    t.tools,
		})
		if err != nil {
			return nil, fmt.Errorf("service request failed: %w", err)
		}

		if len(message.Choices) > 0 {
			response := message.Choices[0].Message

			if len(response.ToolCalls) == 0 {
				if err := t.processChatMessage(response); err != nil {
					return nil, err
				}

				t.history.Toks = message.Usage.TotalTokens
				return t.history, nil
			}

			if err := t.processToolMessage(ctx, response, &retries); err != nil {
				if errorsx.IsCanceled(err) {
					return nil, err
				}
				termx.Error.Println(err)
			}
		}
	}
}
