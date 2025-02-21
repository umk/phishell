package session

import (
	"errors"
	"os"

	"github.com/umk/phishell/util/execx"
)

func (s *Session) Attach(cmd string) error {
	args, err := execx.Parse(cmd)
	if err != nil {
		return err
	}

	if len(args) != 1 {
		return errors.New("cannot pipe from or to the tools provider")
	}

	if _, err := s.Host.Execute(&execx.Cmd{
		Cmd:  args[0][0],
		Args: args[0][1:],
		Env:  append(os.Environ(), "PHI_SHELL=1"),
	}); err != nil {
		return err
	}

	return nil
}
