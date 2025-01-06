package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/umk/phishell/util/execx"
	"gopkg.in/yaml.v3"
)

type HistoryCommand struct {
	context *Context
}

func (c *HistoryCommand) Execute(ctx context.Context, args execx.Arguments) error {
	// Using JSON serialization and deserialization to leverage the client library
	// own JSON serialization logic to extract relevant values.
	j, err := json.Marshal(c.context.session.History.Messages)
	if err != nil {
		return err
	}

	var v any
	if err := json.Unmarshal(j, &v); err != nil {
		return err
	}

	// Then using YAML marshaling for pretty printing the messages.
	y, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	fmt.Println(string(y))

	for _, m := range c.context.session.History.Pending {
		fmt.Println("---")
		fmt.Println(m)
	}

	return nil
}

func (k *HistoryCommand) Info() string {
	return "history: display the chat messages history"
}
