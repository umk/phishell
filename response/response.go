package response

import (
	"encoding/json"
	"strings"

	"github.com/umk/phishell/client"
)

const compactToSumRatio = 3

func mustSummarizeResp(cr *client.Ref, resp string) bool {
	toks := float32(len(resp)) / cr.Samples.BytesPerTok()

	return toks > float32(cr.Config.ContextSize)/compactToSumRatio
}

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
