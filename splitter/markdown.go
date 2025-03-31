package splitter

import (
	"context"

	"github.com/umk/phishell/client"
)

type MarkdownSplitter struct {
	chunkSize, overlapSize int
}

func NewMarkdownSplitter() *MarkdownSplitter {
	var bytesPerTok float32 = client.Default.Samples.BytesPerTok()

	indexing := client.Default.Config.Indexing

	return &MarkdownSplitter{
		chunkSize:   int(float32(indexing.ChunkToks) * bytesPerTok),
		overlapSize: int(float32(indexing.OverlapToks) * bytesPerTok),
	}
}

func (ms *MarkdownSplitter) Split(ctx context.Context, document []byte, metadata Metadata) ([]Chunk, error) {
	return nil, nil
}
