package client

import (
	"sync"

	"golang.org/x/sync/semaphore"
)

var mu sync.Mutex

var clients = make(map[*Ref]*Client)

func Get(ref *Ref) *Client {
	mu.Lock()
	defer mu.Unlock()

	client, ok := clients[ref]
	if !ok {
		s := semaphore.NewWeighted(int64(ref.Config.Concurrency))
		client = &Client{
			Ref: ref,

			s:       s,
			Samples: newSamples(samplesCount, defaultBytesPerTok),
		}

		clients[ref] = client
	}

	return client
}
