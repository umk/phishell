package session

import (
	"context"
	"time"
)

func (s *Session) ProcessEvent(ctx context.Context, id, content, wd string, date time.Time) error {
	s.Inbox.Store(&InboxMessage{
		ID:      id,
		Content: content,
		Wd:      wd,
		Date:    date,
	})

	return nil
}
