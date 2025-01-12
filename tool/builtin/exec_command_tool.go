package builtin

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/response"
	"github.com/umk/phishell/tool"
	"github.com/umk/phishell/util/errorsx"
	"github.com/umk/phishell/util/execx"
	"github.com/umk/phishell/util/fsx"
	"github.com/umk/phishell/util/marshalx"
	"github.com/umk/phishell/util/termx"
)

//go:embed schemas/exec_command.json
var execCommandFunctionBytes []byte

var execCommandFunction = tool.MustUnmarshalFn(execCommandFunctionBytes)

var ExecCommandToolName = execCommandFunction.Name

var ExecCommandTool = openai.ChatCompletionToolParam{
	Type:     openai.F(openai.ChatCompletionToolTypeFunction),
	Function: openai.Raw[openai.FunctionDefinitionParam](execCommandFunction),
}

type ExecCommandArguments struct {
	CommandLine string `json:"command_line" validate:"required"`
	WorkingDir  string `json:"working_dir"`
}

type ExecCommandToolHandler struct {
	arguments  *ExecCommandArguments
	cmds       []*execx.Cmd
	workingDir string
}

func NewExecCommandToolHandler(argsJSON, baseDir string) (*ExecCommandToolHandler, error) {
	// Parsing arguments
	var arguments ExecCommandArguments
	err := marshalx.UnmarshalJSONStruct([]byte(argsJSON), &arguments)
	if err != nil {
		return nil, errorsx.NewRetryableError(fmt.Sprintf("invalid arguments: %v", err))
	}

	// Getting working directory
	wd := fsx.Resolve(baseDir, arguments.WorkingDir)

	// Parsing command parameters
	piped, err := execx.Parse(arguments.CommandLine)
	if err != nil {
		return nil, errorsx.NewRetryableError(fmt.Sprintf("invalid command line: %s", err))
	}

	cmds := make([]*execx.Cmd, len(piped))
	for i, p := range piped {
		cmd, err := execx.AllocArgs(p)
		if err != nil {
			return nil, err
		}

		cmds[i] = cmd
	}

	return &ExecCommandToolHandler{
		arguments:  &arguments,
		cmds:       cmds,
		workingDir: wd,
	}, nil
}

func (h *ExecCommandToolHandler) Execute(ctx context.Context) (any, error) {
	if bootstrap.IsDebug(ctx) {
		termx.Muted.Printf("(call) %s; command=%s\n", ExecCommandToolName, h.arguments.CommandLine)
	}

	// Building command pipeline
	cmds := make([]*exec.Cmd, len(h.cmds))
	for i, c := range h.cmds {
		cmd := c.CommandContext(ctx)
		cmd.Dir = h.workingDir

		cmds[i] = cmd
	}

	if err := execx.Pipe(cmds, nil, os.Stdout, os.Stdout); err != nil {
		return nil, err
	}

	// Running command
	logger := execx.Log(cmds[len(cmds)-1], 15000)

	exitCode, err := execx.RunPipe(cmds)
	if err != nil {
		return nil, err
	}

	// Getting process output
	processOut, err := logger.Output()
	if err != nil {
		return nil, err
	}

	cr := bootstrap.GetClient(ctx)
	output, err := response.GetExecOutput(ctx, cr, &response.ExecOutputParams{
		CommandLine: h.arguments.CommandLine,
		ExitCode:    exitCode,
		Output:      processOut,
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (h *ExecCommandToolHandler) Describe(ctx context.Context) (string, error) {
	return h.arguments.CommandLine, nil
}
