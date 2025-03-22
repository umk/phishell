package client

import (
	"fmt"

	"github.com/openai/openai-go"
)

type ModelTier int

const (
	Tier1 ModelTier = iota + 1
	Tier2
)

func (ref *Ref) Model(tier ModelTier) string {
	if ref.Config.Model != "" {
		return ref.Config.Model
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
