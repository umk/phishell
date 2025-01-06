package session

import (
	"sync"
	"time"
)

type Inbox struct {
	messages sync.Map
}

type InboxMessage struct {
	ID      string    // Unique ID of the message used for deduplication
	Content string    // Content of the message
	Wd      string    // Working directory to use in connection with the message
	Date    time.Time // Date and time when the message was created by provider
}

func (i *Inbox) Store(m *InboxMessage) {
	i.messages.Store(m.ID, m)
}

func (i *Inbox) Delete(id string) {
	i.messages.Delete(id)
}

func (i *Inbox) Messages() []*InboxMessage {
	var messages []*InboxMessage
	i.messages.Range(func(key, value interface{}) bool {
		messages = append(messages, value.(*InboxMessage))
		return true
	})

	return messages
}

func (i *Inbox) Len() int {
	size := 0
	i.messages.Range(func(key, value interface{}) bool {
		size++
		return true
	})

	return size
}
