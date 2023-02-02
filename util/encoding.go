package util

import (
	"encoding/hex"
	"encoding/json"
)

// JSONEncode encodes data to JSON.
func JSONEncode(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}

	return data
}

// JSONDecode decodes data from JSON.
func JSONDecode(data []byte, v any) any {
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil
	}

	return v
}

// HexDecode decodes data from hex.
func HexDecode(data string) []byte {
	b, err := hex.DecodeString(data)
	if err != nil {
		return nil
	}

	return b
}

// HexEncode encodes data to hex.
func HexEncode(data []byte) string {
	return hex.EncodeToString(data)
}
