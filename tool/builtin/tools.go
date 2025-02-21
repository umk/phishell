package builtin

import "github.com/openai/openai-go"

var Tools = []openai.ChatCompletionToolParam{
	ExecCommandTool,
	FsCreateOrUpdateTool,
	FsDeleteTool,
	FsReadTool,
}
