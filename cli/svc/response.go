package svc

import (
	"encoding/json"
	"strings"
)

func getJSONObjectOrArray(s string) (string, bool) {
	if (strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]")) ||
		(strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")) {
		var outputObj any
		if err := json.Unmarshal([]byte(s), &outputObj); err == nil {
			data, err := json.Marshal(outputObj)
			if err == nil {
				s := string(data)
				return s, true
			}
		}
	}

	return "", false
}
