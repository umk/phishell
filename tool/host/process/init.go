package process

import (
	"bufio"
	"errors"
	"fmt"
	"time"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/tool"
)

func (p *Process) Init() error {
	p.scanner = bufio.NewScanner(p.stdout)

	init := make(chan error, 1)

	go func() {
		init <- p.readHeader(p.scanner)
		close(init)
	}()

	select {
	case err := <-init:
		if err != nil {
			return fmt.Errorf("failed to initialize tools: %w", err)
		}
	case <-time.After(10 * time.Second):
		return errors.New("initialization timeout")
	}

	return nil
}

func (p *Process) readHeader(scanner *bufio.Scanner) error {
	for scanner.Scan() {
		t := scanner.Bytes()

		if len(t) == 0 {
			return nil
		}

		if err := p.readHeaderLine(t); err != nil {
			return fmt.Errorf("failed to read header: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	return errors.New("process ended with incomplete header")
}

func (p *Process) readHeaderLine(b []byte) error {
	tool, err := tool.UnmarshalTool(b)
	if err != nil {
		return err
	}

	function := tool.Function

	p.tools[function.Name] = openai.ChatCompletionToolParam{
		Type:     openai.F(openai.ChatCompletionToolTypeFunction),
		Function: openai.Raw[openai.FunctionDefinitionParam](function),
	}

	return nil
}
