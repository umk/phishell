package tool

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared/constant"
)

type Tool struct {
	// The type of the tool. Currently, only `function` is supported.
	Type string `json:"type" validate:"required,oneof=function"`
	// Description of the function.
	Function *Function `json:"function,omitempty"`
}

func (t *Tool) ToChatCompletionToolParam() openai.ChatCompletionToolParam {
	return openai.ChatCompletionToolParam{
		Type:     constant.Function(t.Type),
		Function: t.Function.ToFunctionDefinitionParam(),
	}
}

type Function struct {
	// The name of the function to be called. Must be a-z, A-Z, 0-9, or
	// contain underscores and dashes, with a maximum length of 64.
	Name string `json:"name" validate:"required"`
	// A description of what the function does, used by the model to
	// choose when and how to call the function.
	Description string `json:"description,omitempty"`
	// The parameters the functions accepts, described as a JSON Schema object.
	Parameters map[string]any `json:"parameters" validate:"required"`
}

func (f *Function) ToFunctionDefinitionParam() openai.FunctionDefinitionParam {
	r := openai.FunctionDefinitionParam{
		Name:       f.Name,
		Parameters: f.Parameters,
	}

	if f.Description != "" {
		r.Description = openai.String(f.Description)
	}

	return r
}
