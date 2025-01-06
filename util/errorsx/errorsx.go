package errorsx

import (
	"context"
	"errors"
)

var ErrInterrupted = errors.New("operation interrupted")

func IsCanceled(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, ErrInterrupted)
}
