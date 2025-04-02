package db

import "github.com/umk/phishell/splitter"

var DocumentDB = NewDatabase[Document]()

type Document struct {
	Chunk splitter.Chunk
}
