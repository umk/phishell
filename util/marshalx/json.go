package marshalx

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

var unmarshalJSONVal = validator.New(validator.WithRequiredStructEnabled())

// UnmarshalJSONStruct takes a byte array containing JSON string and
// unmarshals it into the provided value pointer and validates
// the value using the validator package.
func UnmarshalJSONStruct[A any](data []byte, out *A) error {
	if err := json.Unmarshal(data, out); err != nil {
		return err
	}

	if err := unmarshalJSONVal.Struct(out); err != nil {
		return err
	}

	return nil
}
