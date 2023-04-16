package secretfields

import (
	"encoding/json"
	"fmt"
)

// Global settings

var RedactSecretsInJSON bool = true

// Types

type Secret[T any] struct {
	SecretValue   T
	redactedValue T
}

// Factories

func AsSecret[T any](value T, redactedValue ...T) Secret[T] {
	var redacted T
	if len(redactedValue) > 0 {
		redacted = redactedValue[0]
	}
	return Secret[T]{
		SecretValue:   value,
		redactedValue: redacted,
	}
}

// JSON serdes

func (s Secret[T]) MarshalJSON() ([]byte, error) {
	safeValue := s.redactedValue
	if !RedactSecretsInJSON {
		safeValue = s.SecretValue
	}
	return json.Marshal(safeValue)
}

func (s *Secret[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.SecretValue)
}

// Stringer interface

func (s Secret[T]) String() string {
	return fmt.Sprint(s.redactedValue)
}
