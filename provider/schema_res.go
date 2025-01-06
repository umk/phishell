package provider

type ToolResponse struct {
	// ID of the tool call that this message is responding to.
	CallID string `json:"call_id" validate:"required"`
	// The contents of the tool response.
	Content any `json:"content"`
}
