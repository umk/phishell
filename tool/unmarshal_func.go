package tool

import (
	"fmt"

	"github.com/umk/phishell/util/marshalx"
)

// UnmarshalFn unmarshals a JSON-encoded tool function and returns a ToolFunction struct.
func UnmarshalFn(data []byte) (Function, error) {
	var tool Function
	if err := marshalx.UnmarshalJSONStruct(data, &tool); err != nil {
		return Function{}, fmt.Errorf("failed to unmarshal tool function: %w", err)
	}
	return tool, nil
}

// MustUnmarshalFn unmarshals a JSON-encoded tool function and panics on error.
func MustUnmarshalFn(data []byte) Function {
	tool, err := UnmarshalFn(data)
	if err != nil {
		panic(err)
	}
	return tool
}
