package models

import (
	"encoding/json"
	"fmt"
)

func scanJSON(ptr, val interface{}) error {
	v, ok := val.([]byte)
	if !ok {
		return fmt.Errorf("unsupported type: %T", v)
	}
	return json.Unmarshal(v, &ptr)
}
