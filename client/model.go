package client

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
)

type ModelTier int

const (
	Tier1 ModelTier = iota + 1
	Tier2
)

func (c Client) GetModel(ctx context.Context, tier ModelTier) string {
	if c.Config.Model != "" {
		return c.Config.Model
	}

	switch tier {
	case Tier1:
		return openai.ChatModelGPT4o
	case Tier2:
		return openai.ChatModelGPT4oMini
	default:
		panic(fmt.Sprintf("model tier is not supported: %d", tier))
	}
}
