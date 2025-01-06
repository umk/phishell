package provider

type Tool struct {
	// The type of the tool. Currently, only `function` is supported.
	Type string `json:"type" validate:"required,oneof=function"`
	// Description of the function.
	Function *ToolFunction `json:"function,omitempty"`
}

type ToolFunction struct {
	// The name of the function to be called. Must be a-z, A-Z, 0-9, or
	// contain underscores and dashes, with a maximum length of 64.
	Name string `json:"name" validate:"required"`
	// A description of what the function does, used by the model to
	// choose when and how to call the function.
	Description string `json:"description,omitempty"`
	// The parameters the functions accepts, described as a JSON Schema object.
	Parameters map[string]any `json:"parameters" validate:"required"`
}
