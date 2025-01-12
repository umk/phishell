package thread

import (
	"slices"

	"github.com/openai/openai-go"
)

type History struct {
	Messages []openai.ChatCompletionMessageParamUnion

	// Total number of tokens in the history. Doesn't include the
	// pending messages.
	Toks int

	// A collection of pending messages to be delivered with the next
	// user message.
	Pending []string
}

// Reset returns only the part of the history that contains the system message, if any.
func Reset(h *History) *History {
	if s := getSystemMessage(h); s != nil {
		return &History{
			Messages: []openai.ChatCompletionMessageParamUnion{s},
		}
	}

	return &History{}
}

func Clone(h *History) *History {
	return &History{
		Messages: slices.Clone(h.Messages),
		Toks:     h.Toks,

		Pending: slices.Clone(h.Pending),
	}
}

func getSystemMessage(h *History) *openai.ChatCompletionSystemMessageParam {
	if len(h.Messages) == 0 {
		return nil
	}

	m := h.Messages[0]

	if s, ok := m.(openai.ChatCompletionSystemMessageParam); ok {
		return &s
	}

	return nil
}

// getSystemMessageSize gets the system message size in bytes.
func getSystemMessageSize(h *History) int {
	if s := getSystemMessage(h); s != nil {
		var r int

		for _, v := range s.Content.Value {
			r += len(v.Text.Value)
		}

		return r
	}

	return 0
}
