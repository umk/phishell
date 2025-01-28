package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type temporaryLog struct {
	*os.File
}

func (l *temporaryLog) Close() error {
	if err := l.File.Close(); err != nil {
		return err
	}

	if err := os.Remove(l.Name()); err != nil {
		return err
	}

	return nil
}

func setupLogging(name string) (io.Closer, error) {
	l, err := createLog(name)
	if err != nil {
		return nil, err
	}

	const logFileVar = "PHI_LOG"

	if name != "" {
		os.Setenv(logFileVar, name)
	} else if f, ok := l.(interface{ Name() string }); ok {
		os.Setenv(logFileVar, f.Name())
	}

	log.SetOutput(l)

	return l, nil
}

func createLog(name string) (io.WriteCloser, error) {
	if name == "" {
		logName := fmt.Sprintf("phishell.%d.*.log", os.Getpid())
		f, err := os.CreateTemp("", logName)
		if err != nil {
			return nil, err
		}

		return &temporaryLog{f}, nil
	}

	fileInfo, err := os.Stat(name)
	if err != nil {
		return nil, fmt.Errorf("failed to stat path: %w", err)
	}

	mode := fileInfo.Mode()

	if mode&os.ModeSocket != 0 {
		conn, err := net.Dial("unix", name)
		if err != nil {
			return nil, fmt.Errorf("failed to open socket: %w", err)
		}
		return conn, nil
	}

	if mode&os.ModeNamedPipe != 0 {
		file, err := os.OpenFile(name, os.O_WRONLY, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to open named pipe for writing: %w", err)
		}
		return file, nil
	}

	if mode.IsRegular() {
		file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return nil, fmt.Errorf("failed to open file for writing: %w", err)
		}
		return file, nil
	}

	return nil, fmt.Errorf("unsupported file type at path: %s", name)
}
