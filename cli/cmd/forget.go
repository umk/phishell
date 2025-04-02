package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/umk/phishell/db"
	"github.com/umk/phishell/util/execx"
)

type ForgetCommand struct {
	context *Context
}

func (c *ForgetCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) != 1 {
		return ErrInvalidArgs
	}

	batchID, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	b, ok := c.context.documents.batches[batchID]
	if !ok {
		return fmt.Errorf("no such batch with ID %d", batchID)
	}

	db.DocumentDB.DeleteBatch(b.chunks)

	fmt.Println("OK")
	return nil
}

func (c *ForgetCommand) Usage() []string {
	return []string{"forget [batch]"}
}

func (c *ForgetCommand) Info() []string {
	return []string{"forget the previously learned batch"}
}
