package cmd

import "github.com/umk/phishell/db/vector"

type documentsContext struct {
	currentID int
	batches   map[int]*documentsBatch
}

type documentsBatch struct {
	chunks []vector.ID
}

func makeDocumentsContext() documentsContext {
	return documentsContext{
		currentID: 0,
		batches:   make(map[int]*documentsBatch),
	}
}
