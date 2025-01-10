package response

import (
	"context"

	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/cli/msg"
	"github.com/umk/phishell/cli/prompt"
)

type HttpOutputParams struct {
	Url     string
	Status  string
	Headers map[string][]string
	Body    string
}

func GetHttpOutput(ctx context.Context, cr *bootstrap.ClientRef, params *HttpOutputParams) (string, error) {
	body := params.Body
	if s, ok := getJSONObjectOrArray(body); ok {
		body = s
	}

	if !mustSummarizeResp(cr, body) {
		return msg.FormatHttpResponseMessage(&msg.HttpResponseMessageParams{
			Status:  params.Status,
			Headers: params.Headers,
			Body:    params.Body,
			Summary: false,
		})
	}

	summaryCompl, err := prompt.PromptHttpSummary(ctx, &prompt.HttpSummaryPromptParams{
		Client:  cr,
		Url:     params.Url,
		Status:  params.Status,
		Headers: params.Headers,
		Body:    params.Body,
	})
	if err != nil {
		return "", err
	}

	summary, err := summaryCompl.Content()
	if err != nil {
		return "", err
	}

	return msg.FormatHttpResponseMessage(&msg.HttpResponseMessageParams{
		Status:  params.Status,
		Headers: params.Headers,
		Body:    summary,
		Summary: true,
	})
}
