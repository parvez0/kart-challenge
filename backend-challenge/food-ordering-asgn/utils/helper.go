package utils

// Package utils provides utility functions and modules like reusable logger
// helper functions for ErrorWrappers, ToPtrs etc.

// The helper.go file implements helper functions for the utils package.
// It provides functions for creating error wrappers, converting to pointers,
// and other utility functions.

import (
	"fmt"
)


func WrapError(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

func ToPtr[T any](v T) *T {
	return &v
}

func FromPtr[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}