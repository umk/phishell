package builtin

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/cli/msg"
	"github.com/umk/phishell/cli/tool"
	"github.com/umk/phishell/util/errorsx"
	"github.com/umk/phishell/util/marshalx"
)

//go:embed schemas/exec_http_call.json
var execHttpCallFunctionBytes []byte

var execHttpCallFunction = tool.MustUnmarshalFn(execHttpCallFunctionBytes)

var ExecHttpCallToolName = execHttpCallFunction.Name

var ExecHttpCallTool = openai.ChatCompletionToolParam{
	Type:     openai.F(openai.ChatCompletionToolTypeFunction),
	Function: openai.Raw[openai.FunctionDefinitionParam](execHttpCallFunction),
}

type ExecHttpCallArguments struct {
	Url             string            `json:"url" validate:"required"`
	Method          string            `json:"method" validate:"required,oneof=GET POST PUT DELETE PATCH HEAD OPTIONS"`
	Headers         map[string]string `json:"headers"`
	QueryParameters map[string]string `json:"query_parameters"`
	Body            *string           `json:"body"`
}

type ExecHttpCallToolHandler struct {
	arguments *ExecHttpCallArguments
	url       *url.URL
}

func NewExecHttpCallToolHandler(argsJSON string) (*ExecHttpCallToolHandler, error) {
	// Parsing arguments
	var arguments ExecHttpCallArguments
	err := marshalx.UnmarshalJSONStruct([]byte(argsJSON), &arguments)
	if err != nil {
		return nil, errorsx.NewRetryableError(fmt.Sprintf("invalid arguments: %v", err))
	}

	// Parse URL
	u, err := url.Parse(arguments.Url)
	if err != nil {
		return nil, errorsx.NewRetryableError(fmt.Sprintf("invalid arguments: %v", err))
	}

	query := u.Query()
	for key, value := range arguments.QueryParameters {
		query.Set(key, value)
	}
	u.RawQuery = query.Encode()

	return &ExecHttpCallToolHandler{
		arguments: &arguments,
		url:       u,
	}, nil
}

func (h *ExecHttpCallToolHandler) Execute(ctx context.Context) (any, error) {
	// Create the HTTP request
	var body io.Reader
	if h.arguments.Body != nil {
		body = strings.NewReader(*h.arguments.Body)
	}

	req, err := http.NewRequestWithContext(ctx, h.arguments.Method, h.url.String(), body)
	if err != nil {
		return nil, err
	}

	// Set headers
	version := bootstrap.GetVersion(ctx)
	req.Header.Set("User-Agent", fmt.Sprintf("PhiShell/%s", version))

	for key, value := range h.arguments.Headers {
		req.Header.Set(key, value)
	}

	// Execute the HTTP request using the default HTTP client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Println(resp.Status)

	var rb string

	if resp.Body != nil {
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		rb = string(b)

		if rb != "" {
			fmt.Println(rb)
		}
	}

	r, err := msg.FormatHttpResponseMessage(&msg.HttpResponseMessageParams{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       rb,
	})
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (h *ExecHttpCallToolHandler) Describe(ctx context.Context) (string, error) {
	return fmt.Sprintf("%s %s", h.arguments.Method, h.url), nil
}