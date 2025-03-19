package tool

import (
	"context"
)

// Handler defines an interface for executing tools.
type Handler interface {
	Execute(ctx context.Context) (any, error)
}

type DescriptionSource interface {
	// Describe gets a string that describes, what the tool is about
	// to do when executed with given arguments.
	Describe(ctx context.Context) (string, error)
}

// Describe gets a tool description, or an empty string if the tool doesn't
// provide additional description.
func Describe(ctx context.Context, h Handler) (string, error) {
	if s, ok := h.(DescriptionSource); ok {
		return s.Describe(ctx)
	}

	return "", nil
}
