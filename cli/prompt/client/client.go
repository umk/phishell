package client

import (
	"context"
	"sync"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/ssestream"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/util/termx"
	"golang.org/x/sync/semaphore"
)

var mu sync.Mutex

var clients = make(map[*bootstrap.ClientRef]*Client)

func Get(clientRef *bootstrap.ClientRef) *Client {
	mu.Lock()
	defer mu.Unlock()

	client, ok := clients[clientRef]
	if !ok {
		s := semaphore.NewWeighted(int64(clientRef.Config.Concurrency))
		client = &Client{
			ClientRef: clientRef,
			s:         s,
		}

		clients[clientRef] = client
	}

	return client
}

type Client struct {
	*bootstrap.ClientRef

	s *semaphore.Weighted
}

type StreamingCallback func(chunk openai.ChatCompletionChunk) error

func (c Client) Completion(ctx context.Context, params openai.ChatCompletionNewParams) (
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
		messages := len(params.Messages.Value)
		promptToks := compl.Usage.PromptTokens
		complToks := compl.Usage.CompletionTokens
		totalToks := compl.Usage.TotalTokens

		if bootstrap.IsDebug(ctx) {
			termx.Muted.Printf("(messages=%d; prompt=%d; completion=%d; total=%d)\n",
				messages, promptToks, complToks, totalToks,
			)
		}
	}

	return compl, err
}

func (c Client) CompletionStreaming(
	ctx context.Context, params openai.ChatCompletionNewParams, cb StreamingCallback,
) error {
	c.s.Acquire(ctx, 1)
	defer c.s.Release(1)

	paramsCur := params
	paramsCur.StreamOptions.Value.IncludeUsage = openai.F(bootstrap.IsDebug(ctx))

	var promptToks, complToks, totalToks int64
	err := c.Request(ctx, func(client *openai.Client) (err error) {
		stream := func() *ssestream.Stream[openai.ChatCompletionChunk] {
			termx.SpinnerStart()
			defer termx.SpinnerStop()

			return client.Chat.Completions.NewStreaming(ctx, paramsCur)
		}()

		for stream.Next() {
			cur := stream.Current()

			promptToks += cur.Usage.PromptTokens
			complToks += cur.Usage.CompletionTokens
			totalToks += cur.Usage.TotalTokens

			if err := cb(cur); err != nil {
				return err
			}
		}
		return stream.Err()
	})

	if bootstrap.IsDebug(ctx) {
		messages := len(params.Messages.Value)
		termx.Muted.Printf("\n(messages=%d; prompt=%d; completion=%d; total=%d)",
			messages, promptToks, complToks, totalToks,
		)
	}

	return err
}
