package db

import "github.com/umk/phishell/splitter"

var DocumentDB = NewDb[Document]()

type Document struct {
	Chunk splitter.Chunk
}
