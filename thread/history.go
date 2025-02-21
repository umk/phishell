package thread

import (
	"slices"

	"github.com/openai/openai-go"
)

type History struct {
	Frames []MessagesFrame

	// A collection of pending messages to be delivered with the next
	// user message.
	Pending []string
}

type MessagesFrame struct {
	Request  string
	Response string

	Messages []openai.ChatCompletionMessageParamUnion

	Toks int
}

func (h *History) Messages(current ...openai.ChatCompletionMessageParamUnion) []openai.ChatCompletionMessageParamUnion {
	var r []openai.ChatCompletionMessageParamUnion
	for _, f := range h.Frames {
		r = append(r, f.Messages...)
	}
	return append(r, current...)
}

func (h *History) Clone() *History {
	return &History{
		Frames:  slices.Clone(h.Frames),
		Pending: slices.Clone(h.Pending),
	}
}
