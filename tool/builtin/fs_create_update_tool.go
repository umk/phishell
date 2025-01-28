package builtin

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/openai/openai-go"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/tool"
	"github.com/umk/phishell/util/errorsx"
	"github.com/umk/phishell/util/fsx"
	"github.com/umk/phishell/util/marshalx"
	"github.com/umk/phishell/util/stringsx"
	"github.com/umk/phishell/util/termx"
)

//go:embed schemas/fs_create_update.json
var fsCreateOrUpdateFunctionBytes []byte

var fsCreateOrUpdateFunction = tool.MustUnmarshalFn(fsCreateOrUpdateFunctionBytes)

var FsCreateOrUpdateToolName = fsCreateOrUpdateFunction.Name

var FsCreateOrUpdateTool = openai.ChatCompletionToolParam{
	Type:     openai.F(openai.ChatCompletionToolTypeFunction),
	Function: openai.Raw[openai.FunctionDefinitionParam](fsCreateOrUpdateFunction),
}

type FsCreateOrUpdateArguments struct {
	Path        string `json:"path" validate:"required"`
	FileContent string `json:"file_content"`
}

type FsCreateOrUpdateToolHandler struct {
	arguments *FsCreateOrUpdateArguments
	path      string
}

func NewFsCreateOrUpdateToolHandler(argsJSON, baseDir string) (*FsCreateOrUpdateToolHandler, error) {
	var arguments FsCreateOrUpdateArguments
	err := marshalx.UnmarshalJSONStruct([]byte(argsJSON), &arguments)
	if err != nil {
		return nil, errorsx.NewRetryableError(fmt.Sprintf("invalid arguments: %v", err))
	}

	return &FsCreateOrUpdateToolHandler{
		arguments: &arguments,
		path:      fsx.Resolve(baseDir, arguments.Path),
	}, nil
}

func (h *FsCreateOrUpdateToolHandler) Execute(ctx context.Context) (any, error) {
	if bootstrap.IsDebug(ctx) {
		termx.Muted.Printf("(call) %s; path=%s\n", FsCreateOrUpdateToolName, h.arguments.Path)
	}

	d := filepath.Dir(h.path)

	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}

	content := []byte(h.arguments.FileContent)
	if err := os.WriteFile(h.path, content, 0644); err != nil {
		return nil, fmt.Errorf("error writing file: %w", err)
	}

	return nil, nil
}

func (h *FsCreateOrUpdateToolHandler) Describe(ctx context.Context) (string, error) {
	c, err := os.ReadFile(h.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Sprintf("Creating file: %s", h.arguments.Path), nil
		}

		return "", err
	}

	s := string(c)

	d := diffmatchpatch.New()
	diffs := d.DiffMain(s, h.arguments.FileContent, false)
	t := stringsx.RenderDiff(diffs)

	return fmt.Sprintf("Updating file: %s\n%s", h.arguments.Path, t), nil
}
