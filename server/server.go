package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/umk/phishell/server/internal"
)

var instance *http.Server

func Init() error {
	if instance != nil {
		panic("server already initialized")
	}

	socket, server, err := serve(nil)
	if err != nil {
		return err
	}

	instance = server

	os.Setenv("PHISHELL_SOCKET", socket)
	return nil
}

// Close shuts down the server if it was properly initialized
func Close(ctx context.Context) error {
	if instance != nil {
		return instance.Shutdown(ctx)
	}
	return nil
}

// serve starts the server on a Unix domain socket and sends runtime errors to errOut.
// Returns the socket path as a string, an io.Closer (the HTTP server) that can be
// used for immediate shutdown, and an error if any issue occurs during initialization.
func serve(errOut chan<- error) (string, *http.Server, error) {
	sp := getSocketPath()
	os.Remove(sp)

	listener, err := net.Listen("unix", sp)
	if err != nil {
		return "", nil, fmt.Errorf("failed to listen on UDS socket: %w", err)
	}

	h, err := newProxyHandler()
	if err != nil {
		return "", nil, err
	}
	server := &http.Server{
		Handler: internal.Handler(h),
	}

	go func() {
		if err := server.Serve(listener); err != http.ErrServerClosed {
			if errOut != nil {
				errOut <- err
			}
		}
		listener.Close()
		os.Remove(sp)
	}()

	return sp, server, nil
}

func getSocketPath() string {
	socketName := fmt.Sprintf("phishell.%d", os.Getpid())
	return filepath.Join(os.TempDir(), socketName)
}
