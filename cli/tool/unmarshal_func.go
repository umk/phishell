package tool

import (
	"fmt"

	"github.com/umk/phishell/provider"
	"github.com/umk/phishell/util/marshalx"
)

// UnmarshalFn unmarshals a JSON-encoded tool function and returns a ToolFunction struct.
func UnmarshalFn(data []byte) (provider.ToolFunction, error) {
	var tool provider.ToolFunction
	if err := marshalx.UnmarshalJSONStruct(data, &tool); err != nil {
		return provider.ToolFunction{}, fmt.Errorf("failed to unmarshal tool function: %w", err)
	}
	return tool, nil
}

// MustUnmarshalFn unmarshals a JSON-encoded tool function and panics on error.
func MustUnmarshalFn(data []byte) provider.ToolFunction {
	tool, err := UnmarshalFn(data)
	if err != nil {
		panic(err)
	}
	return tool
}
