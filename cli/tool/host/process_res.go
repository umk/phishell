package host

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/umk/phishell/cli/tool"
	"github.com/umk/phishell/provider"
	"github.com/umk/phishell/util/marshalx"
)

func (p *ToolProcess) readHeader(scanner *bufio.Scanner) error {
	for scanner.Scan() {
		t := scanner.Bytes()

		if len(t) == 0 {
			return nil
		}

		if err := p.readHeaderLine(t); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	return errors.New("process ended with incomplete header")
}

func (p *ToolProcess) readHeaderLine(b []byte) error {
	tool, err := tool.UnmarshalTool(b)
	if err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	function := tool.Function

	if _, ok := p.host.tools[function.Name]; ok {
		return fmt.Errorf("function duplicate: %s", function.Name)
	}

	p.Tools[function.Name] = openai.ChatCompletionToolParam{
		Type:     openai.F(openai.ChatCompletionToolTypeFunction),
		Function: openai.Raw[openai.FunctionDefinitionParam](function),
	}

	return nil
}

func (p *ToolProcess) readMessages(ctx context.Context, scanner *bufio.Scanner) error {
	for scanner.Scan() {
		if err := p.processMessage(ctx, scanner.Bytes()); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read mesages: %w", err)
	}

	return nil
}

func (p *ToolProcess) processMessage(ctx context.Context, message []byte) error {
	var raw any
	if err := json.Unmarshal(message, &raw); err != nil {
		// A syntax error is a protocol violation that may result in an undefined
		// behavior, so the error is returned with further termination of the process.
		return err
	}

	properties, ok := raw.(map[string]any)
	if !ok {
		// Invalid message format: not an object. Ignoring.
		return nil
	}

	if _, ok := properties["call_id"]; ok {
		var res provider.ToolResponse
		if err := marshalx.UnmarshalJSONStruct(message, &res); err != nil {
			// Assuming it's not a syntax error. Ignoring.
		} else {
			p.resolveRequest(&res)
		}
	} else {
		var msg provider.Message
		if err := marshalx.UnmarshalJSONStruct(message, &msg); err != nil {
			// Assuming it's not a syntax error. Ignoring.
		} else {
			content := strings.TrimSpace(msg.Content)

			if content != "" {
				id := msg.ID
				content := strings.TrimSpace(msg.Content)
				wd := msg.Dir
				date := msg.Date

				if id == "" {
					id = uuid.NewString()
				}

				if date == nil || date.IsZero() {
					now := time.Now()
					date = &now
				}

				go p.host.events.ProcessEvent(ctx, id, content, wd, *date)
			}
		}
	}

	return nil
}
