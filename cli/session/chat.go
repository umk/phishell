package session

import (
	"context"

	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/thread"
)

func (s *Session) Post(ctx context.Context, client *bootstrap.ClientRef, content string) error {
	t, err := thread.NewThread(s.History, client, s.Host)
	if err != nil {
		return err
	}

	history, err := t.Post(ctx, content)
	if err != nil {
		return err
	}

	s.History = history

	return nil
}
