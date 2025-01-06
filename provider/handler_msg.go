package provider

import (
	"encoding/json"
	"fmt"
)

// Post prints a message outside of the host to provider interactions. The message
// is processed on demand by the end user.
func Post(m Message) error {
	if err := providerVal.Struct(&m); err != nil {
		return err
	}

	b, err := json.Marshal(&m)
	if err != nil {
		return err
	}

	if _, err := fmt.Println(string(b)); err != nil {
		return err
	}

	return nil
}
