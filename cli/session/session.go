package session

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/umk/phishell/cli/thread"
	"github.com/umk/phishell/cli/tool/host"
	"github.com/umk/phishell/util/execx"
)

type Session struct {
	Host *host.Host

	History *thread.History
	Inbox   Inbox

	// Output of a previously executed command.
	PreviousOut *PreviousOut
}

type PreviousOut struct {
	CommandLine string
	ExitCode    int
	Output      execx.ProcessOutput
}

func NewSession() *Session {
	session := &Session{}

	session.Host = host.NewHost(session)

	return session
}

func (s *Session) Chdir(dir string) error {
	dir, err := s.Resolve(dir)
	if err != nil {
		return err
	}

	err = os.Chdir(dir)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) Resolve(path string) (string, error) {
	if strings.HasPrefix(path, "./") || strings.HasPrefix(path, ".\\") {
		// Relative path
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		return absPath, nil
	}

	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		// Path from home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		relPath := path[2:] // Remove ~/
		absPath := filepath.Join(homeDir, relPath)
		return absPath, nil
	}

	return path, nil
}