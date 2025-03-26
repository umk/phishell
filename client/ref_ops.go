package client

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/config"
	"github.com/umk/phishell/util/termx"
)

const (
	defaultBytesPerTok = 3.25
	samplesCount       = 5
	minSampleSize      = 100
)

func (ref *Ref) Completion(ctx context.Context, params openai.ChatCompletionNewParams) (
	*openai.ChatCompletion, error,
) {
	termx.Spinner.Start()
	defer termx.Spinner.Stop()

	if err := ref.S.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer ref.S.Release(1)

	c, err := ref.Client.Chat.Completions.New(ctx, params)

	if c != nil {
		ref.setSamplesFromCompl(c)

		messages := len(params.Messages)
		promptToks := c.Usage.PromptTokens
		complToks := c.Usage.CompletionTokens
		totalToks := c.Usage.TotalTokens
		bytesPerTok := ref.Samples.BytesPerTok()

		if config.Config.Debug {
			termx.Muted.Printf("(messages=%d; prompt=%d; completion=%d; total=%d; bytes per tok=%.2f)\n",
				messages, promptToks, complToks, totalToks, bytesPerTok,
			)
		}
	}

	return c, err
}

func (ref *Ref) Embeddings(ctx context.Context, params openai.EmbeddingNewParams) (
	*openai.CreateEmbeddingResponse, error,
) {
	termx.Spinner.Start()
	defer termx.Spinner.Stop()

	if err := ref.S.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer ref.S.Release(1)

	e, err := ref.Client.Embeddings.New(ctx, params)

	if e != nil {
		ref.setSamplesFromEmbedding(&params, e)

		promptToks := e.Usage.PromptTokens
		totalToks := e.Usage.TotalTokens
		bytesPerTok := ref.Samples.BytesPerTok()

		if config.Config.Debug {
			termx.Muted.Printf("(prompt=%d; total=%d; bytes per tok=%.2f)\n",
				promptToks, totalToks, bytesPerTok,
			)
		}
	}

	return e, err
}

func (ref *Ref) setSamplesFromCompl(c *openai.ChatCompletion) {
	t := c.Usage.CompletionTokens
	if t == 0 {
		return
	}

	var b int

	for _, c := range c.Choices {
		if c.Message.Refusal != "" {
			return
		}

		if len(c.Message.ToolCalls) > 0 {
			return
		}

		b += len(c.Message.Content)
	}

	if b >= minSampleSize {
		ref.Samples.put(float32(b) / float32(t))
	}
}

func (ref *Ref) setSamplesFromEmbedding(params *openai.EmbeddingNewParams, e *openai.CreateEmbeddingResponse) {
	t := e.Usage.PromptTokens
	if t == 0 {
		return
	}

	var b int

	if params.Input.OfArrayOfStrings != nil {
		for _, s := range params.Input.OfArrayOfStrings {
			b += len(s)
		}
	}

	if b >= minSampleSize {
		ref.Samples.put(float32(b) / float32(t))
	}
}
