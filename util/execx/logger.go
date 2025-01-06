package execx

import (
	"errors"
	"os/exec"
)

// CmdLogger captures the Stdout and Stderr output of a process,
// writes it to the console and exposes functions to returns the most
// relevant output based on the status code:
//
//   - stdout if the process was successful
//   - stderr if the process failed and produced an error message
//   - stdout if the process failed but there is no stderr output
type CmdLogger struct {
	cmd *exec.Cmd
}

func Log(cmd *exec.Cmd, maxLen int) *CmdLogger {
	cmd.Stdout = newOutputWrapper(cmd.Stdout, maxLen)
	cmd.Stderr = newOutputWrapper(cmd.Stderr, maxLen)

	return &CmdLogger{cmd}
}

// Output returns the most relevant output based on the exit code.
func (l *CmdLogger) Output() (ProcessOutput, error) {
	if l.cmd.ProcessState == nil {
		return ProcessOutput{}, errors.New("process has not finished yet")
	}

	exitCode := l.cmd.ProcessState.ExitCode()

	if exitCode != 0 {
		w, ok := l.cmd.Stderr.(*outputWrapper)
		if !ok {
			return ProcessOutput{}, errors.New("invalid stderr")
		}

		return w.Get(), nil
	}

	w, ok := l.cmd.Stdout.(*outputWrapper)
	if !ok {
		return ProcessOutput{}, errors.New("invalid stdout")
	}

	return w.Get(), nil
}
