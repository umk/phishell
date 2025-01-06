package provider

import "context"

type contextKey string

const (
	ctxDir contextKey = "dir"
)

// Dir gets the working directory from the context.
func Dir(ctx context.Context) string {
	return getRequestValue(ctx, ctxDir).(string)
}

func getRequestValue(ctx context.Context, key contextKey) any {
	v := ctx.Value(key)
	if v == nil {
		panic("must be called in a context of a tool request")
	}

	return v
}
