package gosecretfields

import (
	"encoding/json"
	"fmt"
)

// Types

// Container for values tagged as secrets
type Secret[T any] struct {
	// This is the secret value, it will be always redacted when
	// converted into a string through the provided stringer interface
	// for the container.
	// Wether or not it is redacted when serializing to JSON is determined
	// by `CleartextJSON` attribute.
	// It can be explicitly accesed through this public attribute.
	SecretValue T

	// This is the value that replaces the secret value upon redaction.
	redactedValue T

	// Flag controling wether this secret should be JSON serialized with
	// its secret or redacted value.
	CleartextJSON bool
}

// Factories

// Factory method to wrap any value around a `Secret` container,
// it, the result is a value tagged as secret that won't leak
// on logs or any other string conversion and that might or might
// not be redacted in JSON serializations
func AsSecret[T any](value T, redactedValue ...T) Secret[T] {
	var redacted T
	if len(redactedValue) > 0 {
		redacted = redactedValue[0]
	}
	return Secret[T]{
		SecretValue:   value,
		redactedValue: redacted,
		CleartextJSON: false,
	}
}

// JSON serdes

func (s Secret[T]) MarshalJSON() ([]byte, error) {
	// Depending on the value of `s.CleartextJSON` flag, JSON serialization of
	// Secret fields will result on the redactedValue or the actual secret value JSON representation
	// but the container will never show in the JSON structure.
	safeValue := s.redactedValue
	if s.CleartextJSON {
		safeValue = s.SecretValue
	}
	return json.Marshal(safeValue)
}

func (s *Secret[T]) UnmarshalJSON(data []byte) error {
	// The fact that a value is tagged as secret is transparent for unmarshalling.
	return json.Unmarshal(data, &s.SecretValue)
}

// Stringer interface

func (s Secret[T]) String() string {
	return fmt.Sprint(s.redactedValue)
}
