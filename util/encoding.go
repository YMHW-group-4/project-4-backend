package util

import (
	"encoding/json"
	"github.com/mr-tron/base58"
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

func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

func Base58Decode(input []byte) []byte {
	decode, _ := base58.Decode(string(input[:]))
	return decode
}
