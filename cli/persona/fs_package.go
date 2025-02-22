package persona

import (
	_ "embed"
	"strings"
	"text/template"
)

type PackageJSONParams struct{}

//go:embed templates/package.tmpl
var packageJSON string

var packageJSONTmpl = template.Must(template.New("package").Parse(packageJSON))

func formatPackageJSON(params PackageJSONParams) (string, error) {
	var sb strings.Builder
	if err := packageJSONTmpl.Execute(&sb, &params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
