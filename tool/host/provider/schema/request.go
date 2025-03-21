package schema

type Request struct {
	// Correlation ID of the tool call.
	CallID string `json:"call_id" validate:"required"`
	// The function that the model called.
	Function Function `json:"function" validate:"required"`
	// Context of the tool call request.
	Context Context `json:"context" validate:"required"`
}

type Function struct {
	// Name of the function to call.
	Name string `json:"name" validate:"required"`
	// Arguments to call the function with, as generated by the model in JSON format.
	Arguments string `json:"arguments" validate:"required"`
}

type Context struct {
	// Full path to the current directory.
	Dir string `json:"dir" validate:"required"`
}
