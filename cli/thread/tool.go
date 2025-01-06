package thread

import (
	"encoding/json"
	"fmt"
)

// getToolResponseContent gets the message content suitable for feeding to LLM
// as a response to the function call.
func getToolResponseContent(response any) (string, error) {
	if response == nil {
		return "OK", nil
	}

	switch response.(type) {
	case string, int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		complex64, complex128,
		bool:
		return fmt.Sprint(response), nil
	}

	j, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(j), nil
}
