package client

import (
	"sync"

	"github.com/umk/phishell/bootstrap"
	"golang.org/x/sync/semaphore"
)

var mu sync.Mutex

var clients = make(map[*bootstrap.Profile]*Client)

func Get(clientRef *bootstrap.Profile) *Client {
	mu.Lock()
	defer mu.Unlock()

	client, ok := clients[clientRef]
	if !ok {
		s := semaphore.NewWeighted(int64(clientRef.Config.Concurrency))
		client = &Client{
			Profile: clientRef,

			s:       s,
			Samples: newSamples(samplesCount, defaultBytesPerTok),
		}

		clients[clientRef] = client
	}

	return client
}
