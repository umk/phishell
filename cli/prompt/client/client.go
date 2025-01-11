package client

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/util/termx"
	"golang.org/x/sync/semaphore"
)

const (
	defaultBytesPerTok = 3.25
	samplesCount       = 5
	minSampleSize      = 100
)

type Client struct {
	*bootstrap.ClientRef

	s       *semaphore.Weighted
	Samples *Samples
}

func (c *Client) Completion(ctx context.Context, params openai.ChatCompletionNewParams) (
	*openai.ChatCompletion, error,
) {
	c.s.Acquire(ctx, 1)
	defer c.s.Release(1)

	var compl *openai.ChatCompletion
	err := c.Request(ctx, func(client *openai.Client) (err error) {
		termx.SpinnerStart()
		defer termx.SpinnerStop()

		compl, err = client.Chat.Completions.New(ctx, params)
		return
	})

	if compl != nil {
		c.setSamplesFromCompl(compl)

		messages := len(params.Messages.Value)
		promptToks := compl.Usage.PromptTokens
		complToks := compl.Usage.CompletionTokens
		totalToks := compl.Usage.TotalTokens
		bytesPerTok := c.Samples.BytesPerTok()

		if bootstrap.IsDebug(ctx) {
			termx.Muted.Printf("(messages=%d; prompt=%d; completion=%d; total=%d; bytes per tok=%.2f)\n",
				messages, promptToks, complToks, totalToks, bytesPerTok,
			)
		}
	}

	return compl, err
}

func (c *Client) setSamplesFromCompl(compl *openai.ChatCompletion) {
	t := compl.Usage.CompletionTokens
	if t == 0 {
		return
	}

	var b int

	for _, c := range compl.Choices {
		if c.Message.Refusal != "" {
			return
		}

		if len(c.Message.ToolCalls) > 0 {
			return
		}

		b += len(c.Message.Content)
	}

	if b < minSampleSize {
		return
	}

	bytesPerTok := float32(b) / float32(t)

	c.Samples.put(bytesPerTok)
}
