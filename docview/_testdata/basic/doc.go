// Package fixture demonstrates documentation extraction.
package fixture

import "fmt"

// Identity returns its argument.
func Identity[T any](v T) T { return v }

// Box stores one value.
type Box[T any] struct {
	Value T
}

// Get returns the stored value.
func (b Box[T]) Get() T { return b.Value }

// Format formats a value.
func Format[T any](v T) string { return fmt.Sprint(v) }
