package thread

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/client"
	"github.com/umk/phishell/prompt"
	"github.com/umk/phishell/prompt/msg"
	"github.com/umk/phishell/tool/host"
	"github.com/umk/phishell/util/errorsx"
)

type Thread struct {
	history *History
	frame   *MessagesFrame

	client *client.Ref
	host   *host.Host

	tools []openai.ChatCompletionToolParam
}

func NewThread(history *History, ref *client.Ref, host *host.Host) (*Thread, error) {
	tools, err := host.Tools()
	if err != nil {
		return nil, err
	}

	h := history.Clone()
	h.Frames = append(h.Frames, MessagesFrame{})

	return &Thread{
		history: h,
		frame:   &h.Frames[len(h.Frames)-1],

		client: ref,
		host:   host,

		// Save for consistency across the rounds of LLM calls even if some of
		// the tools become unavailable.
		tools: tools,
	}, nil
}

func (t *Thread) Post(ctx context.Context, message string) (*History, error) {
	message, err := msg.FormatUserMessage(&msg.UserMessageParams{
		Request: message,
		Context: t.history.Pending,
	})
	if err != nil {
		return nil, err
	}

	t.frame.Messages = append(t.frame.Messages, openai.UserMessage(message))
	t.frame.Request = message

	t.history.Pending = nil

	sys, err := msg.FormatSystemMessage(&msg.SystemMessageParams{
		Prompt: t.client.Config.Prompt,
		OS:     runtime.GOOS,
	})
	if err != nil {
		return nil, err
	}

	for {
		response, err := prompt.PromptUser(ctx, &prompt.UserPromptParams{
			Messages: t.history.Messages(openai.SystemMessage(sys)),
			Tools:    t.tools,
			Client:   t.client,
		})
		if err != nil {
			return nil, fmt.Errorf("service request failed: %w", err)
		}

		if len(response.Choices) == 0 {
			return nil, fmt.Errorf("invalid response")
		}

		responseMsg, err := response.Message()
		if err != nil {
			return nil, nil
		}

		t.frame.Toks += response.RequestToks() + response.ResponseToks()

		if len(responseMsg.ToolCalls) == 0 {
			if err := t.processChatMessage(responseMsg); err != nil {
				return nil, err
			}

			return t.history, nil
		}

		if err := t.processToolMessage(ctx, responseMsg); err != nil {
			if errorsx.IsCanceled(err) {
				return nil, err
			}
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
