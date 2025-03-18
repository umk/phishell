package marshalx

import (
	"encoding/json"
	"maps"
)

type Dynamic[V any] struct {
	V          V
	properties map[string]any
}

func (d Dynamic[V]) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(d.V)
	if err != nil {
		return nil, err
	}
	var properties map[string]any
	if err := json.Unmarshal(b, &properties); err != nil {
		return nil, err
	}
	merged := make(map[string]any)
	if d.properties != nil {
		maps.Copy(merged, d.properties)
	}
	maps.Copy(merged, properties)
	return json.Marshal(merged)
}

func (d *Dynamic[V]) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &d.V); err != nil {
		return err
	}
	var properties map[string]any
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	d.properties = properties
	return nil
}
