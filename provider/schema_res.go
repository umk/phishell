package provider

import "time"

type ToolResponse struct {
	// ID of the tool call that this message is responding to.
	CallID string `json:"call_id" validate:"required"`
	// The contents of the tool response.
	Content any `json:"content"`
}

// Message represents a message initiated by the provider itself.
type Message struct {
	// Unique ID of the message used for deduplication.
	ID string `json:"id,omitempty"`
	// Content of the message.
	Content string `json:"content" validate:"required"`
	// Working directory to use in connection with the message.
	Dir string `json:"dir,omitempty"`
	// Date and time when the message was created by provider. Can be
	// different from the date and time when the message was actually sent to
	// the host.
	Date *time.Time `json:"date,omitempty"`
}
