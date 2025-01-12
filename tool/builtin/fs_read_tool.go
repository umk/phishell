package builtin

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/tool"
	"github.com/umk/phishell/util/errorsx"
	"github.com/umk/phishell/util/fsx"
	"github.com/umk/phishell/util/marshalx"
	"github.com/umk/phishell/util/termx"
)

//go:embed schemas/fs_read.json
var fsReadFunctionBytes []byte

var fsReadFunction = tool.MustUnmarshalFn(fsReadFunctionBytes)

var FsReadToolName = fsReadFunction.Name

var FsReadTool = openai.ChatCompletionToolParam{
	Type:     openai.F(openai.ChatCompletionToolTypeFunction),
	Function: openai.Raw[openai.FunctionDefinitionParam](fsReadFunction),
}

type FsReadArguments struct {
	Path string `json:"path"`
}

type FsReadToolHandler struct {
	arguments *FsReadArguments
	path      string
}

func NewFsReadToolHandler(argsJSON, baseDir string) (*FsReadToolHandler, error) {
	var arguments FsReadArguments
	err := marshalx.UnmarshalJSONStruct([]byte(argsJSON), &arguments)
	if err != nil {
		return nil, errorsx.NewRetryableError(fmt.Sprintf("invalid arguments: %v", err))
	}

	return &FsReadToolHandler{
		arguments: &arguments,
		path:      fsx.Resolve(baseDir, arguments.Path),
	}, nil
}

func (h *FsReadToolHandler) Execute(ctx context.Context) (any, error) {
	if bootstrap.IsDebug(ctx) {
		termx.Muted.Printf("(call) %s; path=%s\n", FsReadToolName, h.arguments.Path)
	}

	f, err := os.ReadFile(h.path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return string(f), nil
}
