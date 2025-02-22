package persona

import (
	_ "embed"
	"strings"
	"text/template"
)

type IndexTSParams struct{}

//go:embed templates/index.tmpl
var indexTS string

var indexTSTmpl = template.Must(template.New("index").Parse(indexTS))

func formatIndexTS(params IndexTSParams) (string, error) {
	var sb strings.Builder
	if err := indexTSTmpl.Execute(&sb, &params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
