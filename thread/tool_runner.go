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

	name    string
	handler tool.Handler
	err     error // An error occurred during handler creation
}

func NewToolRunner(host *host.Host) *ToolRunner {
	return &ToolRunner{host: host}
}

func (c *ToolRunner) Add(toolCall *openai.ChatCompletionMessageToolCall) error {
	h := c.getToolHandler(toolCall)

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

func (c *ToolRunner) getToolHandler(toolCall *openai.ChatCompletionMessageToolCall) *ToolRunnerHandler {
	id := toolCall.ID
	name := toolCall.Function.Name

	h, err := c.host.Get(&toolCall.Function)
	if err != nil {
		return &ToolRunnerHandler{id: id, name: name, err: err}
	}

	if h == nil {
		err := fmt.Errorf("function doesn't exist: %s", toolCall.Function.Name)
		return &ToolRunnerHandler{id: id, name: name, err: err}
	}

	return &ToolRunnerHandler{id: id, name: name, handler: h}
}

func (c *ToolRunner) processError(err error) string {
	return fmt.Sprintf("Error calling function: %v", err)
}

func confirmToolCall(ctx context.Context, h *ToolRunnerHandler) error {
	t := stringsx.Tokens(h.name)

	if len(t) == 0 {
		return nil
	}

	descr, _ := tool.Describe(ctx, h.handler)

	if descr != "" {
		termx.NewPrinter().Println(descr)
	} else {
		dn := strings.ToLower(stringsx.DisplayName(t))
		termx.NewPrinter().Printf("Running %s\n", dn)
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
