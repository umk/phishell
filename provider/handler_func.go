package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type FunctionHandler[A any] func(ctx context.Context, arguments A) (any, error)

// Function registers a handler of a function tool.
func Function[A any](h FunctionHandler[A], f ToolFunction) error {
	if err := providerVal.Struct(&f); err != nil {
		return err
	}

	var value A
	if reflect.TypeOf(value).Kind() != reflect.Struct {
		return fmt.Errorf("type %T is not a structure", value)
	}

	ensureTools()

	tools[f.Name] = func(ctx context.Context, req *ToolRequest) error {
		return processFuncRequest(ctx, h, req)
	}

	b, err := json.Marshal(Tool{
		Type:     "function",
		Function: &f,
	})
	if err != nil {
		return err
	}

	if _, err := fmt.Println(string(b)); err != nil {
		return err
	}

	return nil
}

func processFuncRequest[A any](ctx context.Context, h FunctionHandler[A], req *ToolRequest) error {
	ctx = context.WithValue(ctx, ctxDir, req.Context.Dir)

	var content any

	var arguments A
	if err := json.Unmarshal([]byte(req.Function.Arguments), &arguments); err != nil {
		content = fmt.Sprintf("Bad arguments: %v", err)
	} else if err := providerVal.Struct(&arguments); err != nil {
		content = fmt.Sprintf("Bad arguments: %v", err)
	} else if v, err := h(ctx, arguments); err != nil {
		content = fmt.Sprintf("Error: %v", err)
	} else {
		content = v
	}

	b, err := json.Marshal(ToolResponse{
		CallID:  req.CallID,
		Content: content,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error converting to JSON: %v", err)
	}

	fmt.Println(string(b))

	return nil
}
