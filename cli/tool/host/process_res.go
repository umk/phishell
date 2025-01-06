package host

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"

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
		if err := p.processMessage(scanner.Bytes()); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read mesages: %w", err)
	}

	return nil
}

func (p *ToolProcess) processMessage(message []byte) error {
	var raw any
	if err := json.Unmarshal(message, &raw); err != nil {
		// A syntax error is a protocol violation that may result in an undefined
		// behavior, so the error is returned with further termination of the process.
		return err
	}

	var res provider.ToolResponse
	if err := marshalx.UnmarshalJSONStruct(message, &res); err != nil {
		// Assuming it's not a syntax error. Ignoring.
	} else {
		p.resolveRequest(&res)
	}

	return nil
}
