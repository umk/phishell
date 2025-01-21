package process

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/umk/phishell/provider"
	"github.com/umk/phishell/util/marshalx"
)

func (p *Process) Read() error {
	if p.scanner == nil {
		return errors.New("process is not initialized")
	}

	for p.scanner.Scan() {
		if err := p.processMessage(p.scanner.Bytes()); err != nil {
			return err
		}
	}

	if err := p.scanner.Err(); err != nil {
		return fmt.Errorf("failed to read mesages: %w", err)
	}

	return nil
}

func (p *Process) processMessage(message []byte) error {
	var res provider.ToolResponse
	if err := marshalx.UnmarshalJSONStruct(message, &res); err != nil {
		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			// A syntax error is a protocol violation that may result in an undefined
			// behavior, so the error is returned with further termination of the process.
			return err
		}
		// Ignore error.
	} else {
		p.requestResolve(&res)
	}

	return nil
}
