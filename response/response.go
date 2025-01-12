package response

import (
	"encoding/json"
	"strings"

	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/prompt/client"
)

const compactToSumRatio = 3

func mustSummarizeResp(cr *bootstrap.ClientRef, resp string) bool {
	cl := client.Get(cr)
	toks := float32(len(resp)) / cl.Samples.BytesPerTok()

	return toks > float32(cr.Config.CompactionToks)/compactToSumRatio
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
