package util

import (
	"encoding/json"
)

// MarshalType alias for json.Marshal.
func MarshalType(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}

	return data
}

// UnmarshalType alias for json.Unmarshal.
func UnmarshalType(data []byte, v any) any {
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil
	}

	return v
}
