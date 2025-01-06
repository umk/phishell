package tool

import (
	"fmt"

	"github.com/umk/phishell/provider"
	"github.com/umk/phishell/util/marshalx"
)

// UnmarshalTool unmarshals a JSON-encoded tool and returns a Tool struct.
func UnmarshalTool(data []byte) (provider.Tool, error) {
	var tool provider.Tool
	if err := marshalx.UnmarshalJSONStruct(data, &tool); err != nil {
		return provider.Tool{}, fmt.Errorf("failed to unmarshal tool: %w", err)
	}

	switch tool.Type {
	case "function":
		if tool.Function == nil {
			return provider.Tool{}, fmt.Errorf("missing the function definition")
		}
	}

	return tool, nil
}

// MustUnmarshalTool unmarshals a JSON-encoded tool and panics on error.
func MustUnmarshalTool(data []byte) provider.Tool {
	tool, err := UnmarshalTool(data)
	if err != nil {
		panic(err)
	}
	return tool
}
