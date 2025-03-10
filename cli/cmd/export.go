package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/umk/phishell/util/execx"
)

type ExportCommand struct{}

func (c *ExportCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) > 1 {
		return getUsageError(c)
	}

	if len(args) == 1 {
		return c.setEnvVariable(args[0])
	}

	c.showEnvVariables()

	return nil
}

func (c *ExportCommand) Usage() []string {
	return []string{"export <var>"}
}

func (c *ExportCommand) Info() []string {
	return []string{"set value of an environment variable"}
}

func (c *ExportCommand) showEnvVariables() {
	var names []string
	vars := make(map[string]string)

	for _, v := range os.Environ() {
		s := strings.SplitN(v, "=", 2)
		if len(s) == 2 {
			names = append(names, s[0])
			vars[s[0]] = s[1]
		}
	}

	sort.Strings(names)

	for _, n := range names {
		fmt.Printf("%s=%s\n", n, vars[n])
	}
}

func (c *ExportCommand) setEnvVariable(arg string) error {
	s := strings.SplitN(arg, "=", 2)

	if len(s) < 2 {
		return fmt.Errorf("missing value for the environment variable")
	}

	if err := os.Setenv(s[0], s[1]); err != nil {
		return err
	}

	fmt.Println("OK")

	return nil
}
