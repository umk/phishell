package provider

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/go-playground/validator/v10"
)

var providerVal = validator.New(validator.WithRequiredStructEnabled())

type requestHandler func(ctx context.Context, req *ToolRequest) error

var tools map[string]requestHandler

func ensureTools() {
	if tools == nil {
		tools = make(map[string]requestHandler)
	}
}

// Init indicates that all registrations have been completed. This is a mandatory call
// before processing of incoming messages can begin.
func Init() {
	// Print an empty line to indicate the end of the header block
	fmt.Println()
}

// Serve starts reading tool messages from Stdin and calling registered tools that can
// handle these messages.
func Serve(ctx context.Context) error {
	var wg sync.WaitGroup

	messagesCh := make(chan any)

	wg.Add(1)

	go func() {
		defer wg.Done()

		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			messagesCh <- s.Text()
		}

		if err := s.Err(); err != nil {
			messagesCh <- err
		}

		close(messagesCh)
	}()

	var err error

Messages:
	for message := range messagesCh {
		switch current := message.(type) {
		case string:
			wg.Add(1)

			go func() {
				defer wg.Done()

				if err := processMessage(ctx, current); err != nil {
					messagesCh <- err
				}
			}()
		case error:
			err = current
			break Messages
		default:
			panic("invalid message type")
		}
	}

	wg.Wait()

	return err
}

func processMessage(ctx context.Context, message string) error {
	var req ToolRequest
	if err := json.Unmarshal([]byte(message), &req); err != nil {
		return err
	}

	if err := providerVal.Struct(&req); err != nil {
		fmt.Fprintf(os.Stderr, "Bad request: %v", err)
		return nil
	}

	handler, ok := tools[req.Function.Name]
	if !ok {
		fmt.Fprintf(os.Stderr, "Tool is not supported: %s", req.Function.Name)
		return nil
	}

	return handler(ctx, &req)
}
