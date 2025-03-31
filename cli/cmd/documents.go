package cmd

import "github.com/umk/phishell/db/vector"

type documentsContext map[string]*documentsBatch

type documentsBatch struct {
	chunks []vector.ID
}
