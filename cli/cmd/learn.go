package cmd

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"strings"
	"sync"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/client"
	"github.com/umk/phishell/db"
	"github.com/umk/phishell/db/vector"
	"github.com/umk/phishell/splitter"
	"github.com/umk/phishell/util/execx"
	"github.com/umk/phishell/util/fsx"
)

type LearnCommand struct {
	context *Context
}

func (c *LearnCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) < 1 {
		return ErrInvalidArgs
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	var includes, excludes []string
	for _, pattern := range args {
		if strings.HasPrefix(pattern, "!") {
			excludes = append(excludes, pattern[1:])
		} else {
			includes = append(includes, pattern)
		}
	}

	if len(includes) == 0 {
		return ErrInvalidArgs
	}

	var paths []string
	if err := fsx.GlobsWalk(workingDir, includes, excludes, func(path string, d fs.DirEntry) error {
		if !d.IsDir() {
			paths = append(paths, path)
		}
		return nil
	}); err != nil {
		return err
	}

	currentCtx, cancel := context.WithCancel(ctx)

	chunks := c.getChunks(currentCtx, paths)
	records := c.getRecords(currentCtx, chunks)

	var pendingRecords []db.Record[db.Document]
	for record := range records {
		if err, ok := record.(error); ok {
			cancel()
			return err
		}

		pendingRecords = append(pendingRecords, record.(db.Record[db.Document]))
	}

	var batch documentsBatch
	for _, r := range pendingRecords {
		r = db.DocumentDB.Add(r)
		batch.chunks = append(batch.chunks, r.ID)
	}

	cancel()

	return nil
}

func (c *LearnCommand) getChunks(ctx context.Context, paths []string) <-chan any {
	splitters := splitter.NewSplitters()
	chunks := make(chan any)

	go func() {
		var wg sync.WaitGroup
		for _, p := range paths {
			wg.Add(1)
			go func(p string) {
				s := splitters.GetSplitter(p)
				b, err := os.ReadFile(p)
				if err != nil {
					chunks <- err
				} else {
					c, err := s.Split(ctx, b, splitter.Metadata{})
					if err != nil {
						chunks <- err
					} else {
						chunks <- c
					}
				}

				wg.Done()
			}(p)
		}

		wg.Wait()

		close(chunks)
	}()

	return chunks
}

func (c *LearnCommand) getRecords(ctx context.Context, chunks <-chan any) <-chan any {
	records := make(chan any)

	go func() {
		var wg sync.WaitGroup
		for chunk := range chunks {
			if err, ok := chunk.(error); ok {
				records <- err
				break
			}

			wg.Add(1)
			go func() {
				c := chunk.(splitter.Chunk)
				content := c.Range.Get(c.Document)

				res, err := client.Default.Embeddings(ctx, openai.EmbeddingNewParams{
					Model: client.Default.Model(client.Tier2),
					Input: openai.EmbeddingNewParamsInputUnion{
						OfString: openai.String(content),
					},
				})
				if err != nil {
					records <- err
				} else if len(res.Data) != 1 {
					records <- errors.New("service didn't return an embedding")
				} else {
					e := res.Data[0].Embedding
					record := db.Record[db.Document]{
						Vector: make(vector.Vector, len(e)),
						Data:   db.Document{Chunk: c},
					}
					for i, v := range e {
						record.Vector[i] = float32(v)
					}
					records <- record
				}

				wg.Done()
			}()
		}

		wg.Wait()

		close(records)
	}()

	return records
}

func (c *LearnCommand) Usage() []string {
	return []string{"learn [pattern] ..."}
}

func (c *LearnCommand) Info() []string {
	return []string{"find and learn documents by pattern"}
}
