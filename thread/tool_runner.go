package thread

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/tool"
	"github.com/umk/phishell/tool/host"
	"github.com/umk/phishell/util/errorsx"
	"github.com/umk/phishell/util/stringsx"
	"github.com/umk/phishell/util/termx"
)

type ToolRunner struct {
	host *host.Host

	handlers []*ToolRunnerHandler
}

type ToolRunnerHandler struct {
	id string

	function      openai.ChatCompletionMessageToolCallFunction
	functionDescr string
	handler       tool.Handler
	err           error // An error occurred during handler creation
}

func NewToolRunner(host *host.Host) *ToolRunner {
	return &ToolRunner{host: host}
}

func (c *ToolRunner) Add(toolCall *openai.ChatCompletionMessageToolCall, functionDescr string) error {
	h := c.getToolHandler(toolCall, functionDescr)

	c.handlers = append(c.handlers, h)

	return h.err
}

func (c *ToolRunner) Complete(ctx context.Context) ([]openai.ChatCompletionMessageParamUnion, error) {
	messages := make([]openai.ChatCompletionMessageParamUnion, len(c.handlers))

	hasErrors := false

	for i, h := range c.handlers {
		var handlerErr error

		if !hasErrors {
			if h.err != nil {
				handlerErr = h.err
			} else if content, err := c.callTool(ctx, h); err != nil {
				handlerErr = err
			} else {
				messages[i] = openai.ToolMessage(h.id, content)
			}
		} else {
			handlerErr = errors.New("operation was canceled because the previous operation could not be completed")
		}

		if handlerErr != nil {
			if errorsx.IsCanceled(handlerErr) {
				return nil, handlerErr
			}

			hasErrors = true

			messages[i] = openai.ToolMessage(h.id, c.processError(handlerErr))
		}
	}

	return messages, nil
}

func (c *ToolRunner) callTool(ctx context.Context, h *ToolRunnerHandler) (string, error) {
	if err := confirmToolCall(ctx, h); err != nil {
		return "", err
	}

	res, err := h.handler.Execute(ctx)
	if err != nil {
		return "", err
	}

	content, err := getToolResponseContent(res)
	if err != nil {
		return "", err
	}

	return content, nil
}

func (c *ToolRunner) getToolHandler(
	toolCall *openai.ChatCompletionMessageToolCall, functionDescr string,
) *ToolRunnerHandler {
	id := toolCall.ID

	h, err := c.host.Handler(&toolCall.Function)
	if err != nil {
		return &ToolRunnerHandler{id: id, function: toolCall.Function, err: err}
	}

	if h == nil {
		err := fmt.Errorf("function doesn't exist: %s", toolCall.Function.Name)
		return &ToolRunnerHandler{id: id, function: toolCall.Function, err: err}
	}

	return &ToolRunnerHandler{
		id:            id,
		function:      toolCall.Function,
		functionDescr: functionDescr,
		handler:       h,
	}
}

func (c *ToolRunner) processError(err error) string {
	return fmt.Sprintf("Error calling function: %v", err)
}

func confirmToolCall(ctx context.Context, h *ToolRunnerHandler) error {
	t := stringsx.Tokens(h.function.Name)

	if len(t) == 0 {
		return nil
	}

	if descr, _ := tool.Describe(ctx, h.handler); descr != "" {
		termx.MD.Println(descr)
	} else {
		termx.MD.Println(tool.DescribeCall(h.function.Name, h.function.Arguments, h.functionDescr))
	}

	for {
		s, err := termx.ReadPrompt(ctx, &termx.Static{
			Prompt: ">>> ",
			Hint:   "Press Enter to execute or Ctrl+C to cancel",
		})
		if err != nil {
			return err
		}

		if s == "" {
			return nil
		}

		s = strings.TrimSpace(s)
		if s != "" {
			return errors.New(s)
		}
	}
}
