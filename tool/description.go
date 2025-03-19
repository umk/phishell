package tool

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/umk/phishell/util/stringsx"
)

func DescribeCall(functionName, functionArgs, functionDescr string) string {
	nameToks := stringsx.Tokens(functionName)
	displayName := strings.ToLower(stringsx.DisplayName(nameToks))
	argsDescr, _ := describeArgs(functionArgs)

	var sb strings.Builder

	sb.WriteString("Running ")
	sb.WriteString(displayName)
	sb.WriteString("\n\n")

	if functionDescr != "" {
		sb.WriteString("Description: ")
		sb.WriteString(functionDescr)
		sb.WriteString("\n\n")
	}

	if functionArgs != "" {
		sb.WriteString(argsDescr)
		sb.WriteString("\n\n")
	}

	return sb.String()
}

type argumentsField struct {
	name        string
	value       string
	isMultiline bool
}

func describeArgs(argumentsJSON string) (string, error) {
	var arguments any
	if err := json.Unmarshal([]byte(argumentsJSON), &arguments); err != nil {
		return "", err
	}

	fields := traverseFields([]stringsx.Token{}, arguments)

	var sb strings.Builder

	// First pass: single-line items
	hasSingleLineFields := false
	for _, f := range fields {
		if !f.isMultiline {
			if !hasSingleLineFields {
				hasSingleLineFields = true
			}
			sb.WriteString(" - `")
			sb.WriteString(f.name)
			sb.WriteString("` ")
			sb.WriteString(f.value)
			sb.WriteString("\n")
		}
	}

	// Insert a blank line if we printed any single-line items
	if hasSingleLineFields {
		sb.WriteString("\n")
	}

	// Second pass: multi-line items
	for _, f := range fields {
		if f.isMultiline {
			sb.WriteString("### ")
			sb.WriteString(f.name)
			sb.WriteString("\n`````\n")
			sb.WriteString(f.value)
			if !strings.HasSuffix(f.value, "\n") {
				sb.WriteString("\n")
			}
			sb.WriteString("`````")
		}
	}

	return sb.String(), nil
}

func traverseFields(path []stringsx.Token, data any) []argumentsField {
	var fields []argumentsField

	switch v := data.(type) {
	case map[string]any:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			toks := stringsx.Tokens(k)
			fields = append(fields, traverseFields(append(path, toks...), v[k])...)
		}
	case string:
		isMultiline := strings.Contains(v, "\n")
		fields = append(fields, argumentsField{
			name:        stringsx.DisplayName(path),
			value:       v,
			isMultiline: isMultiline,
		})
	case bool:
		fields = append(fields, argumentsField{
			name:  stringsx.DisplayName(path),
			value: fmt.Sprintf("%t", v),
		})
	case int, int8, int16, int32, int64:
		fields = append(fields, argumentsField{
			name:  stringsx.DisplayName(path),
			value: fmt.Sprintf("%d", v),
		})
	case float32, float64:
		fields = append(fields, argumentsField{
			name:  stringsx.DisplayName(path),
			value: fmt.Sprintf("%f", v),
		})
	case []any:
		for i, val := range v {
			tok := stringsx.Token(fmt.Sprintf("[%d]", i))
			fields = append(fields, traverseFields(append(path, tok), val)...)
		}
	}
	return fields
}
