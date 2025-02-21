package builtin

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/tool"
	"github.com/umk/phishell/util/marshalx"
	"github.com/umk/phishell/util/termx"
)

//go:embed schemas/fs_delete.json
var fsDeleteFunctionBytes []byte

var fsDeleteFunction = tool.MustUnmarshalFn(fsDeleteFunctionBytes)

var FsDeleteToolName = fsDeleteFunction.Name

var FsDeleteTool = openai.ChatCompletionToolParam{
	Type:     openai.F(openai.ChatCompletionToolTypeFunction),
	Function: openai.Raw[openai.FunctionDefinitionParam](fsDeleteFunction),
}

type FsDeleteArguments struct {
	Path      string `json:"path" validate:"required"`
	Recursive bool   `json:"recursive"`
}

type FsDeleteToolHandler struct {
	arguments *FsDeleteArguments
	path      string
}

func NewFsDeleteToolHandler(argsJSON, baseDir string) (*FsDeleteToolHandler, error) {
	var arguments FsDeleteArguments
	err := marshalx.UnmarshalJSONStruct([]byte(argsJSON), &arguments)
	if err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	return &FsDeleteToolHandler{
		arguments: &arguments,
		path:      filepath.Join(baseDir, arguments.Path),
	}, nil
}

func (h *FsDeleteToolHandler) Execute(ctx context.Context) (any, error) {
	if bootstrap.IsDebug(ctx) {
		termx.Muted.Printf("(call) %s; path=%s\n", FsDeleteToolName, h.arguments.Path)
	}

	var err error

	if h.arguments.Recursive {
		err = os.RemoveAll(h.path)
	} else {
		err = os.Remove(h.path)
	}

	if err != nil {
		return nil, fmt.Errorf("error deleting file: %w", err)
	}

	return nil, nil
}

func (h *FsDeleteToolHandler) Describe(ctx context.Context) (string, error) {
	return fmt.Sprintf("Deleting file: %s", h.arguments.Path), nil
}
