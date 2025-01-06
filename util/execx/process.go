package execx

import (
	"errors"
	"os/exec"
)

func Run(cmd *exec.Cmd) (int, error) {
	if err := cmd.Start(); err != nil {
		return -1, err
	}

	if err := cmd.Wait(); err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			return -1, err
		}
	}

	exitCode := cmd.ProcessState.ExitCode()

	return exitCode, nil
}
