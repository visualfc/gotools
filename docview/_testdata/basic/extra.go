package fixture

// Pair contains two values.
type Pair[A, B any] struct {
	First  A
	Second B
}

// NewPair constructs a Pair.
func NewPair[A, B any](a A, b B) Pair[A, B] { return Pair[A, B]{a, b} }
