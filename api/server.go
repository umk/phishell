package api

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Serve starts the server on a Unix domain socket and sends runtime errors to errOut.
// Returns the socket path as a string, an io.Closer (the HTTP server) that can be
// used for immediate shutdown, and an error if any issue occurs during initialization.
func Serve(ctx context.Context, errOut chan<- error) (string, *http.Server, error) {
	sp := getSocketPath()
	os.Remove(sp)

	listener, err := net.Listen("unix", sp)
	if err != nil {
		return "", nil, fmt.Errorf("failed to listen on UDS socket: %w", err)
	}

	var si handler
	server := &http.Server{
		Handler: Handler(&si),
	}

	go func() {
		if err := server.Serve(listener); err != http.ErrServerClosed {
			errOut <- err
		}
		listener.Close()
		os.Remove(sp)
	}()

	return sp, server, nil
}

func getSocketPath() string {
	socketName := fmt.Sprintf("phishell.%d%d", time.Now().Unix()%1e3, rand.Intn(1e5))
	return filepath.Join(os.TempDir(), socketName)
}
