package execx

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

func Pipe(cmds []*exec.Cmd, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(cmds) == 0 {
		panic("pipe must contain at least one command")
	}

	cmds[0].Stdin = stdin

	// Connect the commands
	for i := 0; i < len(cmds)-1; i++ {
		stdout, err := cmds[i].StdoutPipe()
		if err != nil {
			return fmt.Errorf("error creating stdout pipe: %w", err)
		}
		cmds[i+1].Stdin = stdout
	}

	// Set the output of the last command
	cmds[len(cmds)-1].Stdout = stdout
	cmds[len(cmds)-1].Stderr = stderr

	return nil
}

func RunPipe(cmds []*exec.Cmd) (int, error) {
	if len(cmds) == 0 {
		panic("pipe must contain at least one command")
	}

	// Start the commands
	var wg sync.WaitGroup
	var setErr sync.Once
	var cmdErr error

	for _, cmd := range cmds {
		err := cmd.Start()
		if err != nil {
			setErr.Do(func() {
				if cmdErr == nil {
					cmdErr = err
				}
			})
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := cmd.Wait(); err != nil {
				var exitErr *exec.ExitError
				if !errors.As(cmdErr, &exitErr) {
					setErr.Do(func() {
						cmdErr = err
					})
				}
			}
		}()
	}

	wg.Wait()

	if cmdErr != nil {
		return -1, cmdErr
	}

	lastCmd := cmds[len(cmds)-1]
	exitCode := lastCmd.ProcessState.ExitCode()

	return exitCode, nil
}
