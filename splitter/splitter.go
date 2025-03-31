package splitter

import (
	"context"
	"path/filepath"
)

type Splitter interface {
	Split(ctx context.Context, document []byte, metadata Metadata) ([]Chunk, error)
}

type Splitters struct {
	Markdown *MarkdownSplitter
}

func NewSplitters() *Splitters {
	return &Splitters{
		Markdown: NewMarkdownSplitter(),
	}
}

func (s *Splitters) GetSplitter(path string) Splitter {
	ext := filepath.Ext(path)
	switch ext {
	case ".md":
		return s.Markdown
	default:
		return nil
	}
}
