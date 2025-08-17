package timecapsule

import (
	"encoding/json"
)

// JSONCodec implements Codec using JSON encoding
type JSONCodec[T any] struct{}

// NewJSONCodec creates a new JSON codec
func NewJSONCodec[T any]() Codec[T] {
	return &JSONCodec[T]{}
}

// Encode serializes a value to JSON bytes
func (c *JSONCodec[T]) Encode(value T) ([]byte, error) {
	return json.Marshal(value)
}

// Decode deserializes JSON bytes to a value
func (c *JSONCodec[T]) Decode(data []byte) (T, error) {
	var value T
	err := json.Unmarshal(data, &value)
	return value, err
}
