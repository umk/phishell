package builtin

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/cli/msg"
	"github.com/umk/phishell/cli/tool"
	"github.com/umk/phishell/util/errorsx"
	"github.com/umk/phishell/util/fsx"
	"github.com/umk/phishell/util/marshalx"
	"github.com/umk/phishell/util/termx"
)

//go:embed schemas/fs_stat.json
var fsStatFunctionBytes []byte

var fsStatFunction = tool.MustUnmarshalFn(fsStatFunctionBytes)

var FsStatToolName = fsStatFunction.Name

var FsStatTool = openai.ChatCompletionToolParam{
	Type:     openai.F(openai.ChatCompletionToolTypeFunction),
	Function: openai.Raw[openai.FunctionDefinitionParam](fsStatFunction),
}

type FsStatArguments struct {
	Path string `json:"path"`
}

type FsStatToolHandler struct {
	arguments *FsStatArguments
	path      string
}

func NewFsStatToolHandler(argsJSON, baseDir string) (*FsStatToolHandler, error) {
	var arguments FsStatArguments
	err := marshalx.UnmarshalJSONStruct([]byte(argsJSON), &arguments)
	if err != nil {
		return nil, errorsx.NewRetryableError(fmt.Sprintf("invalid arguments: %v", err))
	}

	return &FsStatToolHandler{
		arguments: &arguments,
		path:      fsx.Resolve(baseDir, arguments.Path),
	}, nil
}

func (h *FsStatToolHandler) Execute(ctx context.Context) (any, error) {
	if bootstrap.IsDebug(ctx) {
		termx.Muted.Printf("(call) %s; path=%s\n", FsStatToolName, h.arguments.Path)
	}

	s, err := os.Stat(h.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "File or directory doesn't exist.", nil
		}

		return nil, err
	}

	return msg.FormatStatMessage(&msg.StatMessageParams{
		IsDirectory: s.IsDir(),
		Size:        s.Size(),
	})
}
