package thread

import (
	"slices"

	"github.com/openai/openai-go"
)

type History struct {
	Messages []openai.ChatCompletionMessageParamUnion

	// A collection of pending messages to be delivered with the next
	// user message.
	Pending []string
}

// Reset returns only the part of the history that contains the system message, if any.
func (h *History) Reset() *History {
	if len(h.Messages) > 0 {
		first := h.Messages[0]

		if _, ok := first.(openai.ChatCompletionSystemMessageParam); ok {
			return &History{
				Messages: []openai.ChatCompletionMessageParamUnion{first},
			}
		}
	}

	return &History{}
}

func (h *History) Clone() *History {
	return &History{
		Messages: slices.Clone(h.Messages),
		Pending:  slices.Clone(h.Pending),
	}
}
