package utils

// Package utils provides utility functions and modules like reusable logger
// helper functions for ErrorWrappers, ToPtrs etc.

// The helper.go file implements helper functions for the utils package.
// It provides functions for creating error wrappers, converting to pointers,
// and other utility functions.

import (
	"fmt"
	"reflect"
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

func IsNil(v any) bool {
	return v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}

func IsNotNil(v any) bool {
	return !IsNil(v)
}