package marshalx

import (
	"gopkg.in/yaml.v3"
)

// UnmarshalYAMLStruct takes a byte array containing YAML string and
// unmarshals it into the provided value pointer and validates
// the value using the validator package.
func UnmarshalYAMLStruct[A any](data []byte, out *A) error {
	if err := yaml.Unmarshal(data, out); err != nil {
		return err
	}

	if err := Validator.Struct(out); err != nil {
		return err
	}

	return nil
}
