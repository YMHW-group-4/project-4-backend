package util

import (
	"os"
	"strconv"
)

// GetEnv gets a specified type T value from an environment variable, when value
// cannot be retrieved the specified fallback is given instead.
func GetEnv[T any](key string, fallback T) T {
	var value any

	if val, ok := os.LookupEnv(key); ok && len(val) > 0 {
		switch any(fallback).(type) {
		case string:
			value = val
		case int:
			if v, err := strconv.Atoi(val); err == nil {
				value = v
			}
		case bool:
			if v, err := strconv.ParseBool(val); err == nil {
				value = v
			}
		}
	}

	if v, ok := value.(T); ok {
		return v
	}

	return fallback
}
