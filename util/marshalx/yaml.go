package marshalx

import (
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

var unmarshalYAMLVal = validator.New(validator.WithRequiredStructEnabled())

// UnmarshalYAMLStruct takes a byte array containing YAML string and
// unmarshals it into the provided value pointer and validates
// the value using the validator package.
func UnmarshalYAMLStruct[A any](data []byte, out *A) error {
	if err := yaml.Unmarshal(data, out); err != nil {
		return err
	}

	if err := unmarshalYAMLVal.Struct(out); err != nil {
		return err
	}

	return nil
}
